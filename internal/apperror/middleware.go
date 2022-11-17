package apperror

import (
	"errors"
	"net/http"
)

type applicationHandler func(writer http.ResponseWriter, request *http.Request) error

func Middleware(handler applicationHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		var applicationError *ApplicationError
		err := handler(writer, request)

		if err != nil {
			if errors.As(err, &applicationError) {
				if errors.Is(err, ErrorNotFound) {
					writer.WriteHeader(http.StatusNotFound)
					_, _ = writer.Write(ErrorNotFound.Marshal())
					return
				}
				//Остальные перечисленные ошибки

				err = err.(*ApplicationError)
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write(applicationError.Marshal())
				return
			}

			writer.WriteHeader(http.StatusTeapot)
			_, _ = writer.Write(systemError(err).Marshal())
		}
	}
}
