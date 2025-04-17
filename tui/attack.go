package tui

import (
	"fmt"
	vals "postal/game"
	"postal/utils"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AttackEditor struct {
	*BaseEditorModel
}

var (
	ExpandedAttacksBlockNames = []string{"Move 1", "Move 2", "Move 3", "Move 4", "PP 1", "PP 2", "PP 3", "PP 4"}
)

func (m *AttackEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *AttackEditor) View() string {
	editor := make([]string, len(m.inputs))

	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 0, 1, 2, 3:
		move, err := vals.MoveLookup(utils.SanitizeSearch(m.inputs[m.focusIndex].Value()))
		if err == nil {
			m.vals[m.focusIndex] = uint(move)
		}
	default:
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	}

	s1 := m.pks.GetSubstructDataInOrder(false)[1]
	editLeft := []string{editor[0], editor[1], editor[2], editor[3]}
	editRight := []string{editor[4], editor[5], editor[6], editor[7]}

	ret := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, editLeft...),
		lipgloss.JoinVertical(lipgloss.Center, editRight...),
	)

	ts := make([]string, 8)

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
		makeAttacksSubStructView(&m.pks, subStructArray2String(s1)),
	)
}

func (m *AttackEditor) CommitEdits() {
	for i := range 4 {
		m.pks.Sub1.Moves[i] = uint16(m.vals[i])
		m.pks.Sub1.PP[i] = uint8(m.vals[i+4])
	}
}

func (m *AttackEditor) UpdateValues() {
	for i := range 4 {
		m.vals[i] = uint(m.pks.Sub1.Moves[i])
		m.vals[i+3] = uint(m.pks.Sub1.PP[i])
	}
}

func (m *AttackEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func NewAttackEditor() *AttackEditor {
	return &AttackEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedAttacksBlockNames),
			vals:   make([]uint, len(ExpandedAttacksBlockNames)),
			help:   help.New(),
			keys:   &EditorKeys,
		},
	}
}
