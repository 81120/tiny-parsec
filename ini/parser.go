// Package ini provides a parser for INI configuration files.
package ini

import (
	"strings"

	"github.com/81120/tiny-parsec/parser"
)

// IIni returns a parser that parses an entire INI file.
// It uses ZeroOrMore to parse zero or more sections and constructs an Ini struct.
func IIni() parser.Parser[Ini] {
	return parser.Fmap(
		// Parse zero or more sections
		parser.ZeroOrMore(ISection()),
		// Convert the parsed sections into an Ini struct
		func(sections []Section) Ini {
			return Ini{Sections: sections}
		})
}

// ISectionName returns a parser that parses the name of a section in an INI file.
// It uses Between to parse the text between square brackets.
func ISectionName() parser.Parser[string] {
	return parser.Between(
		// Parse and trim the opening square bracket
		parser.Trim(parser.Char('[')),
		// Parse zero or more characters that are not closing square brackets
		parser.ToString(
			parser.ZeroOrMore(parser.NotChar(']')),
			true,
		),
		// Parse and trim the closing square bracket
		parser.Trim(parser.Char(']')),
	)
}

// ISection returns a parser that parses a section in an INI file.
// It first parses the section name and then parses zero or more entries.
func ISection() parser.Parser[Section] {
	return parser.Bind(
		// Parse the section name
		ISectionName(),
		func(name string) parser.Parser[Section] {
			return parser.Fmap(
				// Parse zero or more entries separated by newlines
				parser.SepBy(
					IEntry(),
					parser.Char('\n')),
				// Convert the parsed entries into a Section struct
				func(entries []Entry) Section {
					return Section{Name: name, Entries: entries}
				})
		})
}

// IEntry returns a parser that parses an entry in an INI file.
// It parses a key-value pair separated by an equal sign.
func IEntry() parser.Parser[Entry] {
	return parser.Fmap(
		// Parse the key, equal sign, and value
		parser.Seq(
			parser.ToString(parser.OneOrMore(parser.NotChar('=')), true),
			parser.ToString(parser.Char('='), true),
			parser.ToString(parser.OneOrMore(parser.NotChar('\n')), false),
		),
		// Convert the parsed strings into an Entry struct
		func(strs []string) Entry {
			return Entry{Key: strings.TrimSpace(strs[0]), Value: strings.TrimSpace(strs[2])}
		},
	)
}

// ParseINI parses an INI string using the IIni parser.
// It returns the result of the parsing operation.
func ParseINI(str string) parser.ParserFuncRet[Ini] {
	return IniParse().Parse(str)
}

func IniParse() parser.Parser[Ini] {
	return parser.NewParser(func(input string) parser.ParserFuncRet[Ini] {
		strs := strings.Split(input, "\n")
		sections := make([]Section, 0)

		for _, s := range strs {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}

			r := ISectionName().Parse(s)
			if r.IsJust() {
				section := Section{Name: r.Get().First}
				sections = append(sections, section)
			} else {
				t := IEntry().Parse(s)
				if t.IsJust() {
					newEntries := append(sections[len(sections)-1].Entries, t.Get().First)
					sections[len(sections)-1].Entries = newEntries
				} else {
					break
				}
			}
		}

		return parser.Just(parser.NewTuple(Ini{Sections: sections}, ""))
	})
}
