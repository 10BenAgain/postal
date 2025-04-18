package pokemon

import (
	"encoding/binary"
	vals "postal/game"
	"postal/utils"
	"reflect"
	"strings"
)

const (
	OriginMax              = 6
	LanguageMax            = 7
	TrainerNameMaxLength   = 7
	NickNameMaxLength      = 10
	SubStructureLineLength = 12
	AbilityIndexMax        = 76
	LocationMax            = 255
	MoveIndexMax           = 354
	ItemIndexMax           = 376
	SpeciesIndexMax        = 413

	BadEggMask     uint8 = 0x1
	HasSpeciesMask uint8 = 0x2
	EggNameMask    uint8 = 0x4
	BlockBoxMask   uint8 = 0x8

	MetLevelMask   uint16 = 0b0000_0000_0111_1111
	GameOriginMask uint16 = 0b0000_0111_1000_0000
	BallMask       uint16 = 0b0111_1000_0000_0000
	TGenderMask    uint16 = 0b1000_0000_0000_0000

	PPBonusMask0 uint8 = 0b0000_0011
	PPBonusMask1 uint8 = 0b0000_1100
	PPBonusMask2 uint8 = 0b0011_0000
	PPBonusMask3 uint8 = 0b1100_0000

	HPMask  uint32 = 0x1F
	AtkMask        = HPMask << 5
	DefMask        = AtkMask << 5
	SpeMask        = DefMask << 5
	SpaMask        = SpeMask << 5
	SpdMask        = SpaMask << 5
	EggMask        = SpdMask << 1
	AblMask        = EggMask << 1

	CoolMask             uint32 = 0x7
	BeautyMask                  = CoolMask << 3
	CuteMask                    = BeautyMask << 3
	SmartMask                   = CuteMask << 3
	ToughMask                   = SmartMask << 3
	ChampionMask         uint32 = 0x8000
	WinningMask                 = ChampionMask << 1
	VictoryMask                 = WinningMask << 1
	ArtistMask                  = VictoryMask << 1
	EffortMask                  = ArtistMask << 1
	BattleChampionMask          = EffortMask << 1
	RegionalChampionMask        = BattleChampionMask << 1
	NationalChampionMask        = RegionalChampionMask << 1
	CountryMask                 = NationalChampionMask << 1
	NationalMask                = CountryMask << 1
	EarthMask                   = NationalMask << 1
	WorldMask                   = EarthMask << 1
	FatefulMask                 = WorldMask << 4
)

type (
	PStructure struct {
		PID      uint32
		OTID     uint32
		NickName string
		Lang     uint8
		Flags    uint8
		OTName   string
		Markings uint8
		Checksum uint16
		filler   uint16
		Sub0     SubStruct0
		Sub1     SubStruct1
		Sub2     SubStruct2
		Sub3     SubStruct3
		Status   uint32
		Level    uint8
		MailID   uint8
		CurHP    uint16
		TotalHP  uint16
		Atk      uint16
		Def      uint16
		Spe      uint16
		Spa      uint16
		Spd      uint16
	}

	SubStruct0 struct {
		Species    uint16
		HeldItem   uint16
		Experience uint32
		PPBonuses  uint8
		FriendShip uint8
		filler     uint16
	}

	SubStruct1 struct {
		Moves [4]uint16
		PP    [4]uint8
	}

	SubStruct2 struct {
		HpEV    uint8
		AtkEV   uint8
		DefEV   uint8
		SpeEV   uint8
		SpAtkEV uint8
		SpDefEV uint8
		Cool    uint8
		Beauty  uint8
		Cute    uint8
		Smart   uint8
		Tough   uint8
		Feel    uint8
	}

	SubStruct3 struct {
		Pokerus                uint8
		MetLocation            uint8
		MetLevel               uint8
		MetGame                uint8
		Ball                   uint8
		OTGender               uint8
		HpIV                   uint8
		AtkIV                  uint8
		DefIV                  uint8
		SpeIV                  uint8
		SpAtkIV                uint8
		SpDIV                  uint8
		IsEgg                  uint8
		AbilityNum             uint8
		CoolRibbon             uint8
		BeautyRibon            uint8
		CuteRibbon             uint8
		SmartRibbon            uint8
		ToughRibbon            uint8
		ChampionRibbon         uint8
		WinningRibbon          uint8
		VictoryRibbon          uint8
		ArtistRibbon           uint8
		EffortRibbon           uint8
		BattleChampionRibbon   uint8
		RegionalChampionRibbon uint8
		NationalChampionRibbon uint8
		CountryRibbon          uint8
		NationalRibbon         uint8
		EarthRibbon            uint8
		WorldRibbon            uint8
		UnusedRibbon           uint8
		FatefulEnc             uint8
	}

	MonMiscFlags struct {
		IsBadEgg   bool
		HasSpecies bool
		UseEggName bool
		BlockRSBox bool
	}
)

var (
	ExperienceMap = map[int][]uint{
		0: GetMediumFastTable(),
		1: GetErraticTable(),
		2: GetFluctuatingTable(),
		3: GetMediumSlowTable(),
		4: GetFastTable(),
		5: GetSlowTable(),
	}

	SubStructOffsetMap = map[int][2]uint8{
		0: {0x20, 0x2C},
		1: {0x2C, 0x38},
		2: {0x38, 0x44},
		3: {0x44, 0x50},
	}
)

func (p *PStructure) GetMonSubStructOrder() []int {
	subStructOrders := [][]int{
		{0, 1, 2, 3}, // GAEM
		{0, 1, 3, 2}, // GAME
		{0, 2, 1, 3}, // GEAM
		{0, 2, 3, 1}, // GEMA
		{0, 3, 1, 2}, // GMAE
		{0, 3, 2, 1}, // GMEA
		{1, 0, 2, 3}, // AGEM
		{1, 0, 3, 2}, // AGME
		{1, 2, 0, 3}, // AEGM
		{1, 2, 3, 0}, // AEMG
		{1, 3, 0, 2}, // AMGE
		{1, 3, 2, 0}, // AMEG
		{2, 0, 1, 3}, // EGAM
		{2, 0, 3, 1}, // EGMA
		{2, 1, 0, 3}, // EAGM
		{2, 1, 3, 0}, // EAMG
		{2, 3, 0, 1}, // EMGA
		{2, 3, 1, 0}, // EMAG
		{3, 0, 1, 2}, // MGAE
		{3, 0, 2, 1}, // MGEA
		{3, 1, 0, 2}, // MAGE
		{3, 1, 2, 0}, // MAEG
		{3, 2, 0, 1}, // MEGA
		{3, 2, 1, 0}, // MEAG
	}

	return subStructOrders[p.PID%24]
}

func GeneratePokemonFromRawData(data []byte, isBoxMon bool, isDecrypted bool) PStructure {
	var mon PStructure

	mon.PID = binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x0, 4))
	mon.OTID = binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x4, 4))
	nick := utils.GetSliceFromRawData(data, 0x8, 10)

	empty := true
	for _, v := range nick {
		if v != 0 {
			empty = false
			break
		}
	}

	if !empty {
		mon.NickName = utils.GenerateNickNameString(nick)
		empty = false
	}

	otName := utils.GetSliceFromRawData(data, 0x14, 7)

	for _, v := range otName {
		if v != 0 {
			empty = false
			break
		}
	}

	if !empty {
		mon.OTName = utils.GenerateOTNameString(otName)
	}

	mon.Lang = utils.GetSliceFromRawData(data, 0x12, 1)[0]
	mon.Flags = utils.GetSliceFromRawData(data, 0x13, 1)[0]
	mon.Markings = utils.GetSliceFromRawData(data, 0x1B, 1)[0]
	mon.Checksum = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x1C, 2))
	mon.filler = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x1E, 2))

	subStructOrder := mon.GetMonSubStructOrder()
	subStructData := GetRawSubStructsFromMon(data)

	if !isDecrypted {
		subStructData = CryptMonSubstructs(
			subStructData,
			mon.PID^mon.OTID,
		)
	}

	{
		pos := 0
		if !isDecrypted {
			pos = FindSubstructPosition(subStructOrder, 0)
		}

		data := utils.GetSliceFromRawData(
			subStructData,
			pos*SubStructureLineLength,
			SubStructureLineLength,
		)
		mon.Sub0 = GenerateSubStructure0(data)
	}

	{
		pos := 1
		if !isDecrypted {
			pos = FindSubstructPosition(subStructOrder, 1)
		}

		data := utils.GetSliceFromRawData(
			subStructData,
			pos*SubStructureLineLength,
			SubStructureLineLength,
		)
		mon.Sub1 = GenerateSubStructure1(data)
	}

	{
		pos := 2
		if !isDecrypted {
			pos = FindSubstructPosition(subStructOrder, 2)
		}

		data := utils.GetSliceFromRawData(
			subStructData,
			pos*SubStructureLineLength,
			SubStructureLineLength,
		)
		mon.Sub2 = GenerateSubStructure2(data)
	}

	{
		pos := 3
		if !isDecrypted {
			pos = FindSubstructPosition(subStructOrder, 3)
		}

		data := utils.GetSliceFromRawData(
			subStructData,
			pos*SubStructureLineLength,
			SubStructureLineLength,
		)
		mon.Sub3 = GenerateSubStructure3(data)
	}

	if !isBoxMon {
		mon.Status = binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x50, 4))
		mon.Level = utils.GetSliceFromRawData(data, 0x54, 1)[0]
		mon.MailID = utils.GetSliceFromRawData(data, 0x55, 1)[0]
		mon.CurHP = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x56, 2))
		mon.TotalHP = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x58, 2))
		mon.Atk = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x5A, 2))
		mon.Def = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x5C, 2))
		mon.Spe = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x5E, 2))
		mon.Spa = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x60, 2))
		mon.Spd = binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x62, 2))
	} else {
		CalculateAllStats(&mon)
	}

	return mon
}

func GetRawSubStructsFromMon(data []byte) []byte {
	return utils.GetSliceFromRawData(data, 0x20, 0x30)
}

func GenerateSubStructure0(data []byte) SubStruct0 {
	if len(data) != 12 {
		return SubStruct0{}
	}

	return SubStruct0{
		Species:    binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x0, 2)),
		HeldItem:   binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x2, 2)),
		Experience: binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x4, 4)),
		PPBonuses:  utils.GetSliceFromRawData(data, 0x8, 1)[0],
		FriendShip: utils.GetSliceFromRawData(data, 0x9, 1)[0],
		filler:     binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0xA, 2)),
	}
}

func GenerateSubStructure1(data []byte) SubStruct1 {

	if len(data) != 12 {
		return SubStruct1{}
	}

	moves := [4]uint16{
		binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x0, 2)),
		binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x2, 2)),
		binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x4, 2)),
		binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x6, 2)),
	}

	pp := [4]uint8{
		utils.GetSliceFromRawData(data, 0x8, 1)[0],
		utils.GetSliceFromRawData(data, 0x9, 1)[0],
		utils.GetSliceFromRawData(data, 0xA, 1)[0],
		utils.GetSliceFromRawData(data, 0xB, 1)[0],
	}

	return SubStruct1{
		Moves: moves,
		PP:    pp,
	}
}

func GenerateSubStructure2(data []byte) SubStruct2 {
	if len(data) != 12 {
		return SubStruct2{}
	}

	return SubStruct2{
		HpEV:    utils.GetSliceFromRawData(data, 0x0, 1)[0],
		AtkEV:   utils.GetSliceFromRawData(data, 0x1, 1)[0],
		DefEV:   utils.GetSliceFromRawData(data, 0x2, 1)[0],
		SpeEV:   utils.GetSliceFromRawData(data, 0x3, 1)[0],
		SpAtkEV: utils.GetSliceFromRawData(data, 0x4, 1)[0],
		SpDefEV: utils.GetSliceFromRawData(data, 0x5, 1)[0],
		Cool:    utils.GetSliceFromRawData(data, 0x6, 1)[0],
		Beauty:  utils.GetSliceFromRawData(data, 0x7, 1)[0],
		Cute:    utils.GetSliceFromRawData(data, 0x8, 1)[0],
		Smart:   utils.GetSliceFromRawData(data, 0x9, 1)[0],
		Tough:   utils.GetSliceFromRawData(data, 0xA, 1)[0],
		Feel:    utils.GetSliceFromRawData(data, 0xB, 1)[0],
	}
}

func GenerateSubStructure3(data []byte) SubStruct3 {
	if len(data) != 12 {
		return SubStruct3{}
	}

	origin := binary.LittleEndian.Uint16(utils.GetSliceFromRawData(data, 0x2, 2))
	stats := binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x4, 4))
	ribbons := binary.LittleEndian.Uint32(utils.GetSliceFromRawData(data, 0x8, 4))

	return SubStruct3{
		Pokerus:     utils.GetSliceFromRawData(data, 0x0, 1)[0],
		MetLocation: utils.GetSliceFromRawData(data, 0x1, 1)[0],
		MetLevel:    uint8(origin & MetLevelMask),
		MetGame:     uint8(origin & GameOriginMask >> 7),
		Ball:        uint8(origin & BallMask >> 11),
		OTGender:    uint8(origin & TGenderMask >> 14),

		HpIV:       uint8(stats & HPMask),
		AtkIV:      uint8(stats & AtkMask >> 5),
		DefIV:      uint8(stats & DefMask >> 10),
		SpeIV:      uint8(stats & SpeMask >> 15),
		SpAtkIV:    uint8(stats & SpaMask >> 20),
		SpDIV:      uint8(stats & SpdMask >> 25),
		IsEgg:      uint8(stats & EggMask >> 30),
		AbilityNum: uint8(stats & AblMask >> 31),

		CoolRibbon:             uint8(ribbons & CoolMask),
		BeautyRibon:            uint8(ribbons & BeautyMask >> 3),
		CuteRibbon:             uint8(ribbons & CuteMask >> 6),
		SmartRibbon:            uint8(ribbons & SmartMask >> 9),
		ToughRibbon:            uint8(ribbons & ToughMask >> 12),
		ChampionRibbon:         uint8(ribbons & ChampionMask >> 15),
		WinningRibbon:          uint8(ribbons & WinningMask >> 16),
		VictoryRibbon:          uint8(ribbons & VictoryMask >> 17),
		ArtistRibbon:           uint8(ribbons & AblMask >> 18),
		EffortRibbon:           uint8(ribbons & EffortMask >> 19),
		BattleChampionRibbon:   uint8(ribbons & BattleChampionMask >> 20),
		RegionalChampionRibbon: uint8(ribbons & RegionalChampionMask >> 21),
		NationalChampionRibbon: uint8(ribbons & NationalChampionMask >> 22),
		CountryRibbon:          uint8(ribbons & CountryMask >> 23),
		NationalRibbon:         uint8(ribbons & NationalMask >> 24),
		EarthRibbon:            uint8(ribbons & EarthMask >> 25),
		WorldRibbon:            uint8(ribbons & WorldMask >> 26),
		FatefulEnc:             uint8(ribbons & FatefulMask >> 31),
	}
}

func (p *PStructure) GetEncryptionKey() uint32 {
	return p.PID ^ p.OTID
}

func calcHPStat(level uint8, base uint, evs uint8, iv uint8) uint16 {
	evc := uint(evs) / 4
	return uint16((((2*base + uint(iv) + evc) * uint(level)) / 100) + uint(level) + 10)
}

func calcOtherStat(level uint8, base uint, nat float32, evs uint8, iv uint8) uint16 {
	evc := uint(evs) / 4
	p1 := (((2*base + uint(iv) + evc) * uint(level)) / 100) + 5
	return uint16(float32(p1) * nat)
}

func DetermineLevelFromExperience(exp uint, tbIndex int) int {
	table := ExperienceMap[tbIndex]
	l := 0
	for range table {
		if exp >= table[l] {
			l++
		}
	}
	return l
}

func CalculateLevel(pk *PStructure) {
	species := pk.Sub0.Species
	var xpr uint
	if species >= SpeciesIndexMax {
		xpr = vals.SpeciesDataList[SpeciesIndexMax].ExpRate
	} else {
		xpr = vals.SpeciesDataList[species].ExpRate
	}
	pk.Level = uint8(DetermineLevelFromExperience(
		uint(pk.Sub0.Experience),
		int(xpr),
	),
	)
}

func CalculateAllStats(pk *PStructure) {

	if pk.Level == 0 && pk.Sub3.IsEgg == 0 {
		CalculateLevel(pk)
	}

	level := pk.Level
	species := pk.Sub0.Species

	var baseStats []uint
	if species >= SpeciesIndexMax {
		baseStats = vals.SpeciesDataList[SpeciesIndexMax].BaseStats
	} else {
		baseStats = vals.SpeciesDataList[species].BaseStats
	}

	natureVals, _ := pk.GetNature()

	evSlice := []uint8{pk.Sub2.HpEV, pk.Sub2.AtkEV, pk.Sub2.DefEV,
		pk.Sub2.SpAtkEV, pk.Sub2.SpDefEV, pk.Sub2.SpeEV,
	}

	ivSlice := []uint8{pk.Sub3.HpIV, pk.Sub3.AtkIV, pk.Sub3.DefIV,
		pk.Sub3.SpeIV, pk.Sub3.SpDIV, pk.Sub3.SpeIV,
	}

	if level > 0 {
		pk.TotalHP = calcHPStat(level, baseStats[0], evSlice[0], ivSlice[0])
	}

	remStats := []uint16{pk.Atk, pk.Def, pk.Spa, pk.Spd, pk.Spe}
	for i := range remStats {
		remStats[i] = calcOtherStat(
			level,
			baseStats[i+1],
			vals.NatureStatModifier[natureVals][i],
			evSlice[i+1],
			ivSlice[i+1],
		)
	}
}

func CryptMonSubstructs(data []byte, xkey uint32) []byte {
	// TODO: Return dummy data
	rs := make([]byte, len(data))
	s := make([]byte, 4)

	for i := 0; i < len(data); i += 4 {
		sec := binary.LittleEndian.Uint32(
			utils.GetSliceFromRawData(
				data,
				i,
				4,
			),
		)
		binary.LittleEndian.PutUint32(s, sec^xkey)
		rs[i] = s[0]
		rs[i+1] = s[1]
		rs[i+2] = s[2]
		rs[i+3] = s[3]
	}
	return rs
}

func FindSubstructPosition(ord []int, ss int) int {
	for i, v := range ord {
		if ss == v {
			return i
		}
	}
	return -1
}

func (p *PStructure) InterpretFlags() MonMiscFlags {
	return MonMiscFlags{
		IsBadEgg:   p.Flags&BadEggMask != 0,
		HasSpecies: p.Flags&HasSpeciesMask != 0,
		UseEggName: p.Flags&EggNameMask != 0,
		BlockRSBox: p.Flags&BlockBoxMask != 0,
	}
}

func (p *PStructure) GetPIDSplit() (uint16, uint16) {
	hi := p.PID >> 16
	lo := p.PID & 0xFFFF
	return uint16(hi), uint16(lo)
}

func (p *PStructure) GetGender() (bool, string) {
	sp, _ := p.GetSpecies()

	ratio := vals.SpeciesDataList[sp].GenderRatio
	g := p.PID&255 > uint32(ratio)

	if g {
		return g, "Male"
	}
	return g, "Female"
}

func (p *PStructure) GetOTGender() (bool, string) {
	b := utils.Uint2Bool(p.Sub3.OTGender)

	if b {
		return b, "Female"
	}
	return b, "Male"
}

func (p *PStructure) GetTIDSIDComboPair() (uint16, uint16) {
	tid := p.OTID & 0xFFFF
	sid := p.OTID >> 16
	return uint16(tid), uint16(sid)
}

func (p *PStructure) GetShinyType() (uint8, string) {
	tid, sid := p.GetTIDSIDComboPair()
	psv := (uint32(tid) ^ uint32(sid)) ^ (p.PID >> 16) ^ (p.PID & 0xFFFF)

	switch psv {
	case 1:
		return 1, "Star"
	case 2:
		return 2, "Square"
	default:
		return 0, "None"
	}
}

func (p *PStructure) GetNature() (uint8, string) {
	nt := uint8(p.PID % 25)
	return nt, vals.Natures[nt]
}

func (p *PStructure) GetBall() (uint8, string) {
	return p.Sub3.Ball, vals.Items[uint16(p.Sub3.Ball)]
}

func (p *PStructure) GetSpecies() (uint16, string) {
	s := p.Sub0.Species
	if s > SpeciesIndexMax {
		return s, "Glitched!"
	}

	e := vals.SpeciesDataList[s]

	return s, e.Name
}

func (p *PStructure) GetAbility() (uint8, string) {
	i, _ := p.GetSpecies()
	ab := p.Sub3.AbilityNum

	if ab > AbilityIndexMax {
		return 0, vals.Abilities[0]
	}

	if i > SpeciesIndexMax {
		i = 0
	}

	s := vals.SpeciesDataList[i]

	if ab == 0 {
		return ab, vals.Abilities[uint8(s.Ab0)]
	}

	return ab, vals.Abilities[uint8(s.Ab1)]
}

func (p *PStructure) GetMoveStrings() []string {
	var ret []string
	moves := p.Sub1.Moves

	for _, v := range moves {
		if v < MoveIndexMax {
			ret = append(ret, vals.MovesList[v])
		} else {
			ret = append(ret, "Glitched!")
		}
	}
	return ret
}

func (p *PStructure) GetPPBonus() []uint {
	return []uint{
		uint(p.Sub0.PPBonuses & PPBonusMask0),
		uint(p.Sub0.PPBonuses & PPBonusMask1 >> 2),
		uint(p.Sub0.PPBonuses & PPBonusMask2 >> 4),
		uint(p.Sub0.PPBonuses & PPBonusMask3 >> 6),
	}
}

func (p *PStructure) GetHeldItem() (uint16, string) {
	if p.Sub0.HeldItem > ItemIndexMax {
		return p.Sub0.HeldItem, "Unknown"
	}

	return p.Sub0.HeldItem, vals.Items[p.Sub0.HeldItem]
}

func (p *PStructure) GetLanguage() (uint8, string) {
	if p.Lang > LanguageMax {
		return p.Lang, "Unknown"
	}
	return p.Lang, vals.LocationMap[p.Lang]
}

func (p *PStructure) GetCleanNickname() string {
	return strings.TrimSpace(strings.ReplaceAll(p.NickName, "\\", ""))
}

func (p *PStructure) GetCleanOTName() string {
	return strings.TrimSpace(strings.ReplaceAll(p.OTName, "\\", ""))
}

func (p *PStructure) GetMetLocation() (uint8, string) {
	l := p.Sub3.MetLocation
	if l > LocationMax {
		return l, "Unknown"
	}

	return l, vals.LocationMap[l]
}

func (p *PStructure) GetOriginGame() (uint8, string) {
	g := p.Sub3.MetGame

	if g > OriginMax {
		return g, "Unknown"
	}

	return g, vals.OriginGame[g]
}

func (p *PStructure) GetFillerBytes() (uint16, uint16) {
	return p.filler, p.Sub0.filler
}

func (p *PStructure) GetMonSubStructOrderString() (uint, string) {

	subStructOrders := []string{
		"GAEM", "GAME", "GEAM", "GEMA", "GMAE", "GMEA",
		"AGEM", "AGME", "AEGM", "AEMG", "AMGE", "AMEG",
		"EGAM", "EGMA", "EAGM", "EAMG", "EMGA", "EMAG",
		"MGAE", "MGEA", "MAGE", "MAEG", "MEGA", "MEAG",
	}

	i := uint(p.PID % 24)

	return i, subStructOrders[i]
}

func (p *PStructure) IsBlank() bool {
	pkr := new(PStructure)
	return reflect.DeepEqual(pkr, p)
}

func (p *PStructure) IsDecrypted() bool {
	return p.OTID == p.PID
}

func (p *PStructure) ReplaceRibbonByValue(n uint32) {
	p.Sub3.CoolRibbon = uint8(n & CoolMask)
	p.Sub3.BeautyRibon = uint8(n & BeautyMask >> 3)
	p.Sub3.CuteRibbon = uint8(n & CuteMask >> 6)
	p.Sub3.SmartRibbon = uint8(n & SmartMask >> 9)
	p.Sub3.ToughRibbon = uint8(n & ToughMask >> 12)
	p.Sub3.ChampionRibbon = uint8(n & ChampionMask >> 15)
	p.Sub3.WinningRibbon = uint8(n & WinningMask >> 16)
	p.Sub3.VictoryRibbon = uint8(n & VictoryMask >> 17)
	p.Sub3.ArtistRibbon = uint8(n & AblMask >> 18)
	p.Sub3.EffortRibbon = uint8(n & EffortMask >> 19)
	p.Sub3.BattleChampionRibbon = uint8(n & BattleChampionMask >> 20)
	p.Sub3.RegionalChampionRibbon = uint8(n & RegionalChampionMask >> 21)
	p.Sub3.NationalChampionRibbon = uint8(n & NationalChampionMask >> 22)
	p.Sub3.CountryRibbon = uint8(n & CountryMask >> 23)
	p.Sub3.NationalRibbon = uint8(n & NationalMask >> 24)
	p.Sub3.EarthRibbon = uint8(n & EarthMask >> 25)
	p.Sub3.WorldRibbon = uint8(n & WorldMask >> 26)
	p.Sub3.FatefulEnc = uint8(n & FatefulMask >> 31)
}

func (p *PStructure) CalculateRibbonValue() uint32 {
	sub := p.Sub3

	var ribbons uint32

	ribbons |= uint32(sub.CoolRibbon)
	ribbons |= uint32(sub.BeautyRibon) << 3
	ribbons |= uint32(sub.CuteRibbon) << 6
	ribbons |= uint32(sub.SmartRibbon) << 9
	ribbons |= uint32(sub.ToughRibbon) << 12

	ribbons |= uint32(sub.ChampionRibbon) << 15
	ribbons |= uint32(sub.WinningRibbon) << 16
	ribbons |= uint32(sub.VictoryRibbon) << 17
	ribbons |= uint32(sub.ArtistRibbon) << 18
	ribbons |= uint32(sub.EffortRibbon) << 19
	ribbons |= uint32(sub.BattleChampionRibbon) << 20
	ribbons |= uint32(sub.RegionalChampionRibbon) << 21
	ribbons |= uint32(sub.NationalChampionRibbon) << 22
	ribbons |= uint32(sub.CountryRibbon) << 23
	ribbons |= uint32(sub.NationalRibbon) << 24
	ribbons |= uint32(sub.EarthRibbon) << 25
	ribbons |= uint32(sub.WorldRibbon) << 26
	ribbons |= uint32(sub.FatefulEnc) << 31

	return ribbons
}

func (p *PStructure) CalculateStatsValue() uint32 {
	sub := p.Sub3
	var stats uint32

	stats |= uint32(sub.HpIV)
	stats |= uint32(sub.AtkIV) << 5
	stats |= uint32(sub.DefIV) << 10
	stats |= uint32(sub.SpeIV) << 15
	stats |= uint32(sub.SpAtkIV) << 20
	stats |= uint32(sub.SpDIV) << 25
	stats |= uint32(sub.IsEgg) << 30
	stats |= uint32(sub.AbilityNum) << 31
	return stats
}

func (p *PStructure) ReplaceStatsByValue(n uint32) {
	p.Sub3.HpIV = uint8(n & HPMask)
	p.Sub3.AtkIV = uint8(n & AtkMask >> 5)
	p.Sub3.DefIV = uint8(n & DefMask >> 10)
	p.Sub3.SpeIV = uint8(n & SpeMask >> 15)
	p.Sub3.SpAtkIV = uint8(n & SpaMask >> 20)
	p.Sub3.SpDIV = uint8(n & SpdMask >> 25)
	p.Sub3.IsEgg = uint8(n & EggMask >> 30)
	p.Sub3.AbilityNum = uint8(n & AblMask >> 31)
}

func (p *PStructure) CalculateOriginValue() uint16 {
	sub := p.Sub3
	var origin uint16

	origin |= uint16(sub.MetLevel)
	origin |= uint16(sub.MetGame) << 7
	origin |= uint16(sub.Ball) << 11
	origin |= uint16(sub.OTGender) << 14

	return origin
}

func (p *PStructure) ReplaceOriginByValue(n uint16) {
	p.Sub3.MetLevel = uint8(n & MetLevelMask)
	p.Sub3.MetGame = uint8(n & GameOriginMask)
	p.Sub3.Ball = uint8(n & BallMask)
	p.Sub3.OTGender = uint8(n & TGenderMask)
}

func (p *PStructure) MakeFirstPartRawMon() []byte {
	// All data before substructs
	pk := make([]byte, 0x20)
	binary.LittleEndian.PutUint32(pk[:0x4], p.PID)
	binary.LittleEndian.PutUint32(pk[0x4:0x8], p.OTID)

	n := utils.GenerateNickNameSlice(p.NickName)
	for i := range NickNameMaxLength {
		pk[0x8+i] = n[i]
	}

	o := utils.GenerateOTNameSlice(p.OTName)
	for i := range TrainerNameMaxLength {
		pk[0x14+i] = o[i]
	}

	pk[0x12] = p.Lang
	pk[0x13] = p.Flags
	pk[0x1B] = p.Markings

	binary.LittleEndian.PutUint16(pk[0x1C:0x1E], p.Checksum)
	binary.LittleEndian.PutUint16(pk[0x1E:0x20], p.filler)

	return pk
}

func (p *PStructure) MakeSecondPartRawMon() []byte {
	// All data after substructs
	pk := make([]byte, 0x14)

	binary.LittleEndian.PutUint32(pk[0x0:0x4], p.Status)
	pk[0x5] = p.Level
	pk[0x6] = p.MailID

	binary.LittleEndian.PutUint16(pk[0x6:0x8], p.CurHP)
	binary.LittleEndian.PutUint16(pk[0x8:0xA], p.TotalHP)
	binary.LittleEndian.PutUint16(pk[0xA:0xC], p.Atk)
	binary.LittleEndian.PutUint16(pk[0xC:0xE], p.Def)
	binary.LittleEndian.PutUint16(pk[0xE:0x10], p.Spe)
	binary.LittleEndian.PutUint16(pk[0x10:0x12], p.Spa)
	binary.LittleEndian.PutUint16(pk[0x12:0x14], p.Spd)

	return pk
}

func (p *PStructure) MakeRawGrowthSubstruct(crypt bool) []byte {
	pk := make([]byte, 0xC)

	sub := p.Sub0
	binary.LittleEndian.PutUint16(pk[:0x2], sub.Species)
	binary.LittleEndian.PutUint16(pk[0x2:0x4], sub.HeldItem)
	binary.LittleEndian.PutUint32(pk[0x4:0x8], sub.Experience)
	pk[0x8] = sub.PPBonuses
	pk[0x9] = sub.FriendShip
	binary.LittleEndian.PutUint16(pk[0xA:0xC], sub.filler)

	if crypt {
		return CryptMonSubstructs(pk, p.GetEncryptionKey())
	}

	return pk
}

func (p *PStructure) MakeRawAttacksSubstruct(crypt bool) []byte {
	pk := make([]byte, 0xC)

	sub := p.Sub1
	binary.LittleEndian.PutUint16(pk[0:0x2], sub.Moves[0])
	binary.LittleEndian.PutUint16(pk[0x2:0x4], sub.Moves[1])
	binary.LittleEndian.PutUint16(pk[0x4:0x6], sub.Moves[2])
	binary.LittleEndian.PutUint16(pk[0x6:0x8], sub.Moves[3])
	pk[0x8] = sub.PP[0]
	pk[0x9] = sub.PP[1]
	pk[0xA] = sub.PP[2]
	pk[0xB] = sub.PP[3]

	if crypt {
		return CryptMonSubstructs(pk, p.GetEncryptionKey())
	}

	return pk
}

func (p *PStructure) MakeRawEVSConditionsSubstruct(crypt bool) []byte {
	pk := make([]byte, 0xC)

	sub := p.Sub2
	pk[0x0] = sub.HpEV
	pk[0x1] = sub.AtkEV
	pk[0x2] = sub.DefEV
	pk[0x3] = sub.SpeEV
	pk[0x4] = sub.SpAtkEV
	pk[0x5] = sub.SpDefEV
	pk[0x6] = sub.Cool
	pk[0x7] = sub.Beauty
	pk[0x8] = sub.Cute
	pk[0x9] = sub.Smart
	pk[0xA] = sub.Tough
	pk[0xB] = sub.Feel

	if crypt {
		return CryptMonSubstructs(pk, p.GetEncryptionKey())
	}

	return pk
}

func (p *PStructure) MakeRawMiscSubstruct(crypt bool) []byte {
	pk := make([]byte, 0xC)

	sub := p.Sub3
	pk[0x0] = sub.Pokerus
	pk[0x1] = sub.MetLocation

	binary.LittleEndian.PutUint16(pk[0x2:0x4], p.CalculateOriginValue())
	binary.LittleEndian.PutUint32(pk[0x4:0x8], p.CalculateStatsValue())
	binary.LittleEndian.PutUint32(pk[0x8:0xC], p.CalculateRibbonValue())

	if crypt {
		return CryptMonSubstructs(pk, p.GetEncryptionKey())
	}

	return pk
}

func (p *PStructure) GetSubstructDataInOrder(crypt bool) [][]byte {
	return [][]byte{
		p.MakeRawGrowthSubstruct(crypt),
		p.MakeRawAttacksSubstruct(crypt),
		p.MakeRawEVSConditionsSubstruct(crypt),
		p.MakeRawMiscSubstruct(crypt),
	}
}

func generateChecksum(d []byte) uint16 {
	var checksum uint16

	for i := 0; i < len(d); i += 2 {
		s := utils.GetSliceFromRawData(d, i, 2)
		checksum += binary.LittleEndian.Uint16(s)
	}
	return checksum
}

func (p *PStructure) ToPK3() []byte {
	pk := make([]byte, 0x64)

	copy(pk[0x0:0x20], p.MakeFirstPartRawMon())
	copy(pk[0x50:0x64], p.MakeSecondPartRawMon())

	sl := p.GetSubstructDataInOrder(false)
	for i := range 4 {
		copy(pk[SubStructOffsetMap[i][0]:SubStructOffsetMap[i][1]], sl[i])
	}

	c := generateChecksum(pk[SubStructOffsetMap[0][0]:SubStructOffsetMap[3][1]])
	binary.LittleEndian.PutUint16(pk[0x1C:0x1E], c)

	return pk
}

func (p *PStructure) ToEK3() []byte {
	pk := make([]byte, 0x64)

	sl := p.GetSubstructDataInOrder(true)
	for i, val := range p.GetMonSubStructOrder() {
		copy(pk[SubStructOffsetMap[i][0]:SubStructOffsetMap[i][1]], sl[val])
	}

	c := generateChecksum(pk[SubStructOffsetMap[0][0]:SubStructOffsetMap[3][1]])
	binary.LittleEndian.PutUint16(pk[0x1C:0x1E], c)

	copy(pk[0x0:0x20], p.MakeFirstPartRawMon())
	copy(pk[0x50:0x64], p.MakeSecondPartRawMon())

	return pk
}

// SaveWithMail Functionally the same process as ToEK3 except we are editing PID, OTID after encryption
// to simulate the mail edits to an encrypted mon while in game
func (p *PStructure) SaveWithMail(PID, OTID uint32) []byte {
	pk := make([]byte, 0x64)

	if p.GetEncryptionKey() == 0 {
		p.PID = PID
		p.OTID = OTID
	}

	sl := p.GetSubstructDataInOrder(true)
	for i, val := range p.GetMonSubStructOrder() {
		copy(pk[SubStructOffsetMap[i][0]:SubStructOffsetMap[i][1]], sl[val])
	}

	c := generateChecksum(pk[SubStructOffsetMap[0][0]:SubStructOffsetMap[3][1]])
	binary.LittleEndian.PutUint16(pk[0x1C:0x1E], c)

	copy(pk[0x0:0x20], p.MakeFirstPartRawMon())
	copy(pk[0x50:0x64], p.MakeSecondPartRawMon())

	binary.LittleEndian.PutUint32(pk[0:4], PID)
	binary.LittleEndian.PutUint32(pk[4:8], OTID)

	return pk
}

func GetErraticTable() []uint {
	return []uint{
		0, 15, 52, 122, 237, 406, 637, 942, 326, 1800,
		2369, 3041, 3822, 4719, 5737, 6881, 8155, 9564,
		11111, 12800, 14632, 16610, 18737, 21012, 23437,
		26012, 28737, 31610, 34632, 37800, 41111, 44564,
		48155, 51881, 55737, 59719, 63822, 68041, 72369,
		76800, 81326, 85942, 90637, 95406, 100237, 105122,
		110052, 115015, 120001, 125000, 131324, 137795,
		144410, 151165, 158056, 165079, 172229, 179503,
		186894, 194400, 202013, 209728, 217540, 225443, 233431,
		241496, 249633, 257834, 267406, 276458, 286328, 296358,
		305767, 316074, 326531, 336255, 346965, 357812, 367807,
		378880, 390077, 400293, 411686, 423190, 433572, 445239,
		457001, 467489, 479378, 491346, 501878, 513934, 526049,
		536557, 548720, 560922, 571333, 583539, 591882, 600000,
	}
}

func GetFastTable() []uint {
	return []uint{
		0, 6, 21, 51, 100, 172, 274, 409, 583, 800,
		1064, 1382, 1757, 2195, 2700, 3276, 3930, 4665,
		5487, 6400, 7408, 8518, 9733, 11059, 12500, 14060,
		15746, 17561, 19511, 21600, 23832, 26214, 28749,
		31443, 34300, 37324, 40522, 43897, 47455, 51200,
		55136, 59270, 63605, 68147, 72900, 77868, 83058,
		88473, 94119, 100000, 106120, 112486, 119101,
		125971, 133100, 140492, 148154, 156089, 164303,
		172800, 181584, 190662, 200037, 209715, 219700,
		229996, 240610, 251545, 262807, 274400, 286328,
		298598, 311213, 324179, 337500, 351180, 365226,
		379641, 394431, 409600, 425152, 441094, 457429,
		474163, 491300, 508844, 526802, 545177, 563975,
		583200, 602856, 622950, 643485, 664467, 685900,
		707788, 730138, 752953, 776239, 800000}
}

func GetMediumFastTable() []uint {
	return []uint{
		8, 27, 64, 125, 216, 343, 512, 729, 1000, 1331,
		1728, 2197, 2744, 3375, 4096, 4913, 5832, 6859,
		8000, 9261, 10648, 12167, 13824, 15625, 17576, 19683,
		21952, 24389, 27000, 29791, 32768, 35937, 39304, 42875,
		46656, 50653, 54872, 59319, 64000, 68921, 74088, 79507,
		85184, 91125, 97336, 103823, 110592, 117649, 125000,
		132651, 140608, 148877, 157464, 166375, 175616, 185193,
		195112, 205379, 216000, 226981, 238328, 250047, 262144,
		274625, 287496, 300763, 314432, 328509, 343000, 357911,
		373248, 389017, 405224, 421875, 438976, 456533, 474552,
		493039, 512000, 531441, 551368, 571787, 592704, 614125,
		636056, 658503, 681472, 704969, 729000, 753571, 778688,
		804357, 830584, 857375, 884736, 912673, 941192, 970299,
		1000000,
	}
}

func GetMediumSlowTable() []uint {
	return []uint{
		0, 9, 57, 96, 135, 179, 236, 314, 419, 560,
		742, 973, 1261, 1612, 2035, 2535, 3120, 3798,
		4575, 5460, 6458, 7577, 8825, 10208, 11735,
		13411, 15244, 17242, 19411, 21760, 24294, 27021,
		29949, 33084, 36435, 40007, 43808, 47846, 52127,
		56660, 61450, 66505, 71833, 77440, 83335, 89523,
		96012, 102810, 109923, 117360, 125126, 133229,
		141677, 150476, 159635, 169159, 179056, 189334,
		199999, 211060, 222522, 234393, 246681, 259392,
		272535, 286115, 300140, 314618, 329555, 344960,
		360838, 377197, 394045, 411388, 429235, 447591,
		466464, 485862, 505791, 526260, 547274, 568841,
		590969, 613664, 636935, 660787, 685228, 710266,
		735907, 762160, 789030, 816525, 844653, 873420,
		902835, 932903, 963632, 995030, 1027103, 1059860,
	}
}

func GetSlowTable() []uint {
	return []uint{
		0, 10, 33, 80, 156, 270, 428, 640, 911,
		1250, 1663, 2160, 2746, 3430, 4218, 5120, 6141,
		7290, 8573, 10000, 11576, 13310, 15208, 17280,
		19531, 21970, 24603, 27440, 30486, 33750, 37238,
		40960, 44921, 49130, 53593, 58320, 63316, 68590,
		74148, 80000, 86151, 92610, 99383, 106480, 113906,
		121670, 129778, 138240, 147061, 156250, 165813,
		175760, 186096, 196830, 207968, 219520, 231491,
		243890, 256723, 270000, 283726, 297910, 312558,
		327680, 343281, 359370, 375953, 393040, 410636,
		428750, 447388, 466560, 486271, 506530, 527343,
		548720, 570666, 593190, 616298, 640000, 664301,
		689210, 714733, 740880, 767656, 795070, 823128,
		851840, 881211, 911250, 941963, 973360, 1005446,
		1038230, 1071718, 1105920, 1140841, 1176490, 1212873,
		1250000,
	}
}

func GetFluctuatingTable() []uint {
	return []uint{
		0, 4, 13, 32, 65, 112, 178, 276, 393, 540, 745, 967,
		1230, 1591, 1957, 2457, 3046, 3732, 4526, 5440, 6482,
		7666, 9003, 10506, 12187, 14060, 16140, 18439, 20974,
		23760, 26811, 30146, 33780, 37731, 42017, 46656, 50653,
		55969, 60505, 66560, 71677, 78533, 84277, 91998, 98415,
		107069, 114205, 123863, 131766, 142500, 151222, 163105,
		172697, 185807, 196322, 210739, 222231, 238036, 250562,
		267840, 281456, 300293, 315059, 335544, 351520, 373744,
		390991, 415050, 433631, 459620, 479600, 507617, 529063,
		559209, 582187, 614566, 639146, 673863, 700115, 737280,
		765275, 804997, 834809, 877201, 908905, 954084, 987754,
		1035837, 1071552, 1122660, 1160499, 1214753, 1254796,
		1312322, 1354652, 1415577, 1460276, 1524731, 1571884, 1640000,
	}
}
