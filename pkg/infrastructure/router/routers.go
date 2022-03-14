package router

import (
	"fmt"
	"github.com/col3name/images-search/pkg/infrastructure/transport"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Controllers struct {
	PictureController *transport.PictureController
	EventsController  *transport.EventController
	AuthController    *transport.AuthController
	TasksController   *transport.TasksController
	OrderController   *transport.OrderController
}

func Router(controllers *Controllers) http.Handler {
	router := mux.NewRouter()
	orderController := controllers.OrderController

	postReq := []string{http.MethodPost, http.MethodOptions}

	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", readyCheckHandler).Methods(http.MethodGet)
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/webhook", orderController.OnEventStripe()).Methods(postReq...)

	apiV1Route := router.PathPrefix("/api/v1").Subrouter()

	for _, h := range routeList(controllers) {
		apiV1Route.HandleFunc(h.Path, h.Func).Methods(h.Methods...)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	return logMiddleware(router)
}

func routeList(controllers *Controllers) []*handleFunc {
	tasksCntrl := controllers.TasksController
	auth := controllers.AuthController
	eventsCntrl := controllers.EventsController
	pictureCntrl := controllers.PictureController
	orderCntrl := controllers.OrderController

	getReq := []string{http.MethodGet, http.MethodOptions}
	postReq := []string{http.MethodPost, http.MethodOptions}
	delReq := []string{http.MethodDelete, http.MethodOptions}

	return []*handleFunc{
		newHandleFunc("/yookassa", getReq, orderCntrl.OnEventYookassa()),
		newHandleFunc("/charges", postReq, withAuth(auth, orderCntrl.Pay())),

		newHandleFunc("/orders/{id}", getReq, withAuth(auth, orderCntrl.GetOrder())),

		newHandleFunc("/tasks/{id}/stats", getReq, withAuth(auth, tasksCntrl.GetTaskStatistic())),
		newHandleFunc("/tasks", getReq, withAuth(auth, tasksCntrl.GetTaskList())),

		newHandleFunc("/events", getReq, withAuth(auth, eventsCntrl.List())),
		newHandleFunc("/events/{id}", delReq, withAuth(auth, eventsCntrl.DeleteEvent())),

		newHandleFunc("/pictures/search", getReq, pictureCntrl.SearchPictures()),
		newHandleFunc("/pictures/dropbox-folders", getReq, withAuth(auth, pictureCntrl.GetDropboxFolders())),
		newHandleFunc("/pictures/detectText/dropbox", postReq, withAuth(auth, pictureCntrl.DetectImageFromDropboxUrl())),
		newHandleFunc("/pictures/{id}", delReq, withAuth(auth, pictureCntrl.DeletePicture())),

		newHandleFunc("/auth/login", postReq, auth.Login),
		newHandleFunc("/auth/token", getReq, withAuth(auth, auth.RefreshToken)),
		newHandleFunc("/auth/token/validate", getReq, auth.ValidateToken),
	}
}
func withAuth(authController *transport.AuthController, next http.HandlerFunc) http.HandlerFunc {
	return authController.CheckTokenHandler(next)
}

type handleFunc struct {
	Path    string
	Methods []string
	Func    func(http.ResponseWriter, *http.Request)
}

func newHandleFunc(path string, methods []string, f func(http.ResponseWriter, *http.Request)) *handleFunc {
	return &handleFunc{Path: path, Methods: methods, Func: f}
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
