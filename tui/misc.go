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

type MiscEditor struct {
	*BaseEditorModel
	keys *EditorKeyMap
}

var (
	ExpandedMiscBlockNames = MiscBlockNames
)

func (m *MiscEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *MiscEditor) View() string {

	editor := make([]string, 5)

	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 1: // Met Location
		loc, err := vals.MetLocLookup(utils.SanitizeSearch(m.inputs[m.focusIndex].Value()))
		if err == nil {
			m.vals[m.focusIndex] = uint(loc)
		}
	default:
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	}

	s3 := m.pks.GetSubstructDataInOrder(false)[3]

	ret := lipgloss.JoinVertical(
		lipgloss.Center,
		editor...,
	)

	ts := make([]string, 5)

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
		makeMiscSubstructreView(&m.pks, subStructArray2String(s3)),
	)
}

func (m *MiscEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub3.Pokerus)
	m.vals[1] = uint(m.pks.Sub3.MetLocation)
	m.vals[2] = uint(m.pks.CalculateOriginValue())
	m.vals[3] = uint(m.pks.CalculateStatsValue())
	m.vals[4] = uint(m.pks.CalculateRibbonValue())
}

func (m *MiscEditor) CommitEdits() {
	m.pks.Sub3.Pokerus = uint8(m.vals[0])
	m.pks.Sub3.MetLocation = uint8(m.vals[1])
	m.pks.ReplaceOriginByValue(uint16(m.vals[2]))
	m.pks.ReplaceStatsByValue(uint32(m.vals[3]))
	m.pks.ReplaceRibbonByValue(uint32(m.vals[4]))
}

func (m *MiscEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func NewMiscEditor() *MiscEditor {
	return &MiscEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedMiscBlockNames),
			vals:   make([]uint, len(ExpandedMiscBlockNames)),
			help:   help.New(),
		},
		keys: &EditorKeys,
	}
}
