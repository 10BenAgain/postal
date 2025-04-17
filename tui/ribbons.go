package tui

import (
	"fmt"
	pk "postal/pokemon"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

type RibbonEditor struct {
	*BaseEditorModel
	keys *EditorKeyMap
}

var (
	ExpandedRibbonNames = mergeRibbonNames()
)

func mergeRibbonNames() []string {
	var r []string
	l := [][]string{RibbonBlockOne, RibbonBlockTwo, RibbonBlockThree, RibbonBlockFour, RibbonBlockFive}

	for i := range l {
		r = append(r, l[i]...)
	}

	return r
}

func (m *RibbonEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *RibbonEditor) View() string {
	editor := make([]string, len(m.inputs))
	for i := range m.inputs {
		e := m.inputs[i].View()
		editor[i] = WordEntryStyle.UnsetPadding().Render(e)
	}

	switch m.focusIndex {
	case 0, 1, 2, 3, 4:
		val, err := strconv.ParseUint(m.inputs[m.focusIndex].Value(), 10, 32)
		if err == nil {
			if val > 4 {
				m.vals[m.focusIndex] = 4
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

	ts := make([]string, len(ExpandedRibbonNames))

	for i := range m.inputs {
		s := lipgloss.NewStyle().Padding(0, 1)
		if i == m.focusIndex {
			ts[i] = s.Background(DarkerPurple).Render(fmt.Sprintf("%d", m.vals[i]))
		} else {
			ts[i] = s.Render(fmt.Sprintf("%d", m.vals[i]))
		}
	}

	t := lipgloss.JoinHorizontal(lipgloss.Center, ts...)
	view := makeFullRibbonView(&m.pks)
	r1 := lipgloss.JoinVertical(lipgloss.Center, editor[0], editor[1], editor[2], editor[3], editor[4])
	r2 := lipgloss.JoinVertical(lipgloss.Center, editor[5], editor[6], editor[7], editor[8], editor[9])
	r3 := lipgloss.JoinVertical(lipgloss.Center, editor[10], editor[11], editor[12])
	r4 := lipgloss.JoinVertical(lipgloss.Center, editor[13], editor[14], editor[15], editor[16])

	edit := lipgloss.JoinHorizontal(lipgloss.Top, r1, r2, r3, r4, editor[17])

	return lipgloss.JoinVertical(
		lipgloss.Center,
		t,
		"",
		edit,
		"",
		view,
	)
}

func (m *RibbonEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func (m *RibbonEditor) UpdateValues() {
	m.vals[0] = uint(m.pks.Sub3.CoolRibbon)
	m.vals[1] = uint(m.pks.Sub3.BeautyRibon)
	m.vals[2] = uint(m.pks.Sub3.CuteRibbon)
	m.vals[3] = uint(m.pks.Sub3.SmartRibbon)
	m.vals[4] = uint(m.pks.Sub3.ToughRibbon)

	m.vals[5] = uint(m.pks.Sub3.ChampionRibbon)
	m.vals[6] = uint(m.pks.Sub3.WinningRibbon)
	m.vals[7] = uint(m.pks.Sub3.VictoryRibbon)
	m.vals[8] = uint(m.pks.Sub3.ArtistRibbon)
	m.vals[9] = uint(m.pks.Sub3.EffortRibbon)

	m.vals[10] = uint(m.pks.Sub3.BattleChampionRibbon)
	m.vals[11] = uint(m.pks.Sub3.RegionalChampionRibbon)
	m.vals[12] = uint(m.pks.Sub3.NationalChampionRibbon)

	m.vals[13] = uint(m.pks.Sub3.CountryRibbon)
	m.vals[14] = uint(m.pks.Sub3.NationalRibbon)
	m.vals[15] = uint(m.pks.Sub3.EarthRibbon)
	m.vals[16] = uint(m.pks.Sub3.WorldRibbon)

	m.vals[17] = uint(m.pks.Sub3.FatefulEnc)
}

func (m *RibbonEditor) CommitEdits() {
	m.pks.Sub3.CoolRibbon = uint8(m.vals[0])
	m.pks.Sub3.BeautyRibon = uint8(m.vals[1])
	m.pks.Sub3.CuteRibbon = uint8(m.vals[2])
	m.pks.Sub3.SmartRibbon = uint8(m.vals[3])
	m.pks.Sub3.ToughRibbon = uint8(m.vals[4])

	m.pks.Sub3.ChampionRibbon = uint8(m.vals[5])
	m.pks.Sub3.WinningRibbon = uint8(m.vals[6])
	m.pks.Sub3.VictoryRibbon = uint8(m.vals[7])
	m.pks.Sub3.ArtistRibbon = uint8(m.vals[8])
	m.pks.Sub3.EffortRibbon = uint8(m.vals[9])

	m.pks.Sub3.BattleChampionRibbon = uint8(m.vals[10])
	m.pks.Sub3.RegionalChampionRibbon = uint8(m.vals[11])
	m.pks.Sub3.NationalChampionRibbon = uint8(m.vals[12])

	m.pks.Sub3.CountryRibbon = uint8(m.vals[13])
	m.pks.Sub3.NationalRibbon = uint8(m.vals[14])
	m.pks.Sub3.EarthRibbon = uint8(m.vals[15])
	m.pks.Sub3.WorldRibbon = uint8(m.vals[16])

	m.pks.Sub3.FatefulEnc = uint8(m.vals[17])
}

func NewRibbonEditor() *RibbonEditor {
	return &RibbonEditor{
		BaseEditorModel: &BaseEditorModel{
			inputs: makeEditorTextModels(ExpandedRibbonNames),
			vals:   make([]uint, len(ExpandedRibbonNames)),
			help:   help.New(),
		},
		keys: &EditorKeys,
	}
}

func BlankEnum(l list.Items, i int) string {
	return ""
}

func blankCheckMarkEnumerator(l list.Items, i int) string {
	val := l.At(i).Value()
	if val[len(val)-1:] == " " {
		return RibbonCheckEnumStyle.Render(CheckMark)
	}
	return RibbonXMarkEnumStyle.Render(XMark)
}

func makeVRibbonList(n []string, v []uint8) *list.List {
	l := list.New().Enumerator(blankCheckMarkEnumerator)
	for i := range n {
		if v[i] > 0 {
			l.Item(fmt.Sprintf("%d %s ", v[i], n[i]))
		} else {
			l.Item(fmt.Sprintf("%d %s", v[i], n[i]))
		}
	}
	return l
}

func makeRibbonList(n []string, v []uint8) *list.List {
	l := list.New().Enumerator(blankCheckMarkEnumerator)
	for i := range n {
		if v[i] > 0 {
			l.Item(fmt.Sprintf("%s ", n[i]))
		} else {
			l.Item(n[i])
		}
	}

	return l
}

func makeFullRibbonView(pks *pk.PStructure) string {
	sub := pks.Sub3
	rs := pks.CalculateRibbonValue()

	rv1 := []uint8{
		sub.CoolRibbon,
		sub.BeautyRibon,
		sub.CuteRibbon,
		sub.SmartRibbon,
		sub.ToughRibbon,
	}

	rv2 := []uint8{
		sub.ChampionRibbon,
		sub.WinningRibbon,
		sub.VictoryRibbon,
		sub.ArtistRibbon,
		sub.EffortRibbon,
	}

	rv3 := []uint8{
		sub.BattleChampionRibbon,
		sub.RegionalChampionRibbon,
		sub.NationalChampionRibbon,
	}

	rv4 := []uint8{
		sub.CountryRibbon,
		sub.NationalRibbon,
		sub.EarthRibbon,
		sub.WorldRibbon,
	}

	rv5 := []uint8{
		sub.FatefulEnc,
	}

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		RibbonListBlock.Render(makeVRibbonList(RibbonBlockOne, rv1).String()),
		RibbonListBlock.Render(makeRibbonList(RibbonBlockTwo, rv2).String()),
		RibbonListBlock.Render(makeRibbonList(RibbonBlockThree, rv3).String()),
		RibbonListBlock.Render(makeRibbonList(RibbonBlockFour, rv4).String()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		RibbonSumStyle.Render(fmt.Sprintf("Ribbons - 0x%X - %b", rs, rs)),
		out,
		RibbonListBlock.Render(makeRibbonList(RibbonBlockFive, rv5).String()),
	)
}
