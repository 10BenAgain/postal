package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	MonView   = 42
	CheckMark = "✓"
	XMark     = "✘"
)

// https://lospec.com/palette-list/soapy-10
var (
	LightGreenish = lipgloss.Color("#54cea7")
	LightBlueish  = lipgloss.Color("#2ba4a6")
	MidBlueish    = lipgloss.Color("#0c6987")
	DarkerBlueish = lipgloss.Color("#054b84")
	DarkBlueish   = lipgloss.Color("#0d2147")
	LightPink     = lipgloss.Color("#ffb0bf")
	DarkerPink    = lipgloss.Color("#ff82bd")
	LightPurple   = lipgloss.Color("#d74ac7")
	DarkerPurple  = lipgloss.Color("#a825ba")
	DarkPurple    = lipgloss.Color("#682b9c")
	RegularText   = lipgloss.Color("#E5E5E5")
	HelperText    = lipgloss.Color("#bab8b1")
	SubText       = lipgloss.Color("#949494")
)

var (
	SubStructBlock = lipgloss.NewStyle().Width(MonView).Padding(0, 1, 1, 0)
	SubStructByte  = lipgloss.NewStyle().Foreground(SubText).Align(lipgloss.Center).Width(MonView).Padding(0, 1)
	SubMovesBlock  = lipgloss.NewStyle().Width(MonView).MarginLeft(2).Padding(0, 0, 1, 0)
)

var (
	MSpecies      = lipgloss.NewStyle().Bold(true).Align(lipgloss.Left).Width(16).Background(MidBlueish).Foreground(LightPink).Padding(0, 1)
	MOTID         = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Background(DarkPurple).Width(10).Foreground(LightPink).Padding(0, 1)
	MPID          = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Background(DarkerBlueish).Foreground(DarkerPink).Width(10)
	StructOrder   = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Background(DarkBlueish).Foreground(LightGreenish).Width(6).Padding(0, 1)
	OuterBorder   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(DarkerPurple)
	MenuStyle     = lipgloss.NewStyle().Inherit(OuterBorder).Padding(0, 1).Align(lipgloss.Center)
	MailMenuStyle = MenuStyle.BorderForeground(LightBlueish).Border(lipgloss.DoubleBorder())
)

var (
	ListHeader        = lipgloss.NewStyle().Foreground(DarkerPink).MarginRight(5).Bold(true)
	SubListHeader     = lipgloss.NewStyle().Foreground(RegularText).Padding(0, 2, 0, 0).Bold(true)
	SearchHeaderStyle = lipgloss.NewStyle().Foreground(SubText).Align(lipgloss.Center)

	KeyValue = lipgloss.NewStyle().Foreground(HelperText)
	KeyPair  = lipgloss.NewStyle().Foreground(SubText)
)

var (
	GrowthFieldStyle    = lipgloss.NewStyle().Bold(true).Width(12).Foreground(DarkerPink).MarginLeft(1)
	GrowthValueStyle    = lipgloss.NewStyle().Bold(true).Width(22).Foreground(RegularText).MarginLeft(1).Align(lipgloss.Left)
	SearchFieldStyle    = lipgloss.NewStyle().Bold(true).Width(10).Foreground(DarkerPink)
	SearchValueStyle    = lipgloss.NewStyle().Bold(true).Width(12).Align(lipgloss.Left)
	ConditionFieldStyle = lipgloss.NewStyle().Inherit(GrowthFieldStyle).Width(10)
	ConditionValueStyle = lipgloss.NewStyle().Inherit(GrowthValueStyle).Width(6)

	FocusedStyle   = lipgloss.NewStyle().Foreground(DarkPurple)
	CatStyle       = lipgloss.NewStyle().Foreground(LightPink)
	FocusedText    = lipgloss.NewStyle().Foreground(RegularText)
	WordEntryStyle = lipgloss.NewStyle().Foreground(LightPink).Padding(1, 2, 1, 2)
)

var (
	RibbonListBlock      = lipgloss.NewStyle().Padding(1)
	RibbonXMarkEnumStyle = lipgloss.NewStyle().Foreground(LightPink).Bold(true)
	RibbonCheckEnumStyle = lipgloss.NewStyle().Foreground(LightGreenish).Bold(true)
	RibbonSumStyle       = lipgloss.NewStyle().Foreground(SubText).Align(lipgloss.Center).Underline(true)
	StatsWrapperStyle    = lipgloss.NewStyle().MarginLeft(8)
	OuterRibbon          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(LightPurple).Padding(1)
)

var (
	XKEYStyle     = lipgloss.NewStyle().Background(DarkerPurple).Foreground(LightPink).Padding(0, 1)
	TIDStyle      = lipgloss.NewStyle().Background(DarkerBlueish).Foreground(LightPink).Padding(0, 1)
	SIDStyle      = lipgloss.NewStyle().Background(LightPurple).Foreground(RegularText).Padding(0, 1)
	NickNameStyle = lipgloss.NewStyle().Background(DarkPurple).Foreground(RegularText).Padding(0, 1)
	EditStyle     = lipgloss.NewStyle().Background(DarkerBlueish).Foreground(LightGreenish).Padding(0, 1)
	OTNameStyle   = lipgloss.NewStyle().Background(DarkerPurple).Foreground(RegularText).Padding(0, 1)
	edit          = []string{"Mail", "Growth", "Attacks", "EV&C", "Misc", "Stats", "Origin", "Ribbons", "Lookup"}
)

var (
	blank           = "-"
	WordEntry       = "_ _ _ _ _ _ _ _"
	BlankWordValues = []string{"0x0", "0x0", "0x0", "0x0"}

	StatNames        = []string{"HP", "Atk", "Def", "Spe", "Spa", "Spd"}
	ConditionNames   = []string{"Coolness", "Beauty", "Cuteness", "Smartness", "Toughness", "Feel"}
	GrowthBlockNames = []string{"Species", "Item", "EXP", "PP 1-4", "Friendship"}
	MiscBlockNames   = []string{"PKRS", "Met Loc", "Origin", "Stats", "Ribbons"}
	OriginNames      = []string{"Met Level", "Origin Game", "Ball", "OT Gender"}

	RibbonBlockOne   = []string{"Cool", "Beauty", "Cute", "Smart", "Tough"}
	RibbonBlockTwo   = []string{"Champion", "Winning", "Victory", "Artist", "Effort"}
	RibbonBlockThree = []string{"Battle", "Regional", "National"}
	RibbonBlockFour  = []string{"Country", "National", "Earth", "World"}
	RibbonBlockFive  = []string{"Fateful Encounter"}
)
