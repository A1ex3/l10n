package common

import "strings"

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
