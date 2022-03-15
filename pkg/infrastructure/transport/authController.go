package transport

import (
	"encoding/json"
	"errors"
	"github.com/col3name/images-search/pkg/domain/auth"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func (c *AuthController) Login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*req).Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var userDto auth.Credentials
	err := json.NewDecoder(req.Body).Decode(&userDto)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(&userDto)
	if err != nil {
		c.WriteError(w, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(w, token)
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*req).Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	authDto, err := c.decodeRefreshTokenReq(req)
	if err != nil {
		msg := err.Error()
		log.Println(err, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	accessTokenResponse, err := c.authService.RefreshToken(authDto)

	if err != nil {
		c.BaseController.WriteError(w, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(w, accessTokenResponse)
}

func (c *AuthController) CheckTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		header := req.Header.Get("Authorization")
		username, err := c.authService.ValidateToken(header)
		if err != nil {
			c.BaseController.WriteError(w, err, TranslateError(err))
			return
		}

		context.Set(req, "username", username)

		next(w, req)
	}
}

func (c *AuthController) ValidateToken(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*req).Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	header := req.Header.Get("Authorization")
	_, err := c.authService.ValidateToken(header)
	if err != nil {
		c.BaseController.WriteError(w, err, TranslateError(err))
		return
	}

	c.BaseController.WriteJsonResponse(w, struct{}{})
}

func (c *AuthController) decodeRefreshTokenReq(req *http.Request) (*auth.RefreshTokenDto, error) {
	username, ok := context.Get(req, "username").(string)
	if !ok {
		context.Clear(req)
		return nil, errors.New("cannot check username")

	}
	accessToken, ok := context.Get(req, "accessToken").(string)
	if !ok {
		context.Clear(req)
		return nil, errors.New("accessToken to preset on context by accessToken checker")
	}
	context.Clear(req)
	return &auth.RefreshTokenDto{Username: username, Token: accessToken}, nil
}
