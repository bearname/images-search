package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/common/util"
	"photofinish/pkg/domain/dto"
	"photofinish/pkg/domain/tasks"
	"strconv"
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

func (c *TasksController) GetTaskStatistic() func(http.ResponseWriter, *http.Request) {
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

		stats, err := c.service.GetTaskStatistic(taskId)
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

func (c *TasksController) GetTaskList() func(http.ResponseWriter, *http.Request) {
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
		limit := 10
		var err error
		if len(limitParameter) != 0 {
			limit, err = strconv.Atoi(limitParameter)
			msg := "Invalid 'limit' query parameter. 'limit' must be in range [0, 100]"
			if err != nil {
				log.Println(err)
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			if limit < 0 || limit > 40 {
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
				log.Println(err)
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			if offset < 0 {
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}

		tasksList, err := c.service.GetTasks(&dto.Page{Offset: offset, Limit: limit})
		if err != nil {
			fmt.Println(err.Error())
			log.Println(err)
			http.Error(w, "Failed found pictures", http.StatusBadRequest)
			return
		}
		c.WriteJsonResponse(w, response{
			Code:    0,
			Message: "",
			Data:    tasksList,
		})
	}
}
