// Package ini provides a parser for INI configuration files.
package ini

import (
	"strings"

	"github.com/81120/tiny-parsec/parser"
)

// ISectionName returns a parser that parses the name of a section in an INI file.
// It uses Between to parse the text between square brackets.
func ISectionName() parser.Parser[string] {
	return parser.Between(
		// Parse and trim the opening square bracket
		parser.Trim(parser.Char('[')),
		// Parse zero or more characters that are not closing square brackets
		parser.Bind(
			parser.ZeroOrMore(parser.NotChar(']')),
			func(rs []rune) parser.Parser[string] {
				s := strings.TrimSpace(string(rs))
				if s == "" {
					return parser.Fail[string]()
				} else {
					return parser.Pure(s)
				}
			}),
		// Parse and trim the closing square bracket
		parser.Trim(parser.Char(']')),
	)
}

// ParseINI parses an INI string using the IIni parser.
// It returns the result of the parsing operation.
func ParseINI(str string) parser.ParserFuncRet[Ini] {
	return IniParse().Parse(str)
}

// Ini represents an INI file with a list of sections.
func IniParse() parser.Parser[Ini] {
	return parser.NewParser(func(input string) parser.ParserFuncRet[Ini] {
		strs := strings.Split(input, "\n")
		sections := make([]Section, 0)
		for _, s := range strs {
			s = strings.TrimSpace(s)
			if s == "" || strings.HasPrefix(s, ";") || strings.HasPrefix(s, "#") {
				continue
			}
			r := ISectionName().Parse(s)
			if r.IsJust() {
				section := Section{Name: r.Get().First}
				sections = append(sections, section)
			} else {
				t := strings.Split(s, "=")
				entry := Entry{
					Key:   strings.TrimSpace(t[0]),
					Value: strings.TrimSpace(t[1]),
				}
				newEntries := append(sections[len(sections)-1].Entries, entry)
				sections[len(sections)-1].Entries = newEntries
			}
		}
		return parser.Just(parser.NewTuple(Ini{Sections: sections}, ""))
	})
}
