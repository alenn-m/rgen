package resp

import (
	"encoding/json"
	"math"
	"net/http"

	"{{Root}}/util/paginate"
)

type SuccessResp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type PaginatedResp struct {
	Code       int         `json:"code"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
}

func ReturnPaginatedSuccess(w http.ResponseWriter, page, total int, data interface{}) {
	totalPages := math.Ceil(float64(total) / float64(paginate.DEFAULT_LIMIT))
	if page < 1 {
		page = 1
	}

	if totalPages < 1 {
		totalPages = 1
	}

	r := PaginatedResp{
		Code:       http.StatusOK,
		Data:       data,
		Page:       page,
		TotalItems: total,
		TotalPages: int(totalPages),
	}

	result, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func ReturnSuccess(w http.ResponseWriter, data interface{}) {
	r := SuccessResp{
		Code: http.StatusOK,
		Data: data,
	}

	result, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

type ErrorResp struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func ReturnError(w http.ResponseWriter, message string, code int) {
	r := ErrorResp{
		Code:  code,
		Error: message,
	}

	result, _ := json.Marshal(r)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
