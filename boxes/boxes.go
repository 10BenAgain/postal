package boxes

import (
	"postal/pokemon"
	"postal/save"
	"postal/utils"
)

const (
	BoxMonCount int = 30
	BoxMonSize      = 0x50
	BoxSize         = BoxMonCount * BoxMonSize
)

type (
	PCBoxBufferMash []byte
	PCBoxMonData    []byte
)

type RWPCBox struct {
	Num       int
	Name      string
	Wallpaper int
	Mons      [30]pokemon.PStructure
}

func GeneratePCBoxBufferMash(t save.RawBoxDataTotal) PCBoxBufferMash {
	var sum []byte

	for i := range t.Boxes {
		sum = append(sum, t.Boxes[i].Data...)
	}
	return sum
}

func GetPCBoxNameString(index int, t save.RawBoxDataTotal) string {
	nameOffset := index * 0x9
	rn := utils.GetSliceFromRawData(t.BoxNames, nameOffset, 0x9)
	return utils.WesternSliceToString(rn)
}

func GetPCBoxWallpaper(index int, t save.RawBoxDataTotal) int {
	return int(utils.GetSliceFromRawData(t.WallPapers, index, 1)[0])
}

func (b PCBoxBufferMash) GenerateBoxMonData(index int) PCBoxMonData {
	bxOffset := index * BoxSize
	return utils.GetSliceFromRawData(b, bxOffset, BoxSize)
}

func (b PCBoxBufferMash) GetMonAtBoxSlot(box int, slot int) []byte {
	if box > 14 || slot >= BoxMonCount {
		return nil
	}

	return utils.GetSliceFromRawData(b, box*BoxSize+slot*BoxMonSize, BoxMonSize)
}

func (b PCBoxMonData) GetMonDataFromSlot(index int) []byte {
	if index >= BoxMonCount {
		return nil
	}

	return utils.GetSliceFromRawData(b, BoxMonSize*index, BoxMonSize)
}

func GenerateRWPCBox(index int, t save.RawBoxDataTotal, mondat PCBoxMonData) RWPCBox {
	var outBox RWPCBox
	outBox.Num = index + 1
	outBox.Name = GetPCBoxNameString(index, t)
	outBox.Wallpaper = GetPCBoxWallpaper(index, t)

	for i := range BoxMonCount {
		monOffset := i * BoxMonSize
		outBox.Mons[i] = pokemon.GeneratePokemonFromRawData(
			utils.GetSliceFromRawData(mondat, monOffset, BoxMonSize),
			true,
			false,
		)
	}
	return outBox
}
