package tui

import (
	"fmt"
	"postal/pokemon"
	"postal/utils"

	vals "postal/game"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	SearchCat = []string{"Pokemon", "Move", "Item", "Location"}
)

type SearcherModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewSearchModel() *SearcherModel {
	m := SearcherModel{
		inputs: make([]textinput.Model, 4),
	}

	var ti textinput.Model
	for i := range m.inputs {
		ti = textinput.New()
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(LightPink)
		ti.CharLimit = 156
		ti.Width = 12

		switch i {
		case 0:
			ti.Placeholder = "Snorlax"
			ti.Focus()
			ti.PromptStyle = FocusedStyle

		case 1:
			ti.Placeholder = "Sky Attack"
		case 2:
			ti.Placeholder = "Figy Berry"
		case 3:
			ti.Placeholder = "Celedon City"
		}

		m.inputs[i] = ti
	}

	return &m
}

func (m *SearcherModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SearcherModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up", "down", "left", "right", "tab":
			s := msg.String()

			if s == "down" || s == "tab" {
				m.focusIndex++
			} else {
				m.focusIndex--
			}

			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = FocusedStyle
					m.inputs[i].TextStyle = FocusedText
					continue
				}

				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = lipgloss.NewStyle()
				m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(SubText)
			}
			return m, tea.Batch(cmds...)
		}
	}

	return m, m.updateInputs(msg)
}

func (m *SearcherModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *SearcherModel) View() string {
	var result string
	out := []string{}

	for i := range m.inputs {
		out = append(
			out,
			fmt.Sprintf("%s %s",
				CatStyle.Render(SearchCat[i]),
				m.inputs[i].View(),
			),
		)
	}

	switch m.focusIndex {
	case 0:
		mon, err := vals.MonLookup(utils.SanitizeSearch(m.inputs[0].Value()))
		if err == nil {
			result = makeSearchResultString(0, uint16(mon))
		}
	case 1:
		item, err := vals.ItemLookup(utils.SanitizeSearch(m.inputs[1].Value()))
		if err == nil {
			result = makeSearchResultString(1, uint16(item))
		}
	case 2:
		move, err := vals.MoveLookup(utils.SanitizeSearch(m.inputs[2].Value()))
		if err == nil {
			result = makeSearchResultString(2, uint16(move))
		}
	case 3:
		loc, err := vals.MetLocLookup(utils.SanitizeSearch(m.inputs[3].Value()))
		if err == nil {
			result = fmt.Sprintf("0x%02X / %02d ", loc, loc)
		}
	}

	so := MailMenuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, out...))
	final := lipgloss.JoinVertical(
		lipgloss.Left,
		so,
	)

	if len(result) > 0 {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			final,
			" ",
			MailMenuStyle.Render(result),
		)
	}

	return final
}

func makeSearchResultString(t int, n uint16) string {
	// 0 Mon, 1 Item, 2 Move
	fields := []string{blank, blank, blank}

	species, err := pokemon.NumToSpeciesName(n)
	if err == nil {
		fields[0] = species
	}

	move, err := pokemon.NumToMoveName(n)
	if err == nil {
		fields[1] = move
	}

	item, err := pokemon.NumToItemName(n)
	if err == nil {
		fields[2] = item
	}

	f, s := pokemon.NumToEVVals(n)

	strs := map[int]string{
		0: SearchHeaderStyle.Render(fmt.Sprintf("%d | %04X | 0x%02X 0x%02x\n", n, n, f, s)),
		1: joinSearchValueStrings(SearchCat[0], fields[0]), // Pokemon
		2: joinSearchValueStrings(SearchCat[1], fields[1]), // Move
		3: joinSearchValueStrings(SearchCat[2], fields[2]), // Item
	}

	m := make([]string, 5)
	m[0] = strs[0]

	switch t {
	case 0:
		m[1] = strs[2]
		m[2] = strs[3]

	case 1:
		m[1] = strs[1]
		m[2] = strs[2]

	case 2:
		m[1] = strs[1]
		m[2] = strs[3]
	}

	m[3] = joinSearchValueStrings("EV 1:", fmt.Sprintf("%02d", f))
	m[4] = joinSearchValueStrings("EV 2:", fmt.Sprintf("%02d", s))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m...,
	)
}
