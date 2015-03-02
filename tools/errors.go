
package tools

// APIError is used by writeError to set the status code and the message when a controller returns an APIError
type APIError struct {
  Err   error `json:"-"`
  Code  int `json:"code"`
  Msg   string `json:"msg"`
}

// ErrorResponse is the json response in case of error, Status must be set to 'error'.
//
// The message should be formated in a way it can be embeded in an error message
// destined to then end-user
type ErrorResponse struct {
  Status  string `json:"status"`
  Message string `json:"message"`
}

// NewError creates an ApiError from the code and msg. err can be nil and is for logging purpose
func NewError(err error, code int, msg string) APIError {
  return APIError{
    Err: err,
    Code: code,
    Msg: msg,
  }
}

func (e APIError) Error() string {
  return e.Msg
}

