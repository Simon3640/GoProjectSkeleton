package handlers

import (
	"io"
	"net/http"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
)

// HTTPHeaderTypeEnum is the type for the HTTP header type
type HTTPHeaderTypeEnum string

func (h HTTPHeaderTypeEnum) String() string {
	return string(h)
}

// HTTPHeaderTypeEnum constants
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

// SerializationTypeEnum is the type for the serialization type
type SerializationTypeEnum string

func (s SerializationTypeEnum) String() string {
	return string(s)
}

// SerializationTypeEnum constants
const (
	JSON         SerializationTypeEnum = "json"
	XML          SerializationTypeEnum = "xml"
	TEXT         SerializationTypeEnum = "text"
	ARRAY_BUFFER SerializationTypeEnum = "arrayBuffer"
	BLOB         SerializationTypeEnum = "blob"
	FORM_DATA    SerializationTypeEnum = "formData"
)

// HTTPContentTypeEnum is the type for the HTTP content type
type HTTPContentTypeEnum string

// HTTPContentTypeEnum constants
const (
	APPLICATION_JSON                  HTTPContentTypeEnum = "application/json"
	APPLICATION_XML                   HTTPContentTypeEnum = "application/xml"
	APPLICATION_X_WWW_FORM_URLENCODED HTTPContentTypeEnum = "application/x-www-form-urlencoded"
	MULTIPART_FORM_DATA               HTTPContentTypeEnum = "multipart/form-data"
	TEXT_HTML                         HTTPContentTypeEnum = "text/html"
	TEXT_PLAIN                        HTTPContentTypeEnum = "text/plain"
	TEXT_XML                          HTTPContentTypeEnum = "text/xml"
)

// HTTPMethodEnum is the type for the HTTP method
type HTTPMethodEnum string

func (m HTTPMethodEnum) String() string {
	return string(m)
}

// HTTPMethodEnum constants

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

// Query is the type for the query
type Query struct {
	Filters  []string
	Sorts    []string
	Page     *int
	PageSize *int
}

// HandlerContext is the context for the handler
type HandlerContext struct {
	Context *app_context.AppContext

	Locale locales.LocaleTypeEnum
	Params map[string]string
	Body   *io.ReadCloser
	Query  *Query

	ResponseWriter http.ResponseWriter
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(
	c *app_context.AppContext,
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
		Context:        c,
		Locale:         locales.LocaleTypeEnum(*locale),
		Params:         params,
		Body:           body,
		Query:          query,
		ResponseWriter: responseWriter,
	}
}
