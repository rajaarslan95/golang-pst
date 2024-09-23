package schemas

type Response struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
}

func NewError(message string, err int16) *Response {
	return &Response{
		Status:  err,
		Message: message,
	}
}
