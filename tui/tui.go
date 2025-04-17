package tui

import (
	"errors"
	"fmt"
	"os"
	"path"
	"postal/pokemon"
	"postal/save"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	mailEdit = iota
	growthEdit
	attackEdit
	evscEdit
	miscEdit
	statEdit
	originEdit
	ribbonEdit
	searcher
	boxView
	picker
)

type editorState int

type EditorModel interface {
	// Tea Model interface methods
	tea.Model

	// All editor models need this to update UI
	UpdateInputs(tea.Msg) tea.Cmd

	// For pushing and pulling Pokemon data up and down the model chain
	GetPokemon() pokemon.PStructure
	SetPokemon(pokemon.PStructure)

	// Update the Pokemon data from user input
	CommitEdits()

	// Update current stored values
	UpdateValues()

	// Generate help key string to send up the model chain
	GetKeys() string
}

type BaseEditorModel struct {
	focusIndex int
	vals       []uint
	inputs     []textinput.Model
	pks        pokemon.PStructure
	help       help.Model
	keys       *EditorKeyMap
}

func (m *BaseEditorModel) GetPokemon() pokemon.PStructure {
	return m.pks
}

func (m *BaseEditorModel) SetPokemon(pks pokemon.PStructure) {
	m.pks = pks
}

func (m *BaseEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *BaseEditorModel) UpdateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func makeEditorTextModels(s []string) []textinput.Model {
	t := make([]textinput.Model, len(s))

	var ti textinput.Model
	for i, v := range s {
		ti = textinput.New()
		ti.Placeholder = v
		ti.Cursor.Style = lipgloss.NewStyle().Foreground(LightPink)
		ti.CharLimit = 156
		ti.Width = 12

		if i == 0 {
			ti.Focus()
			ti.PromptStyle = FocusedStyle
		}
		t[i] = ti
	}

	return t
}

type MainModel struct {
	state editorState

	editors []EditorModel

	searcher tea.Model
	boxView  tea.Model

	picker       filepicker.Model
	selectedFile string
	err          error

	pks    pokemon.PStructure
	keys   *EditorKeyMap
	help   help.Model
	status string

	height int
	width  int
}

type clearErrorMsg struct{}

type clearStatusMsg struct{}

type statusMsg struct{ stat string }

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func clearStatus() tea.Cmd {
	return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func makeFilePicker() filepicker.Model {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowHidden = false
	fp.AllowedTypes = []string{".pk3", ".ek3", ".sav", ".SAV"}

	return fp
}

func InitMainNoFile() MainModel {
	// Start with blank mon, could be changed to something more useful
	blank := pokemon.GeneratePokemonFromRawData(pokemon.BlankSpecies, false, true)
	m := InitMainEditor(blank)

	// Set the initial state to file picker
	m.state = picker

	return m
}

func InitMainBoxData(b save.RawBoxDataTotal) MainModel {
	blank := pokemon.GeneratePokemonFromRawData(pokemon.BlankSpecies, false, true)
	m := InitMainEditor(blank)

	// Make a new box view and switch state to handle it
	m.boxView = NewBoxSelect(b)
	m.state = boxView

	return m
}

func InitMainEditor(pk pokemon.PStructure) MainModel {
	mail := NewMailEditor()
	mail.SetPokemon(pk)
	mail.SetNewPokemon(pk)

	return MainModel{
		state: mailEdit,

		editors: []EditorModel{
			mail,
			NewGrowthEditor(),
			NewAttackEditor(),
			NewEVSCEditor(),
			NewMiscEditor(),
			NewStatEditor(),
			NewOriginEditor(),
			NewRibbonEditor(),
		},

		searcher: NewSearchModel(),
		picker:   makeFilePicker(),

		pks:  pk,
		keys: &EditorKeys,
		help: help.New(),
	}
}

func (m MainModel) Init() tea.Cmd {
	m.updateEditors()

	c := make([]tea.Cmd, len(m.editors))
	for i := range c {
		c[i] = m.editors[i].Init()
	}
	c = append(c, m.picker.Init())
	c = append(c, textinput.Blink)

	return tea.Batch(c...)
}

func (m *MainModel) setMon(pk pokemon.PStructure) {
	m.pks = pk
}

func (m *MainModel) updateEditors() {
	for i := range m.editors {
		m.editors[i].SetPokemon(m.pks)
	}
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height

	case returnMsg:
		m.setMon(msg.pk)
		m.updateEditors()
		m.state = mailEdit

		_, cmd := m.editors[mailEdit].Update(msg)
		cmds = append(cmds, cmd)

	case statusMsg:
		m.status = msg.stat
		cmds = append(cmds, clearStatus())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		// Toggle help on all editors for convenience
		case key.Matches(msg, m.keys.Help):
			for i := range m.editors {
				_, cmd := m.editors[i].Update(msg)
				cmds = append(cmds, cmd)
			}

			_, cmd := m.boxView.Update(msg)
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)

		// Editor hotkeys
		case key.Matches(msg, m.keys.Mail):
			m.state = mailEdit
			return m, nil

		case key.Matches(msg, m.keys.Growth):
			m.state = growthEdit
			return m, nil

		case key.Matches(msg, m.keys.Attacks):
			m.state = attackEdit
			return m, nil

		case key.Matches(msg, m.keys.EVSC):
			m.state = evscEdit
			return m, nil

		case key.Matches(msg, m.keys.Misc):
			m.state = miscEdit
			return m, nil

		case key.Matches(msg, m.keys.Search):
			m.state = searcher
			return m, nil

		case key.Matches(msg, m.keys.Origin):
			m.state = originEdit
			return m, nil

		case key.Matches(msg, m.keys.Ribbon):
			m.state = ribbonEdit
			return m, nil

		case key.Matches(msg, m.keys.Stats):
			m.state = statEdit
			return m, nil

		case key.Matches(msg, m.keys.File):
			m.state = picker

			// Apparently the file picker needs a manual resize here...
			return m, tea.Batch(
				m.picker.Init(),
				func() tea.Msg {
					return tea.WindowSizeMsg{Width: m.width, Height: m.height}
				},
			)

		case key.Matches(msg, m.keys.Save):
			m.setMon(m.editors[m.state].GetPokemon())

			if m.state == mailEdit {
				m.updateEditors()
			}

			m.editors[mailEdit].SetPokemon(m.pks)
			m.status = "Saving active mon data.."
			cmds = append(cmds, clearStatus())
		}

		switch msg.String() {
		// For cycling thru sub editor states without using hotkeys
		// This sould probably be changed to use key matching so.. TODO: (grep)

		case "tab":
			if m.state >= growthEdit && m.state <= ribbonEdit {
				if m.state == ribbonEdit {
					m.state = growthEdit
				} else {
					m.state++
				}
			}

		case "shift+tab":
			if m.state >= growthEdit && m.state <= ribbonEdit {
				if m.state == growthEdit {
					m.state = ribbonEdit
				} else {
					m.state--
				}
			}
		}

	case clearErrorMsg:
		m.err = nil

	case clearStatusMsg:
		m.status = ""
	}

	var cmd tea.Cmd

	switch m.state {

	case mailEdit, growthEdit, miscEdit, attackEdit,
		evscEdit, statEdit, originEdit, ribbonEdit:
		_, cmd = m.editors[m.state].Update(msg)

	case searcher:
		_, cmd = m.searcher.Update(msg)

	case picker:
		m.picker, cmd = m.picker.Update(msg)

		// Did the user select a file?
		if didSelect, p := m.picker.DidSelectFile(msg); didSelect {
			// Get the path of the selected file.
			m.selectedFile = p

			f := func(pk pokemon.PStructure) {
				m.setMon(pk)
				m.editors[mailEdit].SetPokemon(m.pks)
				m.updateEditors()

				m.state = mailEdit
			}

			switch path.Ext(p) {
			case ".pk3":
				pk, err := save.GetRawMonDataFromFile(p)
				if err != nil {
					m.err = errors.New("unable to read .pk3 file")
				} else {
					mon := pokemon.GeneratePokemonFromRawData(pk, false, true)
					f(mon)
				}

			case ".ek3":
				pk, err := save.GetRawMonDataFromFile(p)
				if err != nil {
					m.err = errors.New("unable to read .ek3 file")
				} else {
					mon := pokemon.GeneratePokemonFromRawData(pk, false, false)
					f(mon)
				}

			case ".sav", ".SAV":
				m.state = boxView
				sav := save.GenerateRawSaveData(p)
				block := sav.GenerateSaveBlock()

				m.boxView = NewBoxSelect(block.GetRawBoxData())
			}
		}

		// Did the user select a disabled file?
		// This is only necessary to display an error to the user.
		if didSelect, path := m.picker.DidSelectDisabledFile(msg); didSelect {
			// Let's clear the selectedFile and display an error.
			m.err = errors.New(path + " is not valid.")
			m.selectedFile = ""
			cmds = append(cmds, clearErrorAfter(2*time.Second))
		}

	case boxView:
		m.boxView, cmd = m.boxView.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var (
	EditorBoxSize = lipgloss.NewStyle().
			Width(50).
			Height(17).
			PaddingTop(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder())

	ActiveEditor   = lipgloss.NewStyle().Inherit(EditorBoxSize).BorderForeground(LightGreenish)
	InactiveEditor = lipgloss.NewStyle().Inherit(EditorBoxSize).BorderForeground(DarkerPink)
)

func (m MainModel) View() string {
	var help string
	var active int

	w := lipgloss.Width

	switch m.state {
	case growthEdit, attackEdit, evscEdit,
		miscEdit, originEdit, statEdit, ribbonEdit:
		help = m.editors[m.state].GetKeys()
		active = int(m.state) - 1
		m.UpdateEditorValues()

	case mailEdit:
		r := m.editors[mailEdit].View()
		t := lipgloss.JoinVertical(
			lipgloss.Center,
			r,
			m.editors[mailEdit].GetKeys(),
			m.MakeHeader(w(r)),
		)
		return t

	case searcher:
		return m.searcher.View()
	case picker:
		return m.makeFilePickerView()
	case boxView:
		return m.boxView.View()
	}

	views := make([]string, len(m.editors)-1)
	for i := range len(views) {
		views[i] = m.editors[i+1].View()
	}

	for i := range views {
		if i != active {
			if i == ribbonEdit-1 {
				views[i] = InactiveEditor.Width(102).BorderTop(false).Render(views[i])
			} else {
				views[i] = InactiveEditor.Render(views[i])
			}
		} else {
			if i == ribbonEdit-1 {
				views[i] = ActiveEditor.Width(102).BorderTop(false).Render(views[i])
			} else {
				views[i] = ActiveEditor.Render(views[i])
			}
		}
	}

	if m.state >= statEdit {
		misc := lipgloss.JoinHorizontal(lipgloss.Center, views[4], views[5])
		out := lipgloss.JoinVertical(lipgloss.Center, misc, views[6])
		return lipgloss.JoinVertical(lipgloss.Center, out, help, m.MakeHeader(w(out)))
	} else {
		GA := lipgloss.JoinHorizontal(lipgloss.Center, views[0], views[1])
		EM := lipgloss.JoinHorizontal(lipgloss.Center, views[2], views[3])
		GAEM := lipgloss.JoinVertical(lipgloss.Center, GA, EM)
		return lipgloss.JoinVertical(lipgloss.Center, GAEM, help, m.MakeHeader(w(GAEM)))
	}
}

func (m *MainModel) UpdateEditorValues() {
	for i := range m.editors {
		m.editors[i].UpdateValues()
	}
}

func (m *MainModel) MakeHeader(w int) string {
	editor := edit[m.state]

	width := lipgloss.Width

	xkey := XKEYStyle.Render(fmt.Sprintf("XKEY: %08X", m.pks.GetEncryptionKey()))

	TID, SID := m.pks.GetTIDSIDComboPair()
	tid := TIDStyle.Render(fmt.Sprintf("TID: %04d", TID))
	sid := SIDStyle.Render(fmt.Sprintf("SID: %04d", SID))

	nick := NickNameStyle.Render("Nickname:", m.pks.GetCleanNickname())
	ot := OTNameStyle.Render("OT:", m.pks.GetCleanOTName())

	filler := lipgloss.NewStyle().
		Background(DarkBlueish).
		Foreground(RegularText).
		Width(w - width(xkey) - width(tid) - width(sid) - width(nick) - width(ot) - width(editor))

	out := lipgloss.JoinHorizontal(
		lipgloss.Center,
		xkey,
		tid,
		sid,
		nick,
		ot,
		filler.Render(" "+m.status),
		EditStyle.Render(editor),
	)

	return out
}

func (m *MainModel) makeFilePickerView() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.picker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString(m.picker.CurrentDirectory)
	}
	s.WriteString("\n" + m.picker.View())
	return s.String()
}

func subStructArray2String(s []byte) string {
	str := new(strings.Builder)
	for i := range s {
		str.WriteString(fmt.Sprintf("%02X ", s[i]))
	}
	return str.String()
}

func joinGrowthFieldValueString(f, v string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		GrowthFieldStyle.Render(f),
		GrowthValueStyle.Render(v),
	)
}

func joinConditionFieldValueString(f, v string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		ConditionFieldStyle.Render(f),
		ConditionValueStyle.Render(v),
	)
}

func joinSearchValueStrings(f, v string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		SearchFieldStyle.Render(f),
		SearchValueStyle.Render(v),
	)
}
