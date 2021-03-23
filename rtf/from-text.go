package rtf

import "strings"

// FromPlainText produces RTF-formatted text from plain text file
func FromPlainText(plain string, emptyLineBetweenParagraphs bool) string {
	if plain == "" {
		return ""
	}

	escaped := strings.ReplaceAll(plain, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `{`, `\{`)
	escaped = strings.ReplaceAll(escaped, `}`, `\}`)
	if emptyLineBetweenParagraphs {
		escaped = strings.ReplaceAll(escaped, "\n", "\\par\\par\r\n")
	} else {
		escaped = strings.ReplaceAll(escaped, "\n", "\\par\r\n")
	}

	return `{\rtf1\ansi\ansicpg1250\deff0{\fonttbl\f0\fswiss Helvetica;}\f0\fs18\pard ` + escaped + " }"
}
