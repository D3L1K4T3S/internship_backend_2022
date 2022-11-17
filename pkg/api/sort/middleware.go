package sort

import (
	"context"
	"net/http"
	"strings"
)

const (
	ASC               = "ASC"
	DESC              = "DESC"
	OptionsContextKey = "sort_options"
)

func Middleware(handler http.HandlerFunc, defaultSortFierld, defaultSortOrder string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		sortBy := request.URL.Query().Get("sort_by")
		sordOrder := request.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = defaultSortFierld
		}
		if sordOrder == "" {
			sordOrder = defaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sordOrder)
			if upperSortOrder != ASC && upperSortOrder != DESC {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("bad sort order"))
				return
			}
		}

		options := Options{
			Field: sortBy,
			Order: sordOrder,
		}

		ctx := context.WithValue(request.Context(), OptionsContextKey, options)
		request.WithContext(ctx)

		handler(writer, request)
	}
}

type Options struct {
	Field, Order string
}
