package utils

import (
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

// duplicated from https://github.com/mvdan/xurls/blob/master/xurls.go#L83-L88 as xurls doesn't let the user extend or
// override SchemesNoAuthority
var (
	knownSchemes = []string{
		`cid`,
		`file`,
		`magnet`,
		`mailto`,
		`mid`,
		`sms`,
		`tel`,
		`xmpp`,
		`spotify`,
		`facetime`,
		`facetime-audio`,
	}
	anyScheme = `(?:[a-zA-Z][a-zA-Z.\-+]*://|` + anyOf(knownSchemes...) + `:)`
)

func anyOf(strs ...string) string {
	var b strings.Builder
	b.WriteString("(?:")
	for i, s := range strs {
		if i != 0 {
			b.WriteByte('|')
		}
		b.WriteString(regexp.QuoteMeta(s))
	}
	b.WriteByte(')')
	return b.String()
}

func GetURIRegex() *regexp.Regexp {
	rgx, err := xurls.StrictMatchingScheme(anyScheme)
	if err != nil {
		return xurls.Strict()
	}

	return rgx
}

func ExtractURIs(rg *regexp.Regexp, text string) []string {
	return rg.FindAllString(text, -1)
}
