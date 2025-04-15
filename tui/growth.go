package tui

import (
	"fmt"
	vals "postal/game"
	"postal/utils"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GrowthEditor struct {
	*BaseEditorModel
}

var (
	ExpandedGrowthBlockNames = []string{"Species", "Item", "EXP", "Friendship", "PP 1", "PP 2", "PP 3", "PP 4"}
)

func (m *GrowthEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

func (m *GrowthEditor) View() string {
	editor := make([]string, len(m.inputs))
	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 0: // Species
		mon, err := vals.MonLookup(utils.SanitizeSearch(m.inputs[0].Value()))
		if err == nil {
			m.vals[0] = uint(mon)
		}
	case 1: // Item
		item, err := vals.ItemLookup(utils.SanitizeSearch(m.inputs[1].Value()))
		if err == nil {
			m.vals[1] = uint(item)
		}
	default: // Rest of values that need numbers to parse
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			m.vals[m.focusIndex] = uint(val)
		}
	}

	s0 := m.pks.GetSubstructDataInOrder(false)[0]
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
		makeGrowthSubStructView(&m.pks, subStructArray2String(s0)),
	)
}

func (m *GrowthEditor) CommitEdits() {
	m.pks.Sub0.Species = uint16(m.vals[0])
	m.pks.Sub0.HeldItem = uint16(m.vals[1])
	m.pks.Sub0.Experience = uint32(m.vals[2])

	m1 := uint8(m.vals[4]) & 0b11
	m2 := uint8(m.vals[5]) & 0b11
	m3 := uint8(m.vals[6]) & 0b11
	m4 := uint8(m.vals[7]) & 0b11

	m.pks.Sub0.PPBonuses = (m1 << 0) | (m2 << 2) | (m3 << 4) | (m4 << 6)
	m.pks.Sub0.FriendShip = uint8(m.vals[3])
}

func (m *GrowthEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub0.Species)
	m.vals[1] = uint(m.pks.Sub0.HeldItem)
	m.vals[2] = uint(m.pks.Sub0.Experience)
	m.vals[3] = uint(m.pks.Sub0.FriendShip)

	b := m.pks.GetPPBonus()
	for i := range 3 {
		m.vals[i+4] = b[i]
	}
}

func (m *GrowthEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func NewGrowthEditor() *GrowthEditor {
	return &GrowthEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedGrowthBlockNames),
			vals:   make([]uint, len(ExpandedGrowthBlockNames)),
			keys:   &EditorKeys,
		},
	}
}
