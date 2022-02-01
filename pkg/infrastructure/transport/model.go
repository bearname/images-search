package transport

type responseWithoutData struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Error struct {
	Status   int
	Response responseWithoutData
}

func NewError(status, code int, err error) *Error {
	return &Error{
		Status: status,
		Response: responseWithoutData{
			Code:    uint32(code),
			Message: err.Error(),
		},
	}
}
