package router

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/infrastructure/transport"
)

type Controllers struct {
	PictureController *transport.PictureController
	EventsController  *transport.EventController
	AuthController    *transport.AuthController
	TasksController   *transport.TasksController
	OrderController   *transport.OrderController
}

func Router(controllers Controllers) http.Handler {
	tasksController := controllers.TasksController
	authController := controllers.AuthController
	eventsController := controllers.EventsController
	pictureController := controllers.PictureController
	orderController := controllers.OrderController
	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", readyCheckHandler).Methods(http.MethodGet)

	router.HandleFunc("/webhook", orderController.OnEventStripe()).Methods(http.MethodPost, http.MethodOptions)

	apiV1Route := router.PathPrefix("/api/v1").Subrouter()
	apiV1Route.HandleFunc("/yookassa", orderController.OnEventYookassa()).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/orders/{id}", authController.CheckTokenHandler(orderController.GetOrder())).Methods(http.MethodGet, http.MethodOptions)

	apiV1Route.HandleFunc("/tasks/{id}/stats", authController.CheckTokenHandler(tasksController.GetTaskStatistic())).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/tasks", authController.CheckTokenHandler(tasksController.GetTaskList())).Methods(http.MethodGet, http.MethodOptions)

	apiV1Route.HandleFunc("/charges", authController.CheckTokenHandler(orderController.Pay())).Methods(http.MethodPost, http.MethodOptions)

	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.List())).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.CreateEvent())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/events/{id}", authController.CheckTokenHandler(eventsController.DeleteEvent())).Methods(http.MethodDelete, http.MethodOptions)

	apiV1Route.HandleFunc("/pictures/search", pictureController.SearchPictures()).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/pictures/dropbox-folders", authController.CheckTokenHandler(pictureController.GetDropboxFolders())).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/pictures/detectText/dropbox", authController.CheckTokenHandler(pictureController.DetectImageFromDropboxUrl())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/pictures/{id}", authController.CheckTokenHandler(pictureController.DeletePicture())).Methods(http.MethodDelete, http.MethodOptions)

	apiV1Route.HandleFunc("/auth/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.RefreshToken)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token/validate", authController.ValidateToken).Methods(http.MethodGet, http.MethodOptions)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	return logMiddleware(router)
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func readyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "{\"host\": \"%v\"}", r.Host)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
