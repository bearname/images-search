package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/domain/dto"
	"photofinish/pkg/domain/pictures"
	"strconv"
	"time"
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
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		query := r.URL.Query()
		dropboxUrl := query.Get("path")

		if len(dropboxUrl) == 0 {
			log.Println("err")
			http.Error(w, "Failed get url to dropbox zip file", http.StatusBadRequest)
			return
		}

		eventIdValue := query.Get("eventId")
		if len(eventIdValue) == 0 {
			log.Println("Failed get eventId")
			http.Error(w, "Failed get eventId", http.StatusBadRequest)
			return
		}

		eventId, err := strconv.Atoi(eventIdValue)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid event id", http.StatusBadRequest)
			return
		}

		start := time.Now()
		taskResp, err := c.service.DetectImageFromUrl(dropboxUrl, eventId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("\n\n\n\n\n\n\n\n\n\n\nSeconds:")
		end := time.Since(start)
		fmt.Println(end.Seconds())
		fmt.Println("Minutes:")
		fmt.Println(end.Minutes())
		c.WriteJsonResponse(w, response{
			Code:    http.StatusOK,
			Message: "Folder '" + dropboxUrl + "' processing",
			Data:    taskResp,
		})
	}
}

func (c *PictureController) SearchPictures() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		query := r.URL.Query()
		participantNumberStr := query.Get("number")
		participantNumber := pictures.ValueNotSet
		var err error
		if len(participantNumberStr) != 0 {
			participantNumber, err = strconv.Atoi(participantNumberStr)
			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid 'number' query parameter. 'Number' must be number", http.StatusBadRequest)
				return
			}
		}

		confidenceParameter := query.Get("confidence")
		confidence := 85
		if len(confidenceParameter) != 0 {
			confidence, err = strconv.Atoi(confidenceParameter)
			if err != nil {
				log.Println(err)
				log.Println("Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]")
				http.Error(w, "Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
			if confidence < 0 || confidence > 100 {
				log.Println("Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]")
				http.Error(w, "Invalid 'confidence' query parameter. 'confidence' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
		}
		limitParameter := query.Get("limit")
		limit := 20
		if len(limitParameter) != 0 {
			limit, err = strconv.Atoi(limitParameter)
			if err != nil {
				log.Println(err)
				log.Println("Invalid 'limit' query parameter. 'limit' must be in range [0, 100]")
				http.Error(w, "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
			if limit < 0 || limit > 100 {
				log.Println("Invalid 'limit' query parameter. 'limit' must be in range [0, 100]")
				http.Error(w, "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
		}
		offsetParameter := query.Get("offset")
		offset := 0
		if len(offsetParameter) != 0 {
			offset, err = strconv.Atoi(offsetParameter)
			if err != nil {
				log.Println(err)
				log.Println("Invalid 'offset' query parameter. 'offset' must be in range [0, 100]")
				http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
			if offset < 0 {
				log.Println("Invalid 'offset' query parameter. 'offset' must be in range [0, 100]")
				http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
		}

		eventIdParameter := query.Get("eventId")
		eventId := pictures.ValueNotSet
		if len(eventIdParameter) != 0 {
			eventId, err = strconv.Atoi(eventIdParameter)
			if err != nil {
				log.Println(err)
				log.Println("Invalid 'eventId' query parameter. 'offset' must be in range [0, 100]")
				http.Error(w, "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]", http.StatusBadRequest)
				return
			}
			if eventId < 0 {
				log.Println("Invalid 'eventId' query parameter. 'eventId' must be more than 0")
				http.Error(w, "Invalid 'eventId' query parameter. 'eventId' must be more than 0", http.StatusBadRequest)
				return
			}
		}

		searchDTO := pictures.NewSearchPictureDto(participantNumber, confidence, eventId, dto.Page{Limit: limit, Offset: offset})
		searchDto, err := c.service.Search(searchDTO)
		if err != nil {
			fmt.Println(err.Error())
			log.Println(err)
			http.Error(w, "Failed found searchDto", http.StatusBadRequest)
			return
		}

		fmt.Println(searchDto)
		if searchDto.Pictures == nil {
			searchDto.Pictures = make([]pictures.SearchPictureItem, 0)
		}

		c.WriteJsonResponse(w, searchDto)
	}
}

func (c *PictureController) GetDropboxFolders() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
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
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		vars := mux.Vars(r)
		idString := vars["id"]
		var err error
		if len(idString) == 0 {
			log.Println(err)
			http.Error(w, "Invalid 'number' query parameter. 'Number' must be number", http.StatusBadRequest)
			return
		}

		err = c.service.Delete(idString)
		if err != nil {
			fmt.Println(err.Error())
			log.Println(err)
			http.Error(w, "Failed found pictures", http.StatusBadRequest)
			return
		}

		c.WriteJsonResponse(w, responseWithoutData{1, "Success delete "})
	}
}
