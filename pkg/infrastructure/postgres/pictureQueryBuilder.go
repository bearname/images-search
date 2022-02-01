package postgres

import (
	"photofinish/pkg/common/util/uuid"
	"photofinish/pkg/domain/pictures"
	"strconv"
	"time"
)

type query struct {
	sql  string
	data []interface{}
}

type searchQuery struct {
	CountQuery    query
	PicturesQuery query
	OnlyEventId   bool
}

type pictureQueryBuilder struct {
}

func (s *pictureQueryBuilder) getInsertTextDetectionDeclarationSQL() string {
	return "INSERT INTO pictures_text_detection (picturesId, detected_text, confidence) VALUES "
}
func (s *pictureQueryBuilder) buildStoreAllQuery(picturesList []*pictures.TextDetectionOnImageDto) (*query, error) {
	pictureCount := len(picturesList)
	if pictureCount == 0 {
		return nil, pictures.ErrEmptyImages
	}

	var data []interface{}
	var textData []interface{}
	var sql string
	var detectedTextSql string
	pictureI := 1
	textI := pictureCount*4 + 1

	sql += "INSERT INTO pictures (id, original_path, preview_path, eventId) VALUES "
	detectedTextSql += s.getInsertTextDetectionDeclarationSQL()
	for i, imageTextDetectionDto := range picturesList {
		sql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `, $` + strconv.Itoa(pictureI+3) + `)`
		if pictureCount > 1 && i < pictureCount-1 {
			sql += ", "
		} else {
			sql += "; "
		}
		pictureI += 4
		id := uuid.Generate().String()
		data = append(data, id, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.PreviewPath, imageTextDetectionDto.EventId)

		if imageTextDetectionDto.TextDetection != nil {
			for j, detection := range imageTextDetectionDto.TextDetection {
				detectedTextSql += " ($" + strconv.Itoa(textI) + ", $" + strconv.Itoa(textI+1) + ", $" + strconv.Itoa(textI+2) + ")"
				if i == pictureCount-1 && j == len(imageTextDetectionDto.TextDetection)-1 {
					detectedTextSql += "; "
				} else {
					detectedTextSql += ", "
				}

				textI += 3
				textData = append(textData, id, detection.DetectedText, detection.Confidence)
			}
		}
	}

	sql += detectedTextSql
	data = append(data, textData...)

	return &query{sql: sql, data: data}, nil
}

func (s *pictureQueryBuilder) buildStoreImageQuery(imageTextDetectionDto *pictures.TextDetectionOnImageDto) *query {
	sql := `DO $$
    DECLARE
        my_id pictures.id%TYPE;
    BEGIN
        INSERT INTO pictures (original_path, eventId) VALUES ($1, $2) RETURNING id INTO my_id;`

	i := 3

	var data []interface{}
	data = append(data, imageTextDetectionDto.OriginalPath, imageTextDetectionDto.EventId)

	for _, detection := range imageTextDetectionDto.TextDetection {
		sql += s.getInsertTextDetectionDeclarationSQL() + " (my_id, $" + strconv.Itoa(i) + ", $" + strconv.Itoa(i+1) + ");"
		i += 2
		data = append(data, detection.DetectedText, detection.Confidence)
	}

	sql += `END $$;`

	return &query{data: data, sql: sql}
}
func (s *pictureQueryBuilder) buildUpdateImageHandleQuery(picture *pictures.Picture) *query {
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
		sql += s.getInsertTextDetectionDeclarationSQL()
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

	return &query{data: data, sql: sql}
}

func (s *pictureQueryBuilder) buildSaveInitialPicturesQuery(initial *pictures.InitialDropboxImage) (*query, []uuid.UUID, *uuid.UUID, error) {
	images := initial.Images
	if len(images) == 0 {
		return nil, nil, nil, pictures.ErrEmptyImages
	}
	sql := `INSERT INTO tasks (id, dropbox_path, count_images) VALUES ($1, $2, $3);
            INSERT INTO pictures (id, dropbox_path, eventid, processing_status, task_id) VALUES `
	var data []interface{}

	tasksId := uuid.Generate()
	taskIdString := tasksId.String()
	data = append(data, taskIdString, initial.Path, len(images))
	pictureI := 4
	var ids []uuid.UUID
	for i, image := range images {
		sql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) +
			`, $` + strconv.Itoa(pictureI+3) + `, $` + strconv.Itoa(pictureI+4) + `)`
		if len(images) > 1 && i < len(images)-1 {
			sql += ", "
		} else {
			sql += ";"
		}
		pictureI += 5
		imgId := uuid.Generate()
		id := imgId.String()
		ids = append(ids, imgId)
		data = append(data, id, image, initial.EventId, pictures.Processing, taskIdString)
	}

	return &query{data: data, sql: sql}, ids, &tasksId, nil

}
func (s *pictureQueryBuilder) buildSaveInitialPictureQuery(image *pictures.InitialImage) *query {
	sql := "INSERT INTO pictures (id, dropbox_path, eventid) VALUES "
	pictureI := 1
	var data []interface{}
	sql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `, $` + strconv.Itoa(pictureI+3) + `) RETURNING ID;`
	id := uuid.Generate().String()
	data = append(data, id, image.DropboxPath, image.EventId)

	return &query{
		data: data, sql: sql,
	}
}
func (s *pictureQueryBuilder) buildSearchQuery(dto *pictures.SearchPictureDto) *searchQuery {
	var data []interface{}
	i := 1

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
		sqlWhere += ` WHERE position($` + strconv.Itoa(i) + ` in ptd.detected_text) > 0 AND confidence > $` + strconv.Itoa(i+1)
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
		sqlCount = `SELECT COUNT (DISTINCT id) AS count FROM pictures WHERE eventid = $1;`
		sql = `SELECT id, preview_path FROM pictures WHERE eventid = $1 LIMIT $2 OFFSET $3;`

		data = nil
		data = append(data, dto.EventId)
	} else {
		sqlCount = `SELECT COUNT (DISTINCT picturesId) AS count
            FROM pictures_text_detection AS ptd
            LEFT JOIN pictures p on ptd.picturesId = p.id
            LEFT JOIN events e on e.id = p.eventId ` + sqlWhere
	}
	data = append(data, dto.Page.Limit, dto.Page.Offset)

	return &searchQuery{
		CountQuery:    query{sql: sqlCount, data: dataWhere},
		PicturesQuery: query{sql: sql, data: data},
		OnlyEventId:   onlyEventId,
	}
}
