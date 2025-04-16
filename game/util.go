package game

import (
	"maps"
	"postal/utils"
)

var (
	RevRegularWordMap  = reversedU16ValueMap(RegularWordMap)
	RevWordPokemon1Map = reversedU16ValueMap(Pokemon1Map)
	RevWordPokemon2Map = reversedU16ValueMap(Pokemon2Map)
	RevWordMoves1Map   = reversedU16ValueMap(Moves1Map)
	RevWordMoves2Map   = reversedU16ValueMap(Moves2Map)
	RevItemsMap        = reversedU16ValueMap(Items)
	RevLocationMap     = reversedU8ValueMap(LocationMap)
	RevBallMap         = reversedU16ValueMap(Balls)
	RevOriginGameMap   = reversedU8ValueMap(OriginGame)

	MergedWordMap    = mergeWordMaps()
	ReversedMergeMap = reversedU16ValueMap(MergedWordMap)
)

func mergeWordMaps() map[uint16]string {
	merge := make(map[uint16]string)
	maps.Copy(merge, RegularWordMap)
	maps.Copy(merge, Pokemon1Map)
	maps.Copy(merge, Pokemon2Map)
	maps.Copy(merge, Moves1Map)
	maps.Copy(merge, Moves2Map)
	return merge
}

func reversedU16ValueMap(m map[uint16]string) map[string]uint16 {
	rev := make(map[string]uint16, len(m))

	for k, v := range m {
		rev[v] = k
	}
	return rev
}

func reversedU8ValueMap(m map[uint8]string) map[string]uint8 {
	rev := make(map[string]uint8, len(m))

	for k, v := range m {
		rev[v] = k
	}
	return rev
}

func MoveLookup(s string) (int, error) {
	for _, v := range MoveDataList {
		if s == v.Name {
			return v.Index, nil
		}
	}
	return 0, utils.ErrOutOfRange
}

func MonLookup(s string) (int, error) {
	for _, v := range SpeciesDataList {
		if s == v.Name {
			return v.Dex, nil
		}
	}

	return 0, utils.ErrOutOfRange
}

func WordMon1Lookup(s string) (uint16, error) {
	val, ok := RevWordPokemon1Map[s]

	if ok {
		return val, nil
	}
	return val, utils.ErrOutOfRange
}

func WordMon2Lookup(s string) (uint16, error) {
	val, ok := RevWordPokemon2Map[s]

	if ok {
		return val, nil
	}
	return val, utils.ErrOutOfRange
}

func WordLookup(s string) (uint16, error) {
	val, ok := ReversedMergeMap[s]

	if ok {
		return val, nil
	}

	return val, utils.ErrOutOfRange
}

func ItemLookup(s string) (int, error) {
	val, ok := RevItemsMap[s]
	if ok {
		return int(val), nil
	}

	val, ok = DoubleCapItems[s]
	if ok {
		return int(val), nil
	}

	val, ok = OneLinerItems[s]
	if ok {
		return int(val), nil
	}

	return 0, utils.ErrOutOfRange
}

func BallLookup(s string) (uint16, error) {
	val, ok := RevBallMap[s]

	if ok {
		return val, nil
	}

	return val, utils.ErrOutOfRange
}

func OriginGameLookup(s string) (uint8, error) {
	val, ok := RevOriginGameMap[s]

	if ok {
		return val, nil
	}

	return val, utils.ErrOutOfRange
}

func MetLocLookup(s string) (uint8, error) {
	val, ok := RevLocationMap[s]

	if ok {
		return val, nil
	}

	return 0, utils.ErrOutOfRange
}

func FindMonList1WordPair(OTID, PID uint32, order uint32) [][2]uint16 {
	// Given XKey where XKey = OTID ^ PID, find every word combination that
	// Word 1 ^ Word 2 == Upper half of XKey (XKey >> 16)
	// And the resulting PID order result matches input order

	XKey := uint16((OTID ^ PID) >> 16)
	PIDLower := PID & 0xFFFF

	var pairs [][2]uint16

	for i := range Pokemon1Map {
		p0 := PIDLower | uint32(i)<<16

		if p0%24 != order {
			continue
		}

		for j := range Pokemon1Map {
			p1 := PIDLower | uint32(j)<<16
			if i^j == XKey {
				pairs = append(pairs, [2]uint16{i, j})
				if p1%24 == order {
					pairs = append(pairs, [2]uint16{j, i})
				}
			}
		}
	}

	return pairs
}
