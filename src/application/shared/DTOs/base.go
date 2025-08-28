package dtos

type Link struct {
	Self string `json:"self"`
	Next string `json:"next,omitempty"`
	Last string `json:"last,omitempty"`
}

type MetaMultiResponse struct {
	Count    int   `json:"count"`
	Total    int   `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	HasNext  bool  `json:"has_next"`
	HasPrev  bool  `json:"has_prev"`
	Links    *Link `json:"links"`
}

type MultipleResponse[D any] struct {
	Data    []D               `json:"data"`
	Meta    MetaMultiResponse `json:"meta"`
	Details string            `json:"details,omitempty"`
}

type MetaSingleResponse struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Cached    bool   `json:"cached"`
}

type SingleResponse[D any] struct {
	Data    D      `json:"data"`
	Details string `json:"details,omitempty"`
}
