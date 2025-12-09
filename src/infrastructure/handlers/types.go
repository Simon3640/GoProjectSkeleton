package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
)

type HTTPHeaderTypeEnum string

func (h HTTPHeaderTypeEnum) String() string {
	return string(h)
}

const (
	PROXY_AUTHORIZATION HTTPHeaderTypeEnum = "proxy-authorization"
	TRANSACTION_ID      HTTPHeaderTypeEnum = "x-transaction-id"
	ACCEPT_LANGUAGE     HTTPHeaderTypeEnum = "accept-language"
	FORWARDED_FOR       HTTPHeaderTypeEnum = "x-forwarded-for"
	USER_AGENT          HTTPHeaderTypeEnum = "user-agent"
	ORIGIN              HTTPHeaderTypeEnum = "origin"
	REFERRER            HTTPHeaderTypeEnum = "referrer"
	CONTENT_TYPE        HTTPHeaderTypeEnum = "content-type"
	AUTHORIZATION       HTTPHeaderTypeEnum = "authorization"
)

type SerializationTypeEnum string

const (
	JSON         SerializationTypeEnum = "json"
	XML          SerializationTypeEnum = "xml"
	TEXT         SerializationTypeEnum = "text"
	ARRAY_BUFFER SerializationTypeEnum = "arrayBuffer"
	BLOB         SerializationTypeEnum = "blob"
	FORM_DATA    SerializationTypeEnum = "formData"
)

type HTTPContentTypeEnum string

const (
	APPLICATION_JSON                  HTTPContentTypeEnum = "application/json"
	APPLICATION_XML                   HTTPContentTypeEnum = "application/xml"
	APPLICATION_X_WWW_FORM_URLENCODED HTTPContentTypeEnum = "application/x-www-form-urlencoded"
	MULTIPART_FORM_DATA               HTTPContentTypeEnum = "multipart/form-data"
	TEXT_HTML                         HTTPContentTypeEnum = "text/html"
	TEXT_PLAIN                        HTTPContentTypeEnum = "text/plain"
	TEXT_XML                          HTTPContentTypeEnum = "text/xml"
)

type HTTPMethodEnum string

const (
	GET     HTTPMethodEnum = "GET"
	POST    HTTPMethodEnum = "POST"
	PUT     HTTPMethodEnum = "PUT"
	DELETE  HTTPMethodEnum = "DELETE"
	PATCH   HTTPMethodEnum = "PATCH"
	OPTIONS HTTPMethodEnum = "OPTIONS"
	HEAD    HTTPMethodEnum = "HEAD"
	CONNECT HTTPMethodEnum = "CONNECT"
	TRACE   HTTPMethodEnum = "TRACE"
)

type Query struct {
	Filters  []string
	Sorts    []string
	Page     *int
	PageSize *int
}

type HandlerContext struct {
	c      context.Context
	Locale locales.LocaleTypeEnum
	Params map[string]string
	Body   *io.ReadCloser
	Query  *Query

	ResponseWriter http.ResponseWriter
}

func NewHandlerContext(
	c context.Context,
	locale *string,
	params map[string]string,
	body *io.ReadCloser,
	query *Query,
	responseWriter http.ResponseWriter,
) HandlerContext {
	if locale == nil || *locale == "" {
		defaultLocale := "en-US"
		locale = &defaultLocale
	}

	return HandlerContext{
		c:              c,
		Locale:         locales.LocaleTypeEnum(*locale),
		Params:         params,
		Body:           body,
		Query:          query,
		ResponseWriter: responseWriter,
	}
}
