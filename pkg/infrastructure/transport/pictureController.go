package transport

import (
	"errors"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/col3name/images-search/pkg/infrastructure/transport/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PictureController struct {
	BaseController
	service pictures.Service
}

func NewPictureController(service pictures.Service) *PictureController {
	c := new(PictureController)
	c.service = service
	return c
}

func (c *PictureController) DetectImageFromDropboxUrl() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		dropboxUrl, eventId, err := c.decodeDetectImageFromDropboxUrlReq(req)
		if err != nil {
			msg := err.Error()
			log.Println(msg, err)
			http.Error(w, msg, http.StatusBadRequest)
		}
		taskResp, err := c.service.DetectImageFromUrl(dropboxUrl, eventId)
		if err != nil {
			c.BaseController.WriteError(w, err, TranslateError(err))
			return
		}

		c.WriteJsonResponse(w, response{
			Code:    http.StatusOK,
			Message: "Folder '" + dropboxUrl + "' processing",
			Data:    taskResp,
		})
	}
}

func (c *PictureController) decodeDetectImageFromDropboxUrlReq(req *http.Request) (string, int, error) {
	query := req.URL.Query()
	dropboxUrl := query.Get("path")

	if len(dropboxUrl) == 0 {
		return "", 0, errors.New("failed get url to dropbox zip file")
	}

	eventIdValue := query.Get("eventId")
	if len(eventIdValue) == 0 {
		return "", 0, errors.New("failed get eventId")
	}

	eventId, err := strconv.Atoi(eventIdValue)
	if err != nil {
		return "", 0, errors.New("invalid event id")
	}
	return dropboxUrl, eventId, nil
}

func (c *PictureController) SearchPictures() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		searchDTO, err := c.decodeSearchPicturesReq(req)
		if err != nil {
			msg := err.Error()
			log.Println(err, msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		result, err := c.service.Search(searchDTO)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed found searchDto", http.StatusBadRequest)
			return
		}
		if result.Pictures == nil {
			result.Pictures = make([]pictures.SearchPictureItem, 0)
		}
		c.WriteJsonResponse(w, result)
	}
}

func (c *PictureController) decodeSearchPicturesReq(req *http.Request) (*pictures.SearchPictureDto, error) {
	query := req.URL.Query()
	participantNumber, err := util.GetQueryParameter(query, "number", pictures.ValueNotSet, "invalid 'number' query parameter. 'Number' must be number", func(val int) bool {
		return val < 0 || val > 100
	})

	confidence, err := util.GetQueryParameter(query, "confidence", 85, "invalid 'confidence' query parameter. 'limit' must be in range [0, 100]", func(val int) bool {
		return val < 0 || val > 100
	})
	if err != nil {
		return nil, err
	}
	page, err := util.DecodePageReq(req)
	if err != nil {
		return nil, err
	}
	eventId, err := util.GetPositiveNum(query, "eventId", pictures.ValueNotSet, "invalid 'eventId' query parameter.")
	if err != nil {
		return nil, err
	}
	searchDTO := pictures.NewSearchPictureDto(participantNumber, confidence, eventId, *page)
	return &searchDTO, nil
}

func (c *PictureController) GetDropboxFolders() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		folders, err := c.service.GetDropboxFolders()
		if err != nil {
			http.Error(w, "Failed found pictures"+err.Error(), http.StatusBadRequest)
			return
		}
		c.WriteJsonResponse(w, folders)
	}
}

func (c *PictureController) DeletePicture() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		vars := mux.Vars(req)
		idString := vars["id"]
		var err error
		if len(idString) == 0 {
			log.Println(err)
			http.Error(w, "invalid 'number' query parameter. 'Number' must be number", http.StatusBadRequest)
			return
		}

		err = c.service.Delete(idString)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed found pictures", http.StatusBadRequest)
			return
		}

		c.WriteJsonResponse(w, responseWithoutData{1, "Success delete "})
	}
}
