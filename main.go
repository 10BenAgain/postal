package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"postal/pokemon"
	"postal/save"
	"postal/tui"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	pk3 = iota
	ek3
	sav
	None
)

func determineExtension(s string) int {
	switch path.Ext(s) {
	case ".pk3":
		return pk3
	case ".ek3":
		return ek3
	case ".sav", ".SAV":
		return sav
	default:
		return None
	}
}

func main() {
	clearScreen()

	args := os.Args[1:]

	if len(args) == 1 {
		f := args[0]
		ext := determineExtension(f)

		switch ext {

		case pk3:
			if pks, err := getPK3FromFile(f); err != nil {
				fmt.Printf("err: %v\n", err)
				os.Exit(1)
			} else {
				p := tea.NewProgram(tui.InitMainEditor(pks))
				if _, err := p.Run(); err != nil {
					os.Exit(1)
				}
			}

		case ek3:
			if pks, err := getEK3FromFile(f); err != nil {
				fmt.Printf("err: %v\n", err)
				os.Exit(1)
			} else {
				p := tea.NewProgram(tui.InitMainEditor(pks))
				if _, err := p.Run(); err != nil {
					os.Exit(1)
				}
			}

		case sav:
			sav := getDataFromPath(f)
			p := tea.NewProgram(tui.InitMainBoxData(sav))
			if _, err := p.Run(); err != nil {
				os.Exit(1)
			}

		default:
			p := tea.NewProgram(tui.InitMainBlank())
			if _, err := p.Run(); err != nil {
				os.Exit(1)
			}
		}

	} else {
		p := tea.NewProgram(tui.InitMainBlank())
		if _, err := p.Run(); err != nil {
			os.Exit(1)
		}
	}
}

func getDataFromPath(p string) save.RawBoxDataTotal {
	sav := save.GenerateRawSaveData(p)
	block := sav.GenerateSaveBlock()

	return block.GetRawBoxData()
}

func getPK3FromFile(p string) (pokemon.PStructure, error) {
	pk, err := save.GetRawMonDataFromFile(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return pokemon.PStructure{}, err
	} else {
		return pokemon.GeneratePokemonFromRawData(pk, false, true), nil
	}
}

func getEK3FromFile(p string) (pokemon.PStructure, error) {
	pk, err := save.GetRawMonDataFromFile(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return pokemon.PStructure{}, err
	} else {
		return pokemon.GeneratePokemonFromRawData(pk, false, false), nil
	}
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
