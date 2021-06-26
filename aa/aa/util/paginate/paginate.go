package paginate

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const DEFAULT_LIMIT = 30

type PaginationReq struct {
	Page int `json:"page"`
}

func ParsePaginationReq(r *http.Request) *PaginationReq {
	var pReq PaginationReq
	p := r.URL.Query().Get("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			log.Println("Error while parsing pagination: ", err.Error())
		}

		pReq.Page = page
	}

	return &pReq
}

func Paginate(query string, page int, limit ...int) string {
	l := DEFAULT_LIMIT
	p := page - 1

	if len(limit) > 0 {
		l = limit[0]
	}

	if page <= 0 {
		return query
	}

	return fmt.Sprintf("%s LIMIT %d, %d", query, p*l, l)
}
