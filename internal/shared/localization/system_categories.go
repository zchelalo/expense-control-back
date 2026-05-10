package localization

import "strings"

const (
	LanguageSpanish = "es"
	LanguageEnglish = "en"

	SystemCategoryOpeningBalanceKey          = "opening_balance"
	SystemCategoryOpeningBalanceInternalName = "system.opening_balance"
)

func ResolveLanguage(acceptLanguage string) string {
	for _, part := range strings.Split(strings.ToLower(acceptLanguage), ",") {
		tag := strings.TrimSpace(part)
		if tag == "" {
			continue
		}

		tag = strings.Split(tag, ";")[0]

		switch {
		case strings.HasPrefix(tag, LanguageSpanish):
			return LanguageSpanish
		case strings.HasPrefix(tag, LanguageEnglish):
			return LanguageEnglish
		}
	}

	return LanguageSpanish
}

func LocalizeSystemCategoryName(systemKey, language string) (string, bool) {
	switch systemKey {
	case SystemCategoryOpeningBalanceKey:
		if language == LanguageEnglish {
			return "Opening balance", true
		}

		return "Saldo inicial", true
	default:
		return "", false
	}
}
