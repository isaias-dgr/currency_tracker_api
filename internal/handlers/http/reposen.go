package http

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
)

type Response struct {
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func NewResponse(data interface{}, total int, filter *domain.Filter, msg string) Response {
	resp := Response{
		Data: data,
	}
	if filter != nil {
		resp.Metadata = NewMetadata(total, filter, msg)
	}
	return resp
}

type Metadata struct {
	Offset  int    `json:"offset,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Total   int    `json:"total,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewMetadata(total int, filter *domain.Filter, msg string) *Metadata {
	return &Metadata{
		Offset:  filter.Offset,
		Limit:   filter.Limit,
		Total:   total,
		Message: msg,
	}
}

func NewFilter(qs url.Values) *domain.Filter {
	return &domain.Filter{
		Offset: GetIntDefault(qs, "offset", 0),
		Limit:  GetIntDefault(qs, "limit", 100),
		SortBy: GetDefault(qs, "sort_by", "created_at"),
		Finit:  GetDateDefault(qs, "finit", ""),
		Fend:   GetDateDefault(qs, "fend", ""),
	}
}

func GetIntDefault(qs url.Values, k string, v int) int {
	val, err := strconv.Atoi(GetDefault(qs, k, strconv.Itoa(v)))
	if err != nil {
		log.Printf("> error %s", err)
	}
	return val
}

func GetDateDefault(qs url.Values, k string, v string) *time.Time {
	val, err := time.Parse("2006-01-02 15:04:05", GetDefault(qs, k, v))
	if err != nil {
		log.Printf("> %s %s error %s", k, v, err)
		return nil
	}
	return &val
}

func GetDefault(qs url.Values, k string, v string) string {
	val := qs.Get(k)
	if val == "" {
		return v
	}
	return val
}
