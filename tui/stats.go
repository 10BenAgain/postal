package tui

import (
	"fmt"
	pk "postal/pokemon"
	"postal/utils"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatEditor struct {
	*BaseEditorModel
	keys *EditorKeyMap
}

var (
	ExpandedStatNames = append(StatNames, "Is Egg", "Ability")
)

func (m *StatEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *StatEditor) View() string {
	editor := make([]string, len(m.inputs))
	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 0, 1, 2, 3, 4, 5:
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			if val > 31 {
				m.vals[m.focusIndex] = 31
			} else {
				m.vals[m.focusIndex] = uint(val)
			}
		}
	default:
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			if val > 1 {
				m.vals[m.focusIndex] = 1
			} else {
				m.vals[m.focusIndex] = uint(val)
			}
		}
	}

	ts := make([]string, len(ExpandedStatNames))

	for i := range m.inputs {
		s := lipgloss.NewStyle().Padding(0, 1)
		if i == m.focusIndex {
			ts[i] = s.Background(DarkerPurple).Render(fmt.Sprintf("%d", m.vals[i]))
		} else {
			ts[i] = s.Render(fmt.Sprintf("%d", m.vals[i]))
		}
	}

	t := lipgloss.JoinHorizontal(lipgloss.Center, ts...)

	view := makeFullStatsView(&m.pks)
	s1 := lipgloss.JoinVertical(lipgloss.Center, editor[0], editor[1], editor[2])
	s2 := lipgloss.JoinVertical(lipgloss.Center, editor[3], editor[4], editor[5])
	s3 := lipgloss.JoinVertical(lipgloss.Center, editor[6], editor[7])

	edit := lipgloss.JoinHorizontal(lipgloss.Top, s1, s2, s3)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		t,
		"",
		edit,
		"",
		view,
	)
}

func (m *StatEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func (m *StatEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub3.HpIV)
	m.vals[1] = uint(m.pks.Sub3.AtkIV)
	m.vals[2] = uint(m.pks.Sub3.DefIV)
	m.vals[3] = uint(m.pks.Sub3.SpeIV)
	m.vals[4] = uint(m.pks.Sub3.SpAtkIV)
	m.vals[5] = uint(m.pks.Sub3.SpDIV)
	m.vals[6] = uint(m.pks.Sub3.IsEgg)
	m.vals[7] = uint(m.pks.Sub3.AbilityNum)
}

func (m *StatEditor) CommitEdits() {
	m.pks.Sub3.HpIV = uint8(m.vals[0])
	m.pks.Sub3.AtkIV = uint8(m.vals[1])
	m.pks.Sub3.DefIV = uint8(m.vals[2])
	m.pks.Sub3.SpeIV = uint8(m.vals[3])
	m.pks.Sub3.SpAtkIV = uint8(m.vals[4])
	m.pks.Sub3.SpDIV = uint8(m.vals[5])
	m.pks.Sub3.IsEgg = uint8(m.vals[6])
	m.pks.Sub3.AbilityNum = uint8(m.vals[7])
}

func NewStatEditor() *StatEditor {
	return &StatEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedStatNames),
			vals:   make([]uint, len(ExpandedStatNames)),
			help:   help.New(),
		},
		keys: &EditorKeys,
	}
}

func makeFullStatsView(pks *pk.PStructure) string {
	val := pks.CalculateStatsValue()
	_, ab := pks.GetAbility()

	ivs := []uint8{
		pks.Sub3.HpIV, pks.Sub3.AtkIV, pks.Sub3.DefIV,
		pks.Sub3.SpeIV, pks.Sub3.SpAtkIV, pks.Sub3.SpDIV,
	}

	var joinedIVS []string

	for i := range 6 {
		joinedIVS = append(
			joinedIVS,
			joinConditionFieldValueString(
				StatNames[i], fmt.Sprintf("%d", ivs[i]),
			),
		)
	}

	ivMash := lipgloss.JoinVertical(
		lipgloss.Center,
		joinedIVS...,
	)

	remMash := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s %t",
			lipgloss.NewStyle().Foreground(DarkerPink).Render("Is Egg: "),
			utils.Uint2Bool(pks.Sub3.IsEgg),
		),
		fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(DarkerPink).Render("Ability:"),
			ab,
		),
	)

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		ivMash,
		remMash,
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		RibbonSumStyle.Render(fmt.Sprintf("Stats - 0x%02X", val)),
		" ",
		StatsWrapperStyle.Render(out),
	)
}
