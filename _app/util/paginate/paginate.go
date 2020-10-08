package paginate

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

const DEFAULT_LIMIT = 30
const MIN_PAGE = 1

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

func Paginate(db *gorm.DB, page int, limit ...int) *gorm.DB {
	l := DEFAULT_LIMIT
	p := page - 1

	if len(limit) > 0 {
		l = limit[0]
	}

	if page < MIN_PAGE {
		p = MIN_PAGE - 1
	}

	return db.Offset(p * l).Limit(l)
}

func FilterPaginate(db *gorm.DB, page int, filter map[string]string) *gorm.DB {
	db = Paginate(db, page)

	for k, v := range filter {
		db.Where(k+"= ?", v)
	}

	return db
}
