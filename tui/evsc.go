package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EVSCEditor struct {
	*BaseEditorModel
}

var (
	ExpandedEVSCBlockNames = append(StatNames, ConditionNames...)
)

func (m *EVSCEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.keys.Commit):
			m.CommitEdits()
		case key.Matches(msg, m.keys.Reset):
			m.UpdateValues()
			m.inputs[m.focusIndex].Reset()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Right),
			key.Matches(msg, m.keys.Left),
			key.Matches(msg, m.keys.Up),
			key.Matches(msg, m.keys.Down):

			switch {
			case key.Matches(msg, m.keys.Right):
				m.focusIndex++
			case key.Matches(msg, m.keys.Left):
				m.focusIndex--
			case key.Matches(msg, m.keys.Up):
				m.focusIndex--
			case key.Matches(msg, m.keys.Down):
				m.focusIndex++
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
	return m, m.UpdateInputs(msg)
}

func (m *EVSCEditor) View() string {
	editor := make([]string, 12)

	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	if val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32); err == nil {
		m.vals[m.focusIndex] = uint(val)
	}

	s2 := m.pks.GetSubstructDataInOrder(false)[2]
	editLeft := []string{editor[0], editor[1], editor[2], editor[3], editor[4], editor[5]}
	editRight := []string{editor[6], editor[7], editor[8], editor[9], editor[10], editor[11]}

	ret := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, editLeft...),
		lipgloss.JoinVertical(lipgloss.Center, editRight...),
	)

	ts := make([]string, 12)

	for i := range m.inputs {
		s := lipgloss.NewStyle().Padding(0, 1)
		if i == m.focusIndex {
			ts[i] = s.Background(DarkerPurple).Render(fmt.Sprintf("%d", m.vals[i]))
		} else {
			ts[i] = s.Render(fmt.Sprintf("%d", m.vals[i]))
		}
	}

	t := lipgloss.JoinHorizontal(lipgloss.Center, ts...)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		t,
		" ",
		ret,
		" ",
		makeConditionsSubStructView(&m.pks, subStructArray2String(s2)),
	)
}

func (m *EVSCEditor) CommitEdits() {
	m.pks.Sub2.HpEV = uint8(m.vals[0])
	m.pks.Sub2.AtkEV = uint8(m.vals[1])
	m.pks.Sub2.DefEV = uint8(m.vals[2])
	m.pks.Sub2.SpeEV = uint8(m.vals[3])
	m.pks.Sub2.SpAtkEV = uint8(m.vals[4])
	m.pks.Sub2.SpDefEV = uint8(m.vals[5])

	m.pks.Sub2.Cool = uint8(m.vals[6])
	m.pks.Sub2.Beauty = uint8(m.vals[7])
	m.pks.Sub2.Cute = uint8(m.vals[8])
	m.pks.Sub2.Smart = uint8(m.vals[9])
	m.pks.Sub2.Tough = uint8(m.vals[10])
	m.pks.Sub2.Feel = uint8(m.vals[11])

}

func (m *EVSCEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub2.HpEV)
	m.vals[1] = uint(m.pks.Sub2.AtkEV)
	m.vals[2] = uint(m.pks.Sub2.DefEV)
	m.vals[3] = uint(m.pks.Sub2.SpeEV)
	m.vals[4] = uint(m.pks.Sub2.SpAtkEV)
	m.vals[5] = uint(m.pks.Sub2.SpDefEV)

	m.vals[6] = uint(m.pks.Sub2.Cool)
	m.vals[7] = uint(m.pks.Sub2.Beauty)
	m.vals[8] = uint(m.pks.Sub2.Cute)
	m.vals[9] = uint(m.pks.Sub2.Smart)
	m.vals[10] = uint(m.pks.Sub2.Tough)
	m.vals[11] = uint(m.pks.Sub2.Feel)
}

func (m *EVSCEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func NewEVSCEditor() *EVSCEditor {
	return &EVSCEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedEVSCBlockNames),
			vals:   make([]uint, len(ExpandedEVSCBlockNames)),
			help:   help.New(),
			keys:   &EditorKeys,
		},
	}
}
