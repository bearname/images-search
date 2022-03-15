package transport

import (
	"encoding/json"
	"errors"
	"github.com/col3name/images-search/pkg/domain/event"
	"github.com/col3name/images-search/pkg/infrastructure/transport/util"
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
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		t, err := c.decodeCreateEventReq(req)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		eventId, err := c.service.Create(t)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed create event", http.StatusBadRequest)
			return
		}

		c.WriteJsonResponse(w, responseWithoutData{1, "Success create event. TaskId " + strconv.Itoa(eventId)})
	}
}

func (c *EventController) decodeCreateEventReq(req *http.Request) (*event.CreateEventInputDto, error) {
	all, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("invalid json body. ")
	}
	var t event.CreateEventInputDto
	err = json.Unmarshal(all, &t)
	if err != nil {
		return nil, errors.New("invalid json body. required field: name, location, date ")
	}
	return &t, nil
}

func (c *EventController) DeleteEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		eventId, err := c.decodeDeleteEventReq(req)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
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

func (c *EventController) decodeDeleteEventReq(req *http.Request) (int, error) {
	vars := mux.Vars(req)
	idString := vars["id"]
	var err error
	if len(idString) == 0 {
		return 0, errors.New("invalid 'id' query parameter. 'Number' must be number")
	}

	eventId, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("event id must be positive number")
	}
	return eventId, nil
}

func (c *EventController) List() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		listReq, err := util.DecodePageReq(req)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		events, err := c.service.Search(listReq)
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
