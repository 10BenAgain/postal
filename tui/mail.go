package tui

import (
	"fmt"
	"log"
	"os"
	"postal/pokemon"
	"postal/utils"

	vals "postal/game"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

type SearchMode int

const (
	AllWords = iota
	WordMonOne
	WordMonTwo
)

var ModeNames = []string{"ALL", "PK1", "PK2"}

type MailEditor struct {
	*BaseEditorModel
	words  []string
	tbl    table.Model
	pksNew pokemon.PStructure
	mode   SearchMode
	keys   *MailKeyMap
}

func (m *MailEditor) getCurrentValues() []uint16 {
	h, l := m.pks.GetPIDSplit()
	tid, sid := m.pks.GetTIDSIDComboPair()

	// Temp vals for comparison
	return []uint16{h, l, sid, tid}
}

func (m *MailEditor) GenerateNewPIDOTValues() (uint32, uint32) {
	vl := m.getCurrentValues()
	for i := range vl {
		a := m.vals[i]
		if a > 0 {
			vl[i] = uint16(a)
		}
	}

	return pokemon.CombineU16Split(vl[0], vl[1]), pokemon.CombineU16Split(vl[2], vl[3])
}

func (m *MailEditor) SetNewPokemon(pks pokemon.PStructure) {
	m.pksNew = pks
}

func (m *MailEditor) SwapEdits() {
	m.pks = m.pksNew
}

func (m *MailEditor) GetEditMonView() string {
	return generateMonViewOrder(&m.pks)
}

func (m *MailEditor) GetResultMonView() string {
	return generateMonViewOrder(&m.pksNew)
}

func (m *MailEditor) SaveMonToFile() {
	d, n := m.pks.GetSpecies()
	f := fmt.Sprintf("%03d-%s-%08X.pk3", d, n, m.pks.PID)

	err := os.WriteFile(f, m.pks.ToPK3(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MailEditor) CommitEdits() {
	r := m.pks.SaveWithMail(m.GenerateNewPIDOTValues())
	n := pokemon.GeneratePokemonFromRawData(r, false, false)
	m.SetNewPokemon(n)
}

func (m *MailEditor) UpdateValues() {
	vl := m.getCurrentValues()

	for i := range m.vals {
		m.vals[i] = uint(vl[i])
	}
}

func (m *MailEditor) GetKeys() string {
	return m.help.View(m.keys)
}

func (m *MailEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case returnMsg:
		m.SetNewPokemon(msg.pk)

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Swap):
			m.SwapEdits()
		case key.Matches(msg, m.keys.Mode):
			if m.mode < 2 {
				m.mode++
			} else {
				m.mode = 0
			}
		case key.Matches(msg, m.keys.Commit):
			m.CommitEdits()
		case key.Matches(msg, m.keys.Swap):
			m.SwapEdits()
			return m, nil
		case key.Matches(msg, m.keys.File):
			m.SaveMonToFile()
			return m, nil

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
				m.focusIndex -= 2
			case key.Matches(msg, m.keys.Down):
				m.focusIndex += 2
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

func (m *MailEditor) View() string {
	editor := make([]string, len(m.inputs))

	var val uint16
	var err error

	for i := range m.inputs {
		editor[i] = WordEntryStyle.Render(m.inputs[i].View())

		s := utils.SanitizeWordSearch(m.inputs[i].Value())

		switch m.mode {
		case AllWords:
			val, err = vals.WordLookup(s)
		case WordMonOne:
			val, err = vals.WordMon1Lookup(s)
		case WordMonTwo:
			val, err = vals.WordMon2Lookup(s)
		}

		switch i {
		case 0: //PIDLO
			if err != nil {
				m.vals[1] = 0
				m.words[1] = "0000"
			} else {
				m.vals[1] = uint(val)
				m.words[1] = fmt.Sprintf("%04X", val)
			}
		case 1: // TID
			if err != nil {
				m.vals[0] = 0
				m.words[0] = "0000"
			} else {
				m.vals[0] = uint(val)
				m.words[0] = fmt.Sprintf("%04X", val)
			}
		case 2: // PIDHI
			if err != nil {
				m.vals[3] = 0
				m.words[3] = "0000"
			} else {
				m.vals[3] = uint(val)
				m.words[3] = fmt.Sprintf("%04X", val)
			}
		case 3: // SID
			if err != nil {
				m.vals[2] = 0
				m.words[2] = "0000"
			} else {
				m.vals[2] = uint(val)
				m.words[2] = fmt.Sprintf("%04X", val)
			}
		}
	}

	h, l := m.pks.GetPIDSplit()
	tid, sid := m.pks.GetTIDSIDComboPair()

	stats := []string{
		fmt.Sprintf("%04X", h),
		fmt.Sprintf("%04X", l),
		fmt.Sprintf("%04X", sid),
		fmt.Sprintf("%04X", tid),
	}

	m.tbl.SetRows([]table.Row{
		{"Mon", stats[0], stats[1], stats[2], stats[3]},
		{"Words", m.words[0], m.words[1], m.words[2], m.words[3]},
	})

	editLeft := lipgloss.JoinVertical(lipgloss.Center, editor[0], editor[2])
	editRight := lipgloss.JoinVertical(lipgloss.Center, editor[1], editor[3])

	newPID, newOTID := m.GenerateNewPIDOTValues()
	xkey := newPID ^ newOTID

	editorJoin := MailMenuStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			editLeft,
			editRight,
		),
	)

	res := lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("\nSearch Mode: %s", ModeNames[m.mode]),
		editorJoin,
		m.tbl.View(),
		MenuStyle.Render(lipgloss.JoinVertical(
			lipgloss.Center,
			fmt.Sprintf("XKEY: %08X\n", m.pks.GetEncryptionKey()),
			fmt.Sprintf("PID: %08X", newPID),
			fmt.Sprintf("OTID: %08X", newOTID),
			fmt.Sprintf("XKEY: %08X", xkey),
		)),
	)

	edit := m.GetEditMonView()
	_, height := lipgloss.Size(edit)
	center := lipgloss.NewStyle().
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DarkerPink)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		edit,
		center.Render(res),
		m.GetResultMonView(),
	)
}

func NewMailEditor() *MailEditor {
	m := BaseEditorModel{
		inputs: make([]textinput.Model, 4),
		vals:   make([]uint, 4),
		help:   help.New(),
	}

	var ti textinput.Model
	for i := range m.inputs {
		ti = textinput.New()
		ti.Placeholder = WordEntry
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(LightPink)
		ti.CharLimit = 156
		ti.Width = 12

		if i == 0 {
			ti.Focus()
			ti.PromptStyle = FocusedStyle
		}
		m.inputs[i] = ti
	}

	cols := []table.Column{
		{Title: "", Width: 6},
		{Title: "PIDHI", Width: 8},
		{Title: "PIDLO", Width: 8},
		{Title: "SID", Width: 8},
		{Title: "TID", Width: 8},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithHeight(4),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()

	s.Cell = s.Cell.
		Foreground(RegularText)

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(LightGreenish).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)

	return &MailEditor{
		BaseEditorModel: &m,
		tbl:             t,
		words:           BlankWordValues,
		mode:            AllWords,
		keys:            &MailKeys,
	}
}

func makePokemonTitleBar(pks *pokemon.PStructure) string {
	_, name := pks.GetSpecies()
	_, order := pks.GetMonSubStructOrderString()

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		MSpecies.Render(name),
		MOTID.Render(fmt.Sprintf("%08X", pks.OTID)),
		MPID.Render(fmt.Sprintf("%08X", pks.PID)),
		StructOrder.Render(order),
	)
}

func makeGrowthSubStructView(pks *pokemon.PStructure, s string) string {
	index, _ := pks.GetSpecies()
	itemDex, item := pks.GetHeldItem()
	pb := pks.GetPPBonus()
	exp := pks.Sub0.Experience
	fs := pks.Sub0.FriendShip

	mash := lipgloss.JoinVertical(
		lipgloss.Center,
		SubStructByte.Render(s),
		joinGrowthFieldValueString(GrowthBlockNames[0], fmt.Sprintf("%d (0x%02X)", index, index)),
		joinGrowthFieldValueString(GrowthBlockNames[1], fmt.Sprintf("%s (0x%02X)", item, itemDex)),
		joinGrowthFieldValueString(GrowthBlockNames[2], fmt.Sprintf("%d (0x%02X)", exp, exp)),
		joinGrowthFieldValueString(GrowthBlockNames[3], fmt.Sprintf("%d %d %d %d", pb[0], pb[1], pb[2], pb[3])),
		joinGrowthFieldValueString(GrowthBlockNames[4], fmt.Sprintf("%d (0x%02X)", fs, fs)),
	)

	return SubStructBlock.Render(mash)
}

func makeAttacksSubStructView(pks *pokemon.PStructure, s string) string {
	p := pks.Sub1.PP
	m := pks.Sub1.Moves
	n := pks.GetMoveStrings()

	mList := list.New(
		"Moves", list.New(
			n[0],
			n[1],
			n[2],
			n[3],
		).
			Enumerator(BlankEnum).
			ItemStyle(SubListHeader),
	).
		Enumerator(BlankEnum).
		ItemStyle(ListHeader)

	pList := list.New(
		"PP", list.New(
			fmt.Sprintf("%02d", p[0]),
			fmt.Sprintf("%02d", p[1]),
			fmt.Sprintf("%02d", p[2]),
			fmt.Sprintf("%02d", p[3]),
		).
			Enumerator(BlankEnum).
			ItemStyle(SubListHeader),
	).
		Enumerator(BlankEnum).
		ItemStyle(ListHeader)

	nList := list.New(
		"Index", list.New(
			fmt.Sprintf("0x%02X", m[0]),
			fmt.Sprintf("0x%02X", m[1]),
			fmt.Sprintf("0x%02X", m[2]),
			fmt.Sprintf("0x%02X", m[3]),
		).
			Enumerator(BlankEnum).
			ItemStyle(SubListHeader),
	).
		Enumerator(BlankEnum).
		ItemStyle(ListHeader)

	mash := lipgloss.JoinHorizontal(
		lipgloss.Center,
		mList.String(),
		pList.String(),
		nList.String(),
	)

	return SubStructBlock.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			SubStructByte.Render(s),
			mash,
		),
	)
}

func makeConditionsSubStructView(pks *pokemon.PStructure, s string) string {
	evs := []uint8{
		pks.Sub2.HpEV,
		pks.Sub2.AtkEV,
		pks.Sub2.DefEV,
		pks.Sub2.SpeEV,
		pks.Sub2.SpAtkEV,
		pks.Sub2.SpDefEV,
	}

	conds := []uint8{
		pks.Sub2.Cool,
		pks.Sub2.Beauty,
		pks.Sub2.Cute,
		pks.Sub2.Smart,
		pks.Sub2.Tough,
		pks.Sub2.Feel,
	}

	var joinedEVS []string
	var joinedConditions []string

	for i := range 6 {
		joinedEVS = append(
			joinedEVS,
			joinConditionFieldValueString(
				StatNames[i], fmt.Sprintf("%d", evs[i]),
			),
		)
		joinedConditions = append(
			joinedConditions,
			joinConditionFieldValueString(
				ConditionNames[i], fmt.Sprintf("%d", conds[i]),
			),
		)
	}

	evMash := lipgloss.JoinVertical(
		lipgloss.Center,
		joinedEVS...,
	)

	condMash := lipgloss.JoinVertical(
		lipgloss.Center,
		joinedConditions...,
	)

	mash := lipgloss.JoinHorizontal(
		lipgloss.Center,
		evMash,
		condMash,
	)

	return SubStructBlock.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			SubStructByte.Render(s),
			mash,
		),
	)
}

func makeMiscSubstructreView(pks *pokemon.PStructure, s string) string {
	sub := pks.Sub3
	locIndex, locName := pks.GetMetLocation()

	mash := lipgloss.JoinVertical(
		lipgloss.Center,
		joinGrowthFieldValueString(MiscBlockNames[0], fmt.Sprintf("%d (0x%02X)", sub.Pokerus, sub.Pokerus)),
		joinGrowthFieldValueString(MiscBlockNames[1], fmt.Sprintf("%s (0x%02X)", locName, locIndex)),
		joinGrowthFieldValueString(MiscBlockNames[2], fmt.Sprintf("0x%02X", pks.CalculateOriginValue())),
		joinGrowthFieldValueString(MiscBlockNames[3], fmt.Sprintf("0x%02X", pks.CalculateStatsValue())),
		joinGrowthFieldValueString(MiscBlockNames[4], fmt.Sprintf("0x%02X", pks.CalculateRibbonValue())),
	)

	return SubStructBlock.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			SubStructByte.Render(s),
			mash,
		),
	)
}

func generateMonViewOrder(pks *pokemon.PStructure) string {
	order := pks.GetMonSubStructOrder()
	subs := pks.GetSubstructDataInOrder(false)

	subViews := map[int]string{
		0: makeGrowthSubStructView(pks, subStructArray2String(subs[0])),
		1: makeAttacksSubStructView(pks, subStructArray2String(subs[1])),
		2: makeConditionsSubStructView(pks, subStructArray2String(subs[2])),
		3: makeMiscSubstructreView(pks, subStructArray2String(subs[3])),
	}

	var out []string
	out = append(out, makePokemonTitleBar(pks))
	out = append(out, " ")

	for _, i := range order {
		out = append(out, subViews[i])
	}

	return OuterBorder.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			out...,
		),
	)
}
