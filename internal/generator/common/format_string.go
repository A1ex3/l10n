package common

import (
	"strings"

	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

func EscapeSymbolsString(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))

	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\\':
			builder.WriteString(`\\`)
		case '"':
			builder.WriteString(`\"`)
		case '\n':
			builder.WriteString(`\n`)
		case '\r':
			builder.WriteString(`\r`)
		case '\t':
			builder.WriteString(`\t`)
		default:
			builder.WriteByte(value[i])
		}
	}

	return builder.String()
}

func CapitalizeAfterHyphen(input string) string {
	parts := strings.Split(input, "-")
	var capitalizedParts []string

	for _, part := range parts {
		if part != "" {
			capitalized := cases.Title(lang.English).String(part)
			capitalizedParts = append(capitalizedParts, capitalized)
		}
	}
	return strings.Join(capitalizedParts, "")
}
