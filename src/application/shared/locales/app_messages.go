package locales

import (
	messages "gingoskeleton/src/application/shared/locales/messages"
)

type LocaleTypeEnum string

const (
	EN_US LocaleTypeEnum = "en-US"
	ES_ES LocaleTypeEnum = "es-ES"
)

type Locale struct {
	localeType LocaleTypeEnum
}

func NewLocale(localeType LocaleTypeEnum) *Locale {
	return &Locale{
		localeType: localeType,
	}
}

func (l *Locale) Locales() map[LocaleTypeEnum]map[messages.MessageKeysEnum]string {
	locales := make(map[LocaleTypeEnum]map[messages.MessageKeysEnum]string)
	locales[EN_US] = messages.EnMessages
	locales[ES_ES] = messages.EsMessages
	return locales
}

func (l *Locale) Get(localeType LocaleTypeEnum, key messages.MessageKeysEnum) string {
	if localeType == "" {
		localeType = l.localeType
	}

	local, ok := l.Locales()[localeType]
	if !ok {
		return ""
	}
	return local[key]

}
