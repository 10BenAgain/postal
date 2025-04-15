package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DirectionKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k DirectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit}
}

func (k DirectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}

type EditorKeyMap struct {
	*DirectionKeyMap
	Save    key.Binding
	Commit  key.Binding
	Reset   key.Binding
	Mail    key.Binding
	Growth  key.Binding
	Attacks key.Binding
	EVSC    key.Binding
	Misc    key.Binding
	Stats   key.Binding
	Origin  key.Binding
	Ribbon  key.Binding
	Search  key.Binding
	File    key.Binding
}

func (k EditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit}
}

func (k EditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Mail, k.Growth, k.Attacks, k.EVSC, k.Misc, k.Search, k.File},
		{k.Save, k.Commit, k.Reset},
		{k.Ribbon, k.Origin, k.Stats},
		k.ShortHelp(),
	}
}

type MailKeyMap struct {
	*EditorKeyMap
	Mode key.Binding
	Swap key.Binding
	File key.Binding
}

func (k MailKeyMap) ShortHelp() []key.Binding {
	return k.EditorKeyMap.ShortHelp()
}

func (k MailKeyMap) FullHelp() [][]key.Binding {
	b := k.EditorKeyMap.FullHelp()
	return append(b, []key.Binding{k.Mode, k.Swap, k.File})
}

var (
	SaveKey = tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'S'},
		Alt:   true,
	}

	FileKey = tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'F'},
		Alt:   true,
	}
)

var (
	// RIP Mac Miller
	RibbonViewKey = tea.Key{Type: tea.KeyRunes, Runes: []rune{'r'}, Alt: true}
	OriginViewKey = tea.Key{Type: tea.KeyRunes, Runes: []rune{'o'}, Alt: true}
	StatViewKey   = tea.Key{Type: tea.KeyRunes, Runes: []rune{'s'}, Alt: true}
)

var DirectionKeys = DirectionKeyMap{
	Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp(" 🠕 ", "move up")),
	Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp(" 🠗", "move down")),
	Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp(" 🠔 ", "move left")),
	Right: key.NewBinding(key.WithKeys("right"), key.WithHelp(" 🠖 ", "move right")),
	Help:  key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp(" ctrl+h", "help")),
	Quit:  key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp(" esc", "quit")),
}

var EditorKeys = EditorKeyMap{
	DirectionKeyMap: &DirectionKeys,
	Save:            key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save changes ")),
	Commit:          key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "commit changes ")),
	Reset:           key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "reset entry")),
	Mail:            key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("ctrl+b", "switch to mail mail editor")),
	Growth:          key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("ctrl+g", "switch to growth editor")),
	Attacks:         key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "switch to attack editor")),
	EVSC:            key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "switch to evsc editor")),
	Misc:            key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "switch to misc editor")),
	Stats:           key.NewBinding(key.WithKeys(StatViewKey.String()), key.WithHelp(StatViewKey.String(), "stats view")),
	Origin:          key.NewBinding(key.WithKeys(OriginViewKey.String()), key.WithHelp(OriginViewKey.String(), "origin view")),
	Ribbon:          key.NewBinding(key.WithKeys(RibbonViewKey.String()), key.WithHelp(RibbonViewKey.String(), "ribbon view")),
	Search:          key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("ctrl+f", "switch to search menu")),
	File:            key.NewBinding(key.WithKeys(FileKey.String()), key.WithHelp(FileKey.String(), "open file picker menu")),
}

var MailKeys = MailKeyMap{
	EditorKeyMap: &EditorKeys,
	Mode:         key.NewBinding(key.WithKeys("ctrl+up"), key.WithHelp("ctrl+up", "switch word search mode")),
	Swap:         key.NewBinding(key.WithKeys("ctrl+w"), key.WithHelp("ctrl+w", "swap edit mon with base mon")),
	File:         key.NewBinding(key.WithKeys(SaveKey.String()), key.WithHelp(SaveKey.String(), "save base mon to file")),
}
