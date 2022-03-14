package transport

import (
	"encoding/json"
	"github.com/col3name/images-search/pkg/domain/dto"
	"github.com/col3name/images-search/pkg/domain/event"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type EventController struct {
	BaseController
	service event.Service
}

func NewEventController(service event.Service) *EventController {
	c := new(EventController)
	c.service = service
	return c
}

func (c *EventController) CreateEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		all, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid json body. ", http.StatusBadRequest)
			return
		}
		var t event.CreateEventInputDto
		err = json.Unmarshal(all, &t)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid json body. required field: name, location, date ", http.StatusBadRequest)
			return
		}
		eventId, err := c.service.Create(&t)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed create event", http.StatusBadRequest)
			return
		}

		c.WriteJsonResponse(w, responseWithoutData{1, "Success create event. TaskId " + strconv.Itoa(eventId)})
	}
}

func (c *EventController) DeleteEvent() func(http.ResponseWriter, *http.Request) {
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
			http.Error(w, "Invalid 'id' query parameter. 'Number' must be number", http.StatusBadRequest)
			return
		}

		eventId, err := strconv.Atoi(idString)
		if err != nil {
			log.Println(err)
			http.Error(w, "Event id must be positive number", http.StatusBadRequest)
			return
		}

		err = c.service.DeleteEvent(eventId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c.WriteJsonResponse(w, responseWithoutData{1, "Success delete event "})
	}
}

func (c *EventController) List() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		query := r.URL.Query()
		limitParameter := query.Get("limit")
		limit := 20
		var err error
		if len(limitParameter) != 0 {
			limit, err = strconv.Atoi(limitParameter)
			msg := "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]"
			if err != nil {
				log.Println(err, msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			if limit < 0 || limit > 100 {
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
		offsetParameter := query.Get("offset")
		offset := 0
		if len(offsetParameter) != 0 {
			offset, err = strconv.Atoi(offsetParameter)
			msg := "Invalid 'offset' query parameter. 'offset' must be in range [0, 100]"
			if err != nil {
				log.Println(err, msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			if offset < 0 {
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}

		events, err := c.service.Search(&dto.Page{Offset: offset, Limit: limit})
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed found pictures", http.StatusBadRequest)
			return
		}
		if events == nil {
			events = make([]event.Event, 0)
		}
		responseData := make(map[string]interface{})
		responseData["data"] = events
		c.WriteJsonResponse(w, responseData)
	}
}
