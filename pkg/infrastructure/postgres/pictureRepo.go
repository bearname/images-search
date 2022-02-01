package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/domain/pictures"
	"strconv"
	"strings"
)

type PictureRepositoryImpl struct {
	connPool     *pgx.ConnPool
	queryBuilder pictureQueryBuilder
}

func NewPictureRepository(connPool *pgx.ConnPool) *PictureRepositoryImpl {
	u := new(PictureRepositoryImpl)
	u.connPool = connPool
	u.queryBuilder = pictureQueryBuilder{}
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
		return pictures.ErrNotFound
	}

	return nil
}

func (r *PictureRepositoryImpl) FindPicture(pictureId string) (*pictures.PictureDTO, error) {
	sql := `SELECT id, eventid, task_id, dropbox_path, 
       original_s3_id, original_path, preview_s3_id, preview_path, attempts, 
       processing_status, is_original_saved,
       is_preview_saved, is_text_recognized FROM pictures
WHERE id = $1`
	var data []interface{}
	data = append(data, pictureId)
	i, err := db.Query(r.connPool, sql, data, func(rows *pgx.Rows) (interface{}, error) {
		if !rows.Next() {
			return nil, pictures.ErrNotFound
		}
		var dto pictures.PictureDTO
		err := rows.Scan(&dto.Id,
			&dto.EventId,
			&dto.TaskId,
			&dto.DropboxPath,
			&dto.OriginalS3Id,
			&dto.OriginalPath,
			&dto.PreviewS3Id,
			&dto.PreviewPath,
			&dto.Attempts,
			&dto.ProcessingStatus,
			&dto.IsOriginalSaved,
			&dto.IsPreviewSaved,
			&dto.IsTextRecognized)
		if err != nil {
			return nil, err
		}
		return &dto, nil
	})
	if err != nil {
		return nil, err
	}
	dto := i.(pictures.PictureDTO)

	return &dto, nil
}

func (r *PictureRepositoryImpl) Search(dto *pictures.SearchPictureDto) (*pictures.SearchPictureResultDto, error) {
	fmt.Println("dto")
	fmt.Println(dto)
	queryDTO, err := r.queryBuilder.buildSearchQuery(dto)
	if err != nil {
		return nil, err
	}

	fmt.Println("sqlCount")
	fmt.Println(queryDTO.CountQuery.sql)
	var result pictures.SearchPictureResultDto

	i, err := db.Query(r.connPool, queryDTO.CountQuery.sql, queryDTO.CountQuery.data, r.scanCountSearchImages(err))
	if err != nil {
		return nil, err
	}

	switch i.(type) {
	case int:
		fmt.Println("int")
	}
	result.CountAllItems = i.(int)

	sql := queryDTO.PicturesQuery.sql
	data := queryDTO.PicturesQuery.data
	res, err := db.Query(r.connPool, sql, data, r.scanSearchQueryResult(queryDTO))
	if err != nil {
		return nil, err
	}
	result = res.(pictures.SearchPictureResultDto)

	return &result, nil
}

func (r *PictureRepositoryImpl) SaveInitialPicture(image *pictures.InitialImage) (int, error) {
	pictureQuery := r.queryBuilder.buildSaveInitialPictureQuery(image)

	tx, err := r.connPool.Begin()
	if err != nil {
		return 0, err
	}
	var orderID int
	err = tx.QueryRow(pictureQuery.sql, pictureQuery.data...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, tx.Commit()
}

func (r *PictureRepositoryImpl) SaveInitialPictures(initial *pictures.InitialDropboxImage) (*pictures.InitialDropboxImageResult, error) {
	queryData, ids, tasksId, err := r.queryBuilder.buildSaveInitialPicturesQuery(initial)
	if err != nil {
		return nil, err
	}

	err = db.WithTransactionSQL(r.connPool, queryData.sql, queryData.data)
	if err != nil {
		return nil, err
	}

	return &pictures.InitialDropboxImageResult{
		ImagesId: ids,
		TaskId:   *tasksId,
	}, nil
}

func (r *PictureRepositoryImpl) UpdateImageHandle(picture *pictures.Picture) error {
	handleQuery := r.queryBuilder.buildUpdateImageHandleQuery(picture)

	err := db.WithTransactionSQL(r.connPool, handleQuery.sql, handleQuery.data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *PictureRepositoryImpl) Store(imageTextDetectionDto *pictures.TextDetectionOnImageDto) error {
	queryPg := r.queryBuilder.buildStoreImageQuery(imageTextDetectionDto)
	fmt.Println(queryPg.sql)
	fmt.Println(queryPg.data)

	err := db.WithTransactionSQL(r.connPool, queryPg.sql, queryPg.data)

	if err != nil {
		return err
	}

	return nil
}

func (r *PictureRepositoryImpl) StoreAll(pictures []*pictures.TextDetectionOnImageDto) error {
	pgQuery, err := r.queryBuilder.buildStoreAllQuery(pictures)
	if err != nil {
		return err
	}
	return db.WithTransactionSQL(r.connPool, pgQuery.sql, pgQuery.data)
}

func (r *PictureRepositoryImpl) Delete(imageId string) error {
	//sql := "DELETE FROM pictures WHERE id=$1"
	sql := "UPDATE pictures SET processing_status = $1 WHERE id = $2"

	err := db.WithTransaction(r.connPool, func(tx *pgx.Tx) error {
		_, err := tx.Exec(sql, pictures.Deleted, imageId)
		return err
	})
	return err
}

func (r *PictureRepositoryImpl) scanCountSearchImages(err error) func(rows *pgx.Rows) (interface{}, error) {
	return func(rows *pgx.Rows) (interface{}, error) {
		var count int
		if rows.Next() {
			err = rows.Scan(&count)
			if err != nil {
				return nil, err
			}
		}
		return count, nil
	}
}

func (r *PictureRepositoryImpl) scanSearchQueryResult(queryDTO *searchQuery) func(rows *pgx.Rows) (interface{}, error) {
	return func(rows *pgx.Rows) (interface{}, error) {
		var res pictures.SearchPictureResultDto

		var pictureItem pictures.SearchPictureItem
		for rows.Next() {
			arrs, err := r.buildTextDetectionList(rows, queryDTO.OnlyEventId)
			if err != nil {
				switch err {
				case pictures.ErrCannotAssignNull:
					continue
				default:
					return nil, err
				}
			}
			pictureItem.TextDetections = *arrs
			res.Pictures = append(res.Pictures, pictureItem)
		}

		return res, nil
	}
}

func (r *PictureRepositoryImpl) buildDetectedText(textDetection string) (pictures.TextDetectionDto, error) {
	item := strings.Split(textDetection, ":")
	confidence, err := strconv.ParseFloat(item[1], 64)
	if err != nil {
		return pictures.TextDetectionDto{}, err
	}
	eventId, err := strconv.ParseInt(item[2], 10, 64)
	if err != nil {
		return pictures.TextDetectionDto{}, err
	}
	detectionDto := pictures.TextDetectionDto{
		TextDetection: pictures.TextDetection{
			DetectedText: item[0],
			Confidence:   confidence,
		},
		Event: pictures.Event{
			EventId:   eventId,
			EventName: item[3],
		},
	}
	return detectionDto, nil
}

func (r *PictureRepositoryImpl) buildTextDetectionList(rows *pgx.Rows, onlyEventId bool) (*[]pictures.TextDetectionDto, error) {
	var pictureItem pictures.SearchPictureItem
	var textDetectionsString string

	var arr []pictures.TextDetectionDto
	if onlyEventId {
		err := rows.Scan(
			&pictureItem.PictureId,
			&pictureItem.Path,
		)
		if err != nil {
			return nil, err
		}

		arr = make([]pictures.TextDetectionDto, 0)
	} else {
		err := rows.Scan(
			&pictureItem.PictureId,
			&pictureItem.Path,
			&textDetectionsString,
		)

		if err != nil {
			s := err.Error()
			if !strings.Contains(s, "cannot assign NULL to *string") {
				return nil, pictures.ErrCannotAssignNull
			} else {
				textDetectionsString = ""
			}
		}
		if len(textDetectionsString) > 0 {
			textDetectionItem := strings.Split(textDetectionsString, ",")
			var detectionDto pictures.TextDetectionDto
			for _, textDetection := range textDetectionItem {
				detectionDto, err = r.buildDetectedText(textDetection)
				if err != nil {
					return nil, err
				}
				arr = append(arr, detectionDto)
			}
		} else {
			arr = make([]pictures.TextDetectionDto, 0)
		}
	}

	return &arr, nil
}
