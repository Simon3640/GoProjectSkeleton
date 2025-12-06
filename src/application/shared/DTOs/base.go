package dtos

import (
	"strconv"
)

type Link struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
	Last *string `json:"last"`
}

type MetaMultiResponse struct {
	Count   int   `json:"count"`
	Total   int64 `json:"total"`
	HasNext bool  `json:"hasNext"`
	HasPrev bool  `json:"hasPrev"`
	Links   *Link `json:"links"`
	Cached  bool  `json:"cached"`
}

func NewMetaMultiResponse(count int, total int64, hasNext bool, hasPrev bool, cached bool) MetaMultiResponse {
	return MetaMultiResponse{
		Count:   count,
		Total:   total,
		HasNext: hasNext,
		HasPrev: hasPrev,
		Cached:  cached,
	}
}

func (mr *MetaMultiResponse) BuildLinks(prefix string, page int, pageSize int, queryParamsUrl string) {
	if mr.Links == nil {
		mr.Links = &Link{}
	}

	mr.Links.Self = prefix + "?" + queryParamsUrl
	if mr.HasNext {
		mr.Links.Next = new(string)
		*mr.Links.Next = (prefix +
			"?page=" +
			strconv.Itoa(page+1) +
			"&page_size=" +
			strconv.Itoa(pageSize) +
			"&" + queryParamsUrl)
	}
	if mr.HasPrev {
		mr.Links.Last = new(string)
		*mr.Links.Last = prefix + "?page=" + strconv.Itoa(page-1) + "&page_size=" + strconv.Itoa(pageSize)
	}
}

type MultipleResponse[D any] struct {
	Records []D               `json:"records"`
	Meta    MetaMultiResponse `json:"meta"`
}

type MetaSingleResponse struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Cached    bool   `json:"cached"`
}

type SingleResponse[D any] struct {
	Record D                  `json:"record"`
	Meta   MetaSingleResponse `json:"meta"`
}
