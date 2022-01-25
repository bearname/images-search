package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/common/util"
	"photofinish/pkg/domain/tasks"
)

type TasksController struct {
	BaseController
	service tasks.Service
}

func NewTasksController(service tasks.Service) *TasksController {
	c := new(TasksController)
	c.service = service
	return c
}

//
//func (c *TasksController) CreateEvent() func(http.ResponseWriter, *http.Request) {
//    return func(w http.ResponseWriter, r *http.Request) {
//        w.Header().Set("Access-Control-Allow-Origin", "*")
//        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
//        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//        if (*r).Method == "OPTIONS" {
//            w.WriteHeader(http.StatusNoContent)
//            return
//        }
//
//        all, err := ioutil.ReadAll(r.Body)
//        if err != nil {
//            log.Println(err)
//            http.Error(w, "Invalid json body. ", http.StatusBadRequest)
//            return
//        }
//        var t event.CreateEventInputDto
//        err = json.Unmarshal(all, &t)
//        if err != nil {
//            log.Println(err)
//            http.Error(w, "Invalid json body. required field: name, location, date ", http.StatusBadRequest)
//            return
//        }
//        eventId, err := c.service.Create(&t)
//        if err != nil {
//            log.Println(err)
//            http.Error(w, "Failed create event", http.StatusBadRequest)
//            return
//        }
//
//        c.WriteJsonResponse(w, responseWithoutData{1, "Success create event. TaskId " + strconv.Itoa(eventId)})
//    }
//}
//
//func (c *TasksController) DeleteEvent() func(http.ResponseWriter, *http.Request) {
//    return func(w http.ResponseWriter, r *http.Request) {
//        w.Header().Set("Access-Control-Allow-Origin", "*")
//        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, DELETE")
//        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//        if (*r).Method == "OPTIONS" {
//            w.WriteHeader(http.StatusNoContent)
//            return
//        }
//
//        vars := mux.Vars(r)
//        idString := vars["id"]
//        var err error
//        if len(idString) == 0 {
//            log.Println(err)
//            http.Error(w, "Invalid 'id' query parameter. 'Number' must be number", http.StatusBadRequest)
//            return
//        }
//
//        eventId, err := strconv.Atoi(idString)
//        if err != nil {
//            log.Println(err)
//            http.Error(w, "Event id must be positive number", http.StatusBadRequest)
//            return
//        }
//
//        err = c.service.DeleteEvent(eventId)
//        if err != nil {
//            fmt.Println(err.Error())
//            log.Println(err)
//            http.Error(w, err.Error(), http.StatusBadRequest)
//            return
//        }
//
//        c.WriteJsonResponse(w, responseWithoutData{1, "Success delete event "})
//    }
//}

func (c *TasksController) GetStatistic() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if (*r).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		vars := mux.Vars(r)
		taskId := vars["id"]

		if len(taskId) == 0 || !util.IsUUID(taskId) {
			e := "Invalid 'id' query parameter. 'id' must be uuid"
			log.Println(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}

		stats, err := c.service.GetStatistics(taskId)
		if err != nil {
			fmt.Println(err.Error())
			log.Println(err)
			http.Error(w, "Failed found pictures", http.StatusBadRequest)
			return
		}

		responseData := make(map[string]interface{})
		responseData["data"] = stats
		c.WriteJsonResponse(w, responseData)
	}
}
