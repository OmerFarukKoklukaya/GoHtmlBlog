package models

import (
	"strconv"
)

type ViewData struct {
	Data any
	Meta Meta
}

type Meta struct {
	PageTitle    string
	Pagination   Pagination
	AuthedUser   User
	BaseURL      string
	IsAuthorized bool
	IsAdmin      bool
}

type Pagination struct {
	Page     int
	Limit    int
	NextPage string
	PrevPage string
	LastPage int
}

func (Pagination) Paginate(page int, dataLength int, limit int) Pagination {
	var nextPage = ""
	var prevPage = ""
	var lastPage = 0
	if dataLength%limit == 0 {
		lastPage = dataLength / limit
	} else {
		lastPage = dataLength/limit + 1
	}

	if lastPage > page {
		nextPage = "?limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page+1)
	}
	if page != 1 {
		prevPage = "?limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page-1)
	}

	var paging = Pagination{
		Page:     page,
		Limit:    limit,
		NextPage: nextPage,
		PrevPage: prevPage,
		LastPage: lastPage,
	}

	return paging
}
