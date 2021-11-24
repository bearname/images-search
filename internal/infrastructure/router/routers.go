package router

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"photofinish/internal/infrastructure/transport"
)

func Router(pictureController *transport.PictureController, eventsController *transport.EventController, authController *transport.AuthController) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", readyCheckHandler).Methods(http.MethodGet)

	apiV1Route := router.PathPrefix("/api/v1").Subrouter()
	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.List())).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/events", authController.CheckTokenHandler(eventsController.CreateEvent())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/events/{id}", authController.CheckTokenHandler(eventsController.DeleteEvent())).Methods(http.MethodDelete, http.MethodOptions)

	apiV1Route.HandleFunc("/picture/search", pictureController.SearchPictures()).Methods(http.MethodGet, http.MethodOptions)
	//apiV1Route.HandleFunc("/picture/detectText", pictureController.DetectImageFromArchive()).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/picture/detectText/dropbox", authController.CheckTokenHandler(pictureController.DetectImageFromDropboxUrl())).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/picture/{id}", authController.CheckTokenHandler(pictureController.DeletePicture())).Methods(http.MethodDelete, http.MethodOptions)

	apiV1Route.HandleFunc("/auth/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.RefreshToken)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token/validate", authController.ValidateToken).Methods(http.MethodGet, http.MethodOptions)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	//spa := spaHandler{staticPath: "./html", indexPath: "index.html"}
	//routers.PathPrefix("/").Handler(spa)
	//apiV1Route.HandleFunc("/playlists", middleware.AllowCors(playListController.GetUserPlaylists())).Methods(http.MethodGet, http.MethodOptions)
	//apiV1Route.HandleFunc("/playlists/{playlistId}", middleware.AllowCors(playListController.GetPlayList())).Methods(http.MethodGet, http.MethodOptions)
	//apiV1Route.HandleFunc("/playlists/{playlistId}/modify", middleware.AuthMiddleware(playListController.ModifyVideoToPlaylist(), authServerAddress)).Methods(http.MethodPut, http.MethodOptions)
	//apiV1Route.HandleFunc("/playlists/{playlistId}/change-privacy/{privacyType}", middleware.AllowCors(middleware.AuthMiddleware(playListController.ChangePrivacy(), authServerAddress))).Methods(http.MethodPut, http.MethodOptions)
	//apiV1Route.HandleFunc("/playlists/{playlistId}", middleware.AllowCors(middleware.AuthMiddleware(playListController.DeletePlaylist(), authServerAddress))).Methods(http.MethodDelete, http.MethodOptions)
	//
	//apiV1Route.HandleFunc("/videos/", videoController.GetVideos()).Methods(http.MethodGet)
	//apiV1Route.HandleFunc("/videos/search", videoController.SearchVideo()).Methods(http.MethodGet)
	//apiV1Route.HandleFunc("/videos/{videoId}", videoController.GetVideo()).Methods(http.MethodGet)
	//apiV1Route.HandleFunc("/videos/", videoController.UploadVideo()).Methods(http.MethodPost, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos/{videoId}", videoController.UpdateTitleAndDescription()).Methods(http.MethodPut, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos/{videoId}/add-quality", videoController.AddQuality()).Methods(http.MethodPut, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos/{videoId}", videoController.DeleteVideo()).Methods(http.MethodDelete, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos/{videoId}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos-liked", middleware.AuthMiddleware(videoController.FindUserLikedVideo(), authServerAddress)).Methods(http.MethodGet, http.MethodOptions)
	//apiV1Route.HandleFunc("/videos/{videoId}/like/{isLike:[0-1]}", middleware.AuthMiddleware(videoController.LikeVideo(), authServerAddress)).Methods(http.MethodPost, http.MethodOptions)

	return logMiddleware(router)
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	fmt.Println(path)
	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
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
