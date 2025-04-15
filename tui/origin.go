package tui

import (
	"fmt"
	vals "postal/game"
	pk "postal/pokemon"
	"postal/utils"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OriginEditor struct {
	*BaseEditorModel
	keys *EditorKeyMap
}

func (m *OriginEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *OriginEditor) View() string {
	editor := make([]string, len(m.inputs))
	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 0: // Met Level
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	case 1: // Origin Game
		val, err := vals.OriginGameLookup(utils.SanitizeSearch(m.inputs[m.focusIndex].Value()))
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	case 2: // Ball
		val, err := vals.BallLookup(utils.SanitizeSearch(m.inputs[m.focusIndex].Value()))
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	case 3:
		i := utils.SanitizeSearch(m.inputs[m.focusIndex].Value())
		if i == "Male" {
			m.vals[m.focusIndex] = 0
		} else if i == "Female" {
			m.vals[m.focusIndex] = 1
		}
	}

	ts := make([]string, 4)

	for i := range m.inputs {
		s := lipgloss.NewStyle().Padding(0, 1)
		if i == m.focusIndex {
			ts[i] = s.Background(DarkerPurple).Render(fmt.Sprintf("%d", m.vals[i]))
		} else {
			ts[i] = s.Render(fmt.Sprintf("%d", m.vals[i]))
		}
	}

	t := lipgloss.JoinHorizontal(lipgloss.Center, ts...)

	view := makeFullOriginView(&m.pks)
	edit := lipgloss.JoinVertical(lipgloss.Center, editor...)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		t,
		"",
		edit,
		"",
		view,
	)
}

func (m *OriginEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func (m *OriginEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub3.MetLevel)
	m.vals[1] = uint(m.pks.Sub3.MetGame)
	m.vals[2] = uint(m.pks.Sub3.Ball)
	m.vals[3] = uint(m.pks.Sub3.OTGender)
}

func (m *OriginEditor) CommitEdits() {
	m.pks.Sub3.MetLevel = uint8(m.vals[0])
	m.pks.Sub3.MetGame = uint8(m.vals[1])
	m.pks.Sub3.Ball = uint8(m.vals[2])
	m.pks.Sub3.OTGender = uint8(m.vals[3])

}

func NewOriginEditor() *OriginEditor {
	return &OriginEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(OriginNames),
			vals:   make([]uint, len(OriginNames)),
		},
		keys: &EditorKeys,
	}
}

func makeFullOriginView(pks *pk.PStructure) string {
	val := pks.CalculateOriginValue()

	_, ball := pks.GetBall()
	_, otGen := pks.GetOTGender()
	iGame, nGame := pks.GetOriginGame()

	out := lipgloss.JoinVertical(
		lipgloss.Center,
		joinGrowthFieldValueString(OriginNames[0], fmt.Sprintf("%d", pks.Sub3.MetLevel)),
		joinGrowthFieldValueString(OriginNames[1], fmt.Sprintf("%s (#%d)", nGame, iGame)),
		joinGrowthFieldValueString(OriginNames[2], ball),
		joinGrowthFieldValueString(OriginNames[3], otGen),
	)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		RibbonSumStyle.Render(fmt.Sprintf("Origin - 0x%02X", val)),
		" ",
		out,
	)
}
