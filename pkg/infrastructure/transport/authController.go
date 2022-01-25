package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"photofinish/pkg/domain/auth"
	"strings"
)

type AuthController struct {
	BaseController
	authService auth.Service
}

func NewAuthController(authService auth.Service) *AuthController {
	v := new(AuthController)
	v.authService = authService
	return v
}

func (c *AuthController) Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	var userDto auth.Credentials
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		log.Error(err.Error())
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(userDto)
	if err != nil {
		c.WriteError(writer, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(writer, token)
}

func (c *AuthController) RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	username, ok := context.Get(request, "username").(string)
	if !ok {
		context.Clear(request)
		http.Error(writer, "cannot check username", http.StatusBadRequest)
		return
	}
	accessToken, ok := context.Get(request, "accessToken").(string)
	if !ok {
		context.Clear(request)
		http.Error(writer, "accessToken to preset on context by accessToken checker", http.StatusBadRequest)
		return
	}
	context.Clear(request)

	accessTokenResponse, err := c.authService.RefreshToken(auth.RefreshTokenDto{Username: username, Token: accessToken})

	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(writer, accessTokenResponse)
}

func (c *AuthController) CheckTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fmt.Println(req.URL.String())
		if strings.Contains(req.URL.String(), "dropbox") {
			fmt.Println(req.URL)

		}
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		log.Println("check token handler")

		header := req.Header.Get("Authorization")
		username, err := c.authService.ValidateToken(header)
		if err != nil {
			c.BaseController.WriteError(w, err, TranslateError(err))
			return
		}

		context.Set(req, "username", username)
		log.Println("success")

		next(w, req)
	}
}

func (c *AuthController) ValidateToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("check token handler")

	header := request.Header.Get("Authorization")
	_, err := c.authService.ValidateToken(header)
	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.BaseController.WriteJsonResponse(writer, struct {
	}{})
}
