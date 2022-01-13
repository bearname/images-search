package router

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/infrastructure/transport"
)

func Router(pictureController *transport.PictureController,
	eventsController *transport.EventController,
	authController *transport.AuthController) http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", readyCheckHandler).Methods(http.MethodGet)

	apiV1Route := router.PathPrefix("/api/v1").Subrouter()
	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.List())).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.CreateEvent())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/events/{id}", authController.CheckTokenHandler(eventsController.DeleteEvent())).Methods(http.MethodDelete, http.MethodOptions)

	apiV1Route.HandleFunc("/picture/search", pictureController.SearchPictures()).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/picture/detectText/dropbox", authController.CheckTokenHandler(pictureController.DetectImageFromDropboxUrl())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/picture/{id}", authController.CheckTokenHandler(pictureController.DeletePicture())).Methods(http.MethodDelete, http.MethodOptions)

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
