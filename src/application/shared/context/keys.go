package app_context

type AppContextKey string

type AppContextDictionary map[AppContextKey]string

const (
	UserKey AppContextKey = "user"
)
