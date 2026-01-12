package app_context

import (
	"context"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// Ensure AppContext implements TraceContextCarrier interface
var _ contractsobservability.TraceContextCarrier = (*AppContext)(nil)

// AppContext is the context for the application layer
type AppContext struct {
	context.Context
	User         *models.UserWithRole
	OneTimeToken *dtos.OneTimeTokenUser
	trace        *Trace
	traceCtx     contractsobservability.TraceContext
}

// NewContextWithUser creates a new AppContext with a user
func NewContextWithUser(user *models.UserWithRole) *AppContext {
	return &AppContext{
		Context: context.Background(),
		User:    user,
	}
}

// NewVoidAppContext creates a new AppContext with a background context
func NewVoidAppContext() *AppContext {
	return &AppContext{
		Context: context.Background(),
	}
}

// AddUserToContext adds a user to the AppContext
func (a *AppContext) AddUserToContext(user *models.UserWithRole) {
	a.User = user
}

// AddOneTimeTokenToContext adds a one-time token to the AppContext
func (a *AppContext) AddOneTimeTokenToContext(oneTimeToken dtos.OneTimeTokenUser) {
	a.OneTimeToken = &oneTimeToken
}

// AddTraceToContext adds a trace to the AppContext
func (a *AppContext) AddTraceToContext(trace Trace) {
	a.trace = &trace
}

// GetTrace returns the trace from the AppContext
func (a *AppContext) GetTrace() *Trace {
	return a.trace
}

// TraceContext returns the TraceContext from the AppContext
func (a *AppContext) TraceContext() contractsobservability.TraceContext {
	return a.traceCtx
}

// WithTraceContext injects a TraceContext into the AppContext
// Implements TraceContextCarrier interface
func (a *AppContext) WithTraceContext(tc contractsobservability.TraceContext) interface{} {
	a.traceCtx = tc
	return a
}

// HasTrace checks if the AppContext has a valid TraceContext
func (a *AppContext) HasTrace() bool {
	return a.traceCtx != nil && a.traceCtx.IsValid()
}
