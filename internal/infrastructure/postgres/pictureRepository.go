package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"photofinish/internal/common/uuid"
	"photofinish/internal/domain/picture"
	"strconv"
	"strings"
)

type PictureRepositoryImpl struct {
	connPool *pgx.ConnPool
}

func NewPictureRepository(connPool *pgx.ConnPool) *PictureRepositoryImpl {
	u := new(PictureRepositoryImpl)
	u.connPool = connPool
	return u
}

func (r *PictureRepositoryImpl) Search(dto picture.SearchPictureDto) (picture.SearchPictureResultDto, error) {
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
	if dto.ParticipantNumber != picture.ValueNotSetted {
		data = append(data, strconv.Itoa(dto.ParticipantNumber), dto.Confidence)
		dataWhere = append(dataWhere, strconv.Itoa(dto.ParticipantNumber), dto.Confidence)
		sqlWhere += ` WHERE position($` + strconv.Itoa(i) + ` in ptd.detected_text) > 0 AND  confidence > $` + strconv.Itoa(i+1)
		i += 2
		if dto.EventId != picture.ValueNotSetted {
			data = append(data, dto.EventId)
			dataWhere = append(dataWhere, dto.EventId)
			sqlWhere += ` AND e.id = $` + strconv.Itoa(i)
			i++
		}
	} else if dto.EventId != picture.ValueNotSetted {
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

	var sqlCount = ""
	if onlyEventId {
		sqlCount = `SELECT COUNT (DISTINCT id) AS COUNT FROM pictures WHERE eventid = $1;`
		sql = `SELECT id, preview_path FROM pictures WHERE eventid = $1 LIMIT $2 OFFSET $3;`

		data = nil
		data = append(data, dto.EventId)
	} else {
		sqlCount = `SELECT COUNT (DISTINCT picturesId) AS COUNT
            FROM pictures_text_detection
            LEFT JOIN pictures p on pictures_text_detection.picturesId = p.id
            LEFT JOIN events e on e.id = p.eventId ` + sqlWhere
	}
	rowsCount, err := r.connPool.Query(sqlCount, dataWhere...)
	var result picture.SearchPictureResultDto

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
	var pictureItem picture.SearchPictureItem
	var textDetectionsString string
	for rows.Next() {
		var arr []picture.TextDetectionDto
		if onlyEventId {
			err = rows.Scan(
				&pictureItem.PictureId,
				&pictureItem.Path,
			)
			if err != nil {
				return result, err
			}
			arr = make([]picture.TextDetectionDto, 0)
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
					atoi, err := strconv.ParseFloat(item[1], 64)
					if err != nil {
						return result, err
					}
					atoi1, err := strconv.ParseInt(item[2], 10, 64)
					if err != nil {
						return result, err
					}
					arr = append(arr, picture.TextDetectionDto{
						TextDetection: picture.TextDetection{
							DetectedText: item[0],
							Confidence:   atoi,
						},
						Event: picture.Event{
							EventId:   atoi1,
							EventName: item[3],
						},
					})
				}
			} else {
				arr = make([]picture.TextDetectionDto, 0)
			}

		}

		pictureItem.TextDetections = arr
		result.Pictures = append(result.Pictures, pictureItem)
	}

	return result, nil
}

func (r *PictureRepositoryImpl) Delete(pictureId string) error {
	sql := "DELETE FROM pictures WHERE id=$1"
	rows, err := r.connPool.Query(sql, pictureId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (r *PictureRepositoryImpl) FindById(pictureId string) error {
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
		return errors.New("picture not exist")
	}

	return nil
}

func (r *PictureRepositoryImpl) StoreAll(pictures []*picture.TextDetectionOnImageDto) error {
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
	//
	//pictureI := 1
	//for _, imageTextDetectionDto := range pictures {
	//    picturesSql += `DO $$
	//        DECLARE
	//            myid pictures.id%TYPE;
	//        BEGIN
	//            INSERT INTO pictures (original_path, preview_path, eventId) VALUES ($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `) RETURNING id INTO myid; `
	//
	//    pictureI += 2
	//    data = append(data, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.PreviewPath,  imageTextDetectionDto.EventId)
	//
	//    if imageTextDetectionDto.TextDetection != nil {
	//        picturesSql += "INSERT INTO pictures_text_detection (picturesId, detectedText, confidence) VALUES "
	//
	//        for j, detection := range imageTextDetectionDto.TextDetection {
	//            picturesSql += " (myid, $" + strconv.Itoa(pictureI) + ", $" + strconv.Itoa(pictureI+1) + ")"
	//            if j < len(imageTextDetectionDto.TextDetection)-1 {
	//                picturesSql += ", "
	//            } else {
	//                picturesSql += "; "
	//            }
	//            pictureI += 2
	//            data = append(data, detection.DetectedText, detection.Confidence)
	//        }
	//    }
	//
	//    picturesSql += `END $$; `
	//}

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
	defer tx.Rollback()
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
func (r *PictureRepositoryImpl) Store(imageTextDetectionDto *picture.TextDetectionOnImageDto) error {
	sql := `DO $$
    DECLARE
        myid pictures.id%TYPE;
    BEGIN
        INSERT INTO pictures (original_path, eventId) VALUES ($1, $2) RETURNING id INTO myid;`

	i := 3

	var data []interface{}
	data = append(data, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.EventId)

	for _, detection := range imageTextDetectionDto.TextDetection {
		sql += "INSERT INTO pictures_text_detection (picturesId, detected_text, confidence) VALUES (myid, $" + strconv.Itoa(i) + ", $" + strconv.Itoa(i+1) + ");"
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
	defer tx.Rollback()
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
