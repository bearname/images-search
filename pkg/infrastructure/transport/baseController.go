package transport

import (
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrBadRouting = errors.New("bad routing")
	ErrBadRequest = errors.New("bad request")
)

type BaseController struct {
}

func (c *BaseController) WriteError(w http.ResponseWriter, err error, responseError Error) {
	log.Error(err.Error())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(responseError.Status)
	_ = json.NewEncoder(w).Encode(responseError.Response)
}

func (c *BaseController) SetupCors(w *http.ResponseWriter, _ *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (c *BaseController) WriteJsonResponse(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(jsonData)
	if err != nil {
		return
	}
}

//type transportError struct {
//	Status   int
//	Response responseWithoutData
//}
//
//func (c *BaseController) translateError(err error) transportError {
//	if errors.Is(err, ErrBadRequest) {
//		return transportError{
//			Status: http.StatusBadRequest,
//			Response: responseWithoutData{
//				Code:    101,
//				Message: err.Error(),
//			},
//		}
//		//} else if errors.Is(err, app.ErrUrlNotFound) {
//		//    return transportError{
//		//        Status: http.StatusNotFound,
//		//        Response: responseWithoutData{
//		//            Code:    102,
//		//            Message: err.Error(),
//		//        },
//		//    }
//		//} else if errors.Is(err, app.ErrDuplicateUrl) {
//		//    return transportError{
//		//        Status: http.StatusConflict,
//		//        Response: responseWithoutData{
//		//            Code:    103,
//		//            Message: err.Error(),
//		//        },
//		//    }
//		//} else if errors.Is(err, app.ErrInvalidUrl) {
//		//    return transportError{
//		//        Status: http.StatusBadRequest,
//		//        Response: responseWithoutData{
//		//            Code:    104,
//		//            Message: err.Error(),
//		//        },
//		//    }
//	}
//
//	return transportError{
//		Status: http.StatusInternalServerError,
//		Response: responseWithoutData{
//			Code:    100,
//			Message: "unexpected error",
//		},
//	}
//}
