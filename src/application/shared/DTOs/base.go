package dtos

type Link struct {
	Self string `json:"self"`
	Next string `json:"next,omitempty"`
	Last string `json:"last,omitempty"`
}

type MetaMultiResponse struct {
	Count   int   `json:"count"`
	Total   int   `json:"total"`
	HasNext bool  `json:"has_next"`
	HasPrev bool  `json:"has_prev"`
	Links   *Link `json:"links"`
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
