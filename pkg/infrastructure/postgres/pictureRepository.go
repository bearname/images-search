package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/common/util/uuid"
	"photofinish/pkg/domain/pictures"
	"strconv"
	"strings"
	"time"
)

type PictureRepositoryImpl struct {
	connPool *pgx.ConnPool
}

func NewPictureRepository(connPool *pgx.ConnPool) *PictureRepositoryImpl {
	u := new(PictureRepositoryImpl)
	u.connPool = connPool
	return u
}

func (r *PictureRepositoryImpl) IsExists(pictureId string) error {
	sql := "SELECT id FROM pictures WHERE id=$1"
	rows, err := r.connPool.Query(sql, pictureId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	if !rows.Next() {
		return errors.New("pictures not exist")
	}

	return nil
}

func (r *PictureRepositoryImpl) Search(dto pictures.SearchPictureDto) (pictures.SearchPictureResultDto, error) {
	var data []interface{}
	i := 1
	fmt.Println("dto")
	fmt.Println(dto)
	sql := `SELECT picturesId, preview_path, 
        string_agg(DISTINCT detected_text || ':' || confidence::varchar(255) || ':' || eventId|| ':' || e.name, ','::varchar(255)) AS detected_texts
        FROM pictures
            LEFT JOIN pictures_text_detection ptd on pictures.id = ptd.picturesId
            LEFT JOIN events e on e.id = pictures.eventId`

	var sqlWhere string
	var dataWhere []interface{}
	onlyEventId := false
	if dto.ParticipantNumber != pictures.ValueNotSet {
		data = append(data, strconv.Itoa(dto.ParticipantNumber), dto.Confidence)
		dataWhere = append(dataWhere, strconv.Itoa(dto.ParticipantNumber), dto.Confidence)
		sqlWhere += ` WHERE position($` + strconv.Itoa(i) + ` in ptd.detected_text) > 0 AND  confidence > $` + strconv.Itoa(i+1)
		i += 2
		if dto.EventId != pictures.ValueNotSet {
			data = append(data, dto.EventId)
			dataWhere = append(dataWhere, dto.EventId)
			sqlWhere += ` AND e.id = $` + strconv.Itoa(i)
			i++
		}
	} else if dto.EventId != pictures.ValueNotSet {
		data = append(data, dto.EventId)
		dataWhere = append(dataWhere, dto.EventId)
		sqlWhere += ` WHERE e.id = $` + strconv.Itoa(i)
		i++
		onlyEventId = true
	} else {
		sqlWhere += ` WHERE confidence > $` + strconv.Itoa(i)
		data = append(data, dto.Confidence)
		dataWhere = append(dataWhere, dto.Confidence)
		i++
	}
	sql += sqlWhere + ` GROUP BY picturesId, preview_path LIMIT $` + strconv.Itoa(i) + ` OFFSET $` + strconv.Itoa(i+1) + `;`

	var sqlCount string
	if onlyEventId {
		sqlCount = `SELECT COUNT (DISTINCT id) AS COUNT FROM pictures WHERE eventid = $1;`
		sql = `SELECT id, preview_path FROM pictures WHERE eventid = $1 LIMIT $2 OFFSET $3;`

		data = nil
		data = append(data, dto.EventId)
	} else {
		sqlCount = `SELECT COUNT (DISTINCT picturesId) AS COUNT
            FROM pictures_text_detection AS ptd
            LEFT JOIN pictures p on ptd.picturesId = p.id
            LEFT JOIN events e on e.id = p.eventId ` + sqlWhere
	}
	fmt.Println("sqlCount")
	fmt.Println(sqlCount)
	rowsCount, err := r.connPool.Query(sqlCount, dataWhere...)
	var result pictures.SearchPictureResultDto

	if err != nil {
		return result, err
	}

	defer rowsCount.Close()

	if rowsCount.Err() != nil {
		return result, rowsCount.Err()
	}

	if rowsCount.Next() {
		err = rowsCount.Scan(&result.CountAllItems)
		if err != nil {
			return result, err
		}
	}

	//"SELECT picturesId, detectedText, confidence, eventId, name\n\t\tFROM pictures\n\t\t\tLEFT JOIN pictures_text_detection ptd on pictures.id = ptd.picturesId\n\t\t\tLEFT JOIN events e on e.id = pictures.eventId\n\t\tWHERE ptd.detectedText = $1 AND  confidence > $2 LIMIT $3 OFFSET $4;"
	data = append(data, dto.Page.Limit, dto.Page.Offset)
	fmt.Println("sql")
	fmt.Println(sql)
	fmt.Println("data")
	fmt.Println(data)
	rows, err := r.connPool.Query(sql, data...)

	if err != nil {
		return result, err
	}

	defer rows.Close()

	if rows.Err() != nil {
		return result, rows.Err()
	}
	var pictureItem pictures.SearchPictureItem
	var textDetectionsString string
	for rows.Next() {
		var arr []pictures.TextDetectionDto
		if onlyEventId {
			err = rows.Scan(
				&pictureItem.PictureId,
				&pictureItem.Path,
			)
			if err != nil {
				return result, err
			}
			arr = make([]pictures.TextDetectionDto, 0)
		} else {
			err = rows.Scan(
				&pictureItem.PictureId,
				&pictureItem.Path,
				&textDetectionsString,
			)

			if err != nil {
				s := err.Error()
				if !strings.Contains(s, "cannot assign NULL to *string") {
					continue
				} else {
					textDetectionsString = ""
				}
			}
			if len(textDetectionsString) > 0 {
				textDetectionItem := strings.Split(textDetectionsString, ",")
				for _, textDetection := range textDetectionItem {
					item := strings.Split(textDetection, ":")
					confidence, err := strconv.ParseFloat(item[1], 64)
					if err != nil {
						return result, err
					}
					eventId, err := strconv.ParseInt(item[2], 10, 64)
					if err != nil {
						return result, err
					}
					arr = append(arr, pictures.TextDetectionDto{
						TextDetection: pictures.TextDetection{
							DetectedText: item[0],
							Confidence:   confidence,
						},
						Event: pictures.Event{
							EventId:   eventId,
							EventName: item[3],
						},
					})
				}
			} else {
				arr = make([]pictures.TextDetectionDto, 0)
			}

		}

		pictureItem.TextDetections = arr
		result.Pictures = append(result.Pictures, pictureItem)
	}

	return result, nil
}

func (r *PictureRepositoryImpl) SaveInitialPicture(image *pictures.InitialImage) (int, error) {
	picturesSql := "INSERT INTO pictures (id, dropbox_path, eventid) VALUES"
	pictureI := 1
	var data []interface{}
	picturesSql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `, $` + strconv.Itoa(pictureI+3) + `) RETURNING ID;`
	id := uuid.Generate().String()
	data = append(data, id, image.DropboxPath, image.EventId)

	tx, err := r.connPool.Begin()
	if err != nil {
		return 0, err
	}
	var orderID int
	err = tx.QueryRow(picturesSql, data...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, tx.Commit()
}

func (r *PictureRepositoryImpl) SaveInitialPictures(initial pictures.InitialDropboxImage) (*pictures.InitialDropboxImageResult, error) {
	images := initial.Images
	if len(images) == 0 {
		return nil, errors.New("empty images")
	}
	picturesSql := `INSERT INTO tasks (id, dropbox_path, count_images) VALUES ($1, $2, $3);
                    INSERT INTO pictures (id, dropbox_path, eventid, processing_status, task_id) VALUES `
	var data []interface{}

	tasksId := uuid.Generate()
	taskIdString := tasksId.String()
	data = append(data, taskIdString, initial.Path, len(images))
	pictureI := 4
	var ids []uuid.UUID
	for i, image := range images {
		picturesSql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) +
			`, $` + strconv.Itoa(pictureI+3) + `, $` + strconv.Itoa(pictureI+4) + `)`
		if len(images) > 1 && i < len(images)-1 {
			picturesSql += ", "
		} else {
			picturesSql += ";"
		}
		pictureI += 5
		imgId := uuid.Generate()
		id := imgId.String()
		ids = append(ids, imgId)
		data = append(data, id, image, initial.EventId, pictures.Processing, taskIdString)
	}

	tx, err := r.connPool.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(picturesSql, data...)
	if err != nil {
		return nil, err
	}

	var result pictures.InitialDropboxImageResult
	result.ImagesId = ids
	result.TaskId = tasksId
	return &result, tx.Commit()
}

func (r *PictureRepositoryImpl) UpdateImageHandle(picture *pictures.Picture) error {
	sql := `UPDATE pictures SET
                original_path = $1,
                preview_path = $2,
                is_original_saved = $3, 
                is_preview_saved= $4,
                is_text_recognized = $5,
                attempts = $6,
                processing_status = $7,
                execute_after = $8,
                update_at = $9
            WHERE id = $10;`

	var data []interface{}

	data = append(data, picture.OriginalPath,
		picture.PreviewPath,
		picture.IsOriginalSaved,
		picture.IsPreviewSaved,
		picture.IsTextRecognized,
		picture.Attempts,
		picture.ProcessingStatus,
		picture.ExecuteAfter,
		time.Now(),
		picture.Id.String(),
	)

	textI := 11
	texts := picture.DetectedTexts
	if len(texts) > 0 {
		sql += "INSERT INTO pictures_text_detection (picturesId, detected_text, confidence) VALUES "
		for i, detection := range texts {
			sql += " ($" + strconv.Itoa(textI) + ", $" + strconv.Itoa(textI+1) + ", $" + strconv.Itoa(textI+2) + ")"
			if i == len(texts)-1 {
				sql += "; "
			} else {
				sql += ", "
			}

			textI += 3
			data = append(data, picture.Id.String(), detection.DetectedText, detection.Confidence)
		}
	}

	err := db.WithTransaction(r.connPool, func(tx *pgx.Tx) error {
		_, err := tx.Exec(sql, data...)
		if err != nil {
			fmt.Println(err)
		}
		return err
	})
	//
	//tx, err := r.connPool.Begin()
	//if err != nil {
	//    return err
	//}
	//defer func(tx *pgx.Tx) {
	//    err = tx.Rollback()
	//    if err != nil {
	//        log.Println(err)
	//    }
	//}(tx)
	//_, err = tx.Exec(sql, data...)
	//if err != nil {
	//    return err
	//}
	//
	//err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *PictureRepositoryImpl) Store(imageTextDetectionDto *pictures.TextDetectionOnImageDto) error {
	sql := `DO $$
    DECLARE
        my_id pictures.id%TYPE;
    BEGIN
        INSERT INTO pictures (original_path, eventId) VALUES ($1, $2) RETURNING id INTO my_id;`

	i := 3

	var data []interface{}
	data = append(data, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.EventId)

	for _, detection := range imageTextDetectionDto.TextDetection {
		sql += "INSERT INTO pictures_text_detection (picturesId, detected_text, confidence) VALUES (my_id, $" + strconv.Itoa(i) + ", $" + strconv.Itoa(i+1) + ");"
		i += 2
		data = append(data, detection.DetectedText, detection.Confidence)
	}

	sql += `END $$;`

	fmt.Println(sql)
	fmt.Println(data)
	tx, err := r.connPool.Begin()
	if err != nil {
		return err
	}
	defer func(tx *pgx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Println(err)
		}
	}(tx)
	_, err = tx.Exec(sql, data...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func (r *PictureRepositoryImpl) StoreAll(pictures []*pictures.TextDetectionOnImageDto) error {
	if len(pictures) == 0 {
		return errors.New("empty pictures")
	}

	var data []interface{}
	var textData []interface{}
	var picturesSql string
	var detectedTextSql string
	pictureI := 1
	textI := len(pictures)*4 + 1

	picturesSql += "INSERT INTO pictures (id, original_path, preview_path, eventId) VALUES "
	detectedTextSql += "INSERT INTO pictures_text_detection (picturesId, detected_text, confidence) VALUES "
	for i, imageTextDetectionDto := range pictures {
		picturesSql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `, $` + strconv.Itoa(pictureI+3) + `)`
		if len(pictures) > 1 && i < len(pictures)-1 {
			picturesSql += ", "
		} else {
			picturesSql += "; "
		}
		pictureI += 4
		id := uuid.Generate().String()
		data = append(data, id, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.PreviewPath, imageTextDetectionDto.EventId)

		if imageTextDetectionDto.TextDetection != nil {
			for j, detection := range imageTextDetectionDto.TextDetection {
				detectedTextSql += " ($" + strconv.Itoa(textI) + ", $" + strconv.Itoa(textI+1) + ", $" + strconv.Itoa(textI+2) + ")"
				if i == len(pictures)-1 && j == len(imageTextDetectionDto.TextDetection)-1 {
					detectedTextSql += "; "
				} else {
					detectedTextSql += ", "
				}

				textI += 3
				textData = append(textData, id, detection.DetectedText, detection.Confidence)
			}
		}
	}

	picturesSql += detectedTextSql
	fmt.Println(len(data))
	data = append(data, textData...)
	fmt.Println(picturesSql)
	fmt.Println(data)
	fmt.Println(len(data))

	tx, err := r.connPool.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(picturesSql, data...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *PictureRepositoryImpl) Delete(imageId string) error {
	sql := "DELETE FROM pictures WHERE id=$1"
	rows, err := r.connPool.Query(sql, imageId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}
