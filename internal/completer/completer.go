package completer

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func FromShortcuts(shc map[string]string) prompt.Completer {
	s := make([]prompt.Suggest, len(shc))
	i := 0
	for k, v := range shc {
		s[i].Text = k
		s[i].Description = v
		i = i + 1
	}

	return func(d prompt.Document) []prompt.Suggest {
		if d.TextBeforeCursor() == "" {
			return []prompt.Suggest{}
		}

		args := strings.Split(d.TextBeforeCursor(), " ")
		if len(args) == 1 {
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}
		return []prompt.Suggest{}
	}
}
