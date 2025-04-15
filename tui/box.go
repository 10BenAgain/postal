package tui

import (
	"fmt"
	"postal/boxes"
	"postal/pokemon"
	"postal/save"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type returnMsg struct {
	pk pokemon.PStructure
}

func returnToMain(p pokemon.PStructure) tea.Cmd {
	return func() tea.Msg { return returnMsg{pk: p} }
}

type boxSelect struct {
	index int

	table table.Model

	data save.RawBoxDataTotal

	mash boxes.PCBoxBufferMash

	box boxes.RWPCBox

	pks pokemon.PStructure

	name string
}

var cols = []table.Column{
	{Title: "Slot", Width: 5},
	{Title: "Species", Width: 10},
	{Title: "Item", Width: 10},
	{Title: "PID", Width: 10},
	{Title: "Lvl", Width: 4},
	{Title: "Nature", Width: 10},
	{Title: "Ability", Width: 16},
}

func (m boxSelect) Init() tea.Cmd { return nil }

func (m boxSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "tab", "shift+tab":
			s := msg.String()

			switch s {
			case "tab":
				if m.index == 13 {
					m.index = 0
				} else {
					m.index++
				}
			case "shift+tab":
				if m.index == 0 {
					m.index = 13
				} else {
					m.index--
				}
			}

			m.UpdateBox()
			m.UpdateTable()

		case "enter":
			if len(m.table.SelectedRow()) > 0 {
				i, err := strconv.Atoi(m.table.SelectedRow()[0])
				if err != nil {
					return m, cmd
				} else {
					m.pks = m.box.Mons[i-1]
					return m, returnToMain(m.pks)
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m boxSelect) View() string {
	s := fmt.Sprintf(" Box #%d %s", m.index+1, m.name)

	// TODO: Fix this to not regen all the data on every update
	// This will probably require a lot of fixing of the box read
	// Code since its super jank currently
	// This is pretty lazy but it works for now

	// Holy fuck this is bad
	if len(m.table.SelectedRow()) > 0 {
		i, err := strconv.Atoi(m.table.SelectedRow()[0])
		if err != nil {
			goto end
		} else {
			p := generateMonViewOrder(&m.box.Mons[i-1])
			o := lipgloss.JoinVertical(lipgloss.Left, s, baseStyle.Render(m.table.View()))
			return lipgloss.JoinHorizontal(lipgloss.Center, o, p)
		}
	}

end:
	return lipgloss.JoinVertical(lipgloss.Left, s, baseStyle.Render(m.table.View()))
}

func (m *boxSelect) GetPokemon() pokemon.PStructure { return m.pks }

func (m *boxSelect) UpdateBox() {
	m.box = boxes.GenerateRWPCBox(
		m.index, m.data, m.mash.GenerateBoxMonData(m.index),
	)
}

func (m *boxSelect) UpdateTable() {
	m.name = m.box.Name
	m.table = makeTableFromBox(m.box)
}

func NewBoxSelect(data save.RawBoxDataTotal) boxSelect {
	// User box one for default
	buffer := boxes.GeneratePCBoxBufferMash(data)
	def := boxes.GenerateRWPCBox(0, data, buffer.GenerateBoxMonData(0))

	return boxSelect{
		index: 0,
		table: makeTableFromBox(def),
		data:  data,
		mash:  boxes.GeneratePCBoxBufferMash(data),
		box:   def,
		name:  def.Name,
	}
}

func monToTableRow(s string, p pokemon.PStructure) table.Row {
	_, species := p.GetSpecies()
	_, nature := p.GetNature()
	_, item := p.GetHeldItem()
	_, ab := p.GetAbility()
	pid := fmt.Sprintf("%08X", p.PID)
	lvl := fmt.Sprintf("%d", p.Level)

	return table.Row{s, species, item, pid, lvl, nature, ab}
}

func makeTableFromBox(b boxes.RWPCBox) table.Model {
	rows := []table.Row{}

	for i := range b.Mons {
		if !b.Mons[i].IsBlank() {
			rows = append(rows, monToTableRow(fmt.Sprintf("%d", i+1), b.Mons[i]))
		}
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(30),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(RegularText).
		Background(DarkPurple).
		Bold(false)

	t.SetStyles(s)

	return t
}
