package transport

import (
	"github.com/col3name/images-search/pkg/domain/tasks"
	transpUtil "github.com/col3name/images-search/pkg/infrastructure/transport/util"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		taskId, err := transpUtil.GetUUIDParam(mux.Vars(req))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stats, err := c.service.GetTaskStatistic(taskId)
		if err != nil {
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
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		page, err := transpUtil.DecodePageReq(req)
		if err != nil {
			msg := err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		tasksList, err := c.service.GetTasks(page)
		if err != nil {
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
