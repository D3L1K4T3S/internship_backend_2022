package apperror

import "encoding/json"

var (
	ErrorNotFound = NewApplicationError(nil, "not found", "error not found", "1234")
)

type ApplicationError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage"`
	Code             string `json:"code"`
}

func (error *ApplicationError) Error() string {
	return error.Message
}

func (error *ApplicationError) Unwrap() error {
	return error.Err
}

func (error *ApplicationError) Marshal() []byte {
	marshal, err := json.Marshal(error)
	if err != nil {
		return nil
	}
	return marshal
}

func NewApplicationError(err error, message, developerMessage, code string) *ApplicationError {
	return &ApplicationError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *ApplicationError {
	return NewApplicationError(err, "system error", err.Error(), "001")
}
