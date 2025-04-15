package pokemon

import (
	vals "postal/game"
	"postal/utils"
)

// These are just convenient wrappers that provide errors
// when searching a list or map that is commonly used

func NumToSpeciesName(n uint16) (string, error) {
	if n > SpeciesIndexMax {
		return "", utils.ErrOutOfRange
	}

	return vals.SpeciesDataList[n].Name, nil
}

func NumToMoveName(n uint16) (string, error) {
	if n > MoveIndexMax {
		return "", utils.ErrOutOfRange
	}

	return vals.MovesList[n], nil
}

func NumToItemName(n uint16) (string, error) {
	if val, ok := vals.Items[n]; ok {
		return val, nil
	}

	return "", utils.ErrOutOfRange
}

func NumToEVVals(n uint16) (uint8, uint8) {
	return uint8(n & 0xFF), uint8((n & 0xFF00) >> 8)
}

func CombineU16Split(hi, lo uint16) uint32 {
	return uint32(hi)<<16 | uint32(lo)
}
