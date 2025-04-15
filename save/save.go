package save

import (
	"bytes"
	"encoding/binary"
	"os"
	"postal/utils"
)

const (
	SaveSize        uint = 0x20000
	SaveSectionSize uint = 0x1000
	PkSize          uint = 0x64
)

func getSignatureSlice() []byte {
	return []byte{0x25, 0x20, 0x01, 0x08}
}

func getBlockOffsets() []uint {
	return []uint{0x0, 0x00E000, 0x01C000, 0x01E000, 0x01F000}
}

func getSectionFooterOffsets() []uint {
	return []uint{0x0, 0x0FF4, 0x0FF6, 0x0FF8, 0x0FFC}
}

func GetSectionSizes() []uint {
	return []uint{3884, 3968, 3968, 3968, 3848, 3968, 3968, 3968, 3968, 3968, 3968, 3968, 3968, 2000}
}

func getFRLGTrainerOffsets() []uint {
	return []uint{0x0, 0x8, 0x9, 0xA, 0xE, 0x13, 0xAC, 0xAF8}
}

func getRSTrainerOffsets() []uint {
	return []uint{0x0, 0x8, 0x9, 0xA, 0xE, 0x13, 0xAC}
}

func getETrainerOffsets() []uint {
	return []uint{0x0, 0x8, 0x9, 0xA, 0xE, 0x13, 0xAC}
}

func getFRLGTeamsAndItemsOffsets() []uint {
	return []uint{0x0034, 0x0038, 0x0290, 0x0294, 0x0298, 0x0310, 0x03B8, 0x0430, 0x0464, 0x054C}
}

func getRSTeamAndItemsOffsets() []uint {
	return []uint{0x0234, 0x0238, 0x0490, 0x0494, 0x0498, 0x0560, 0x05B0, 0x0600, 0x0640, 0x740}
}

func getETeamAndItemsOffsets() []uint {
	return []uint{0x0234, 0x0238, 0x0490, 0x0494, 0x0498, 0x0560, 0x05D8, 0x0650, 0x0690, 0x790}
}

type (
	RawMonData  []byte
	RawSaveFile struct {
		Data          []byte
		isCorrectSize bool
	}

	SaveBlock struct {
		offset   uint
		sections [14]SaveBlockSection
	}

	SaveBlockSection struct {
		offset    uint
		data      []byte
		sectionID []byte
		checksum  []byte
		signature []byte
		saveIndex []byte
		isValid   bool
	}

	RawTrainerInfo struct {
		Offset  uint
		Name    []byte
		Gender  []byte
		TID     []byte
		SID     []byte
		Time    []byte
		Options []byte
		Code    []byte
		Key     []byte
	}

	RawTeamAndItems struct {
		Offset   uint
		TeamSize []byte
		MonList  []byte
		Money    []byte
		Coins    []byte
		PCItems  []byte
		Items    []byte
		KeyItems []byte
		Balls    []byte
		TMCase   []byte
		Berries  []byte
	}

	RawBoxDataTotal struct {
		CurrentBox uint
		Boxes      [9]RawBoxDataBuffer
		BoxNames   []byte
		WallPapers []byte
	}

	RawBoxDataBuffer struct {
		Offset uint
		Data   []byte
	}
)

func GenerateRawSaveData(p string) RawSaveFile {
	save := RawSaveFile{
		nil,
		false,
	}

	if checkFileSize(p) {
		f, err := os.Open(p)
		checkError(err)

		defer f.Close()

		buffer := make([]byte, SaveSize)

		_, err = f.Read(buffer)
		checkError(err)

		save.Data = buffer
		save.isCorrectSize = true
	}

	return save
}

func GetRawMonDataFromFile(p string) (RawMonData, error) {

	if checkPkSize(p) {

		f, err := os.Open(p)
		checkError(err)

		buffer := make([]byte, PkSize)

		_, err = f.Read(buffer)
		checkError(err)

		return buffer, nil
	}

	return RawMonData{}, utils.ErrPKIncorrectFileSize
}

func (r RawMonData) GetSubStructureArray() [][]byte {
	b := make([][]byte, 4)
	b[0] = utils.GetSliceFromRawData(r, 0x20, 0xC)
	b[1] = utils.GetSliceFromRawData(r, 0x2C, 0xC)
	b[2] = utils.GetSliceFromRawData(r, 0x38, 0xC)
	b[3] = utils.GetSliceFromRawData(r, 0x44, 0xC)
	return b
}

func (s *RawSaveFile) getActiveSaveBlockOffsetAndSaveCount() (uint, uint) {
	// TODO: Add save file size check error
	blockASaveIndex := getBlockOffsets()[0] + getSectionFooterOffsets()[4]
	blockBSaveIndex := getBlockOffsets()[1] + getSectionFooterOffsets()[4]

	blockASaveVal := binary.LittleEndian.Uint16(
		utils.GetSliceFromRawData(
			s.Data,
			int(blockASaveIndex),
			4,
		),
	)
	blockBSaveVal := binary.LittleEndian.Uint16(
		utils.GetSliceFromRawData(
			s.Data,
			int(blockBSaveIndex),
			4,
		),
	)

	if blockBSaveVal >= blockASaveVal {
		return getBlockOffsets()[1], uint(blockBSaveVal)
	}

	return getBlockOffsets()[0], uint(blockASaveVal)
}

func (s *RawSaveFile) GenerateSaveBlock() SaveBlock {
	saveBlock := SaveBlock{}
	block := SaveBlockSection{}

	offset, _ := s.getActiveSaveBlockOffsetAndSaveCount()

	saveBlock.offset = offset
	for range GetSectionSizes() {
		block.offset = offset
		block.data = utils.GetSliceFromRawData(
			s.Data,
			int(offset),
			int(SaveSectionSize),
		)
		block.sectionID = utils.GetSliceFromRawData(
			s.Data,
			int(offset)+int(getSectionFooterOffsets()[1]),
			2,
		)
		block.checksum = utils.GetSliceFromRawData(
			s.Data,
			int(offset)+int(getSectionFooterOffsets()[2]),
			2,
		)
		block.signature = utils.GetSliceFromRawData(
			s.Data,
			int(offset)+int(getSectionFooterOffsets()[3]),
			4,
		)
		block.saveIndex = utils.GetSliceFromRawData(
			s.Data,
			int(offset)+int(getSectionFooterOffsets()[4]),
			4,
		)

		sigValid := bytes.Equal(getSignatureSlice(), block.signature)

		checkSumValid := generateChecksum(
			block.data,
			int(GetSectionSizes()[binary.LittleEndian.Uint16(block.sectionID)]),
		) == binary.LittleEndian.Uint16(block.checksum)

		block.isValid = (sigValid && checkSumValid)
		sectionIDVal := binary.LittleEndian.Uint16(block.sectionID)
		saveBlock.sections[sectionIDVal] = block
		offset += SaveSectionSize
	}
	return saveBlock
}

func (b *SaveBlock) GetDirectPlayerXKey() uint16 {
	sec := b.sections[0].data
	return binary.LittleEndian.Uint16(
		utils.GetSliceFromRawData(
			sec,
			int(getFRLGTrainerOffsets()[7]),
			4,
		),
	)
}

func (b *SaveBlock) GetRawTrainerInfo() RawTrainerInfo {
	tb := b.sections[0].data
	return RawTrainerInfo{
		Offset: b.sections[0].offset,
		Name: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[0]),
			7,
		),

		Gender: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[1]),
			1,
		),

		TID: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[3]),
			2,
		),

		SID: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[3]+2),
			2,
		),

		Time: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[4]),
			5,
		),
		Options: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[5]),
			3,
		),
		Code: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[6]),
			4,
		),
		Key: utils.GetSliceFromRawData(
			tb,
			int(getFRLGTrainerOffsets()[7]),
			4,
		),
	}
}

func (b *SaveBlock) GetRawTeamAndItems() RawTeamAndItems {
	ti := b.sections[1]

	return RawTeamAndItems{
		Offset: ti.offset,

		TeamSize: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[0]),
			4,
		),

		MonList: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[1]),
			600,
		),

		Money: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[2]),
			4,
		),

		Coins: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[3]),
			2,
		),

		PCItems: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[4]),
			120,
		),

		Items: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[5]),
			168,
		),

		KeyItems: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[6]),
			120,
		),

		Balls: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[7]),
			52,
		),

		TMCase: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[8]),
			232,
		),

		Berries: utils.GetSliceFromRawData(
			ti.data,
			int(getFRLGTeamsAndItemsOffsets()[7]),
			52,
		),
	}
}

func (b *SaveBlock) GetRawBoxData() RawBoxDataTotal {
	var bdt RawBoxDataTotal
	boxData := b.sections[5:]

	bdt.CurrentBox = uint(binary.LittleEndian.Uint16(
		utils.GetSliceFromRawData(
			boxData[0].data,
			0x0,
			0x4,
		),
	))

	var buf RawBoxDataBuffer
	for i := range boxData {
		buf.Offset = boxData[i].offset

		if i == 0 {
			buf.Data = utils.GetSliceFromRawData(boxData[i].data, 0x4, 0xF7C)
		} else if i == 8 {
			buf.Data = utils.GetSliceFromRawData(boxData[i].data, 0x4, 0x7D0)
		} else {
			buf.Data = utils.GetSliceFromRawData(boxData[i].data, 0x0, 0xF80)
		}
		bdt.Boxes[i] = buf
	}

	bdt.BoxNames = utils.GetSliceFromRawData(boxData[8].data, 0x744, 0x7E)
	bdt.WallPapers = utils.GetSliceFromRawData(boxData[8].data, 0x87D, 0xE)
	return bdt
}

func generateChecksum(data []byte, size int) uint16 {
	var checksum uint32 = 0

	for i := 0; i < size; i += 4 {
		s := utils.GetSliceFromRawData(data, i, 4)
		checksum += binary.LittleEndian.Uint32(s)
	}

	return uint16((checksum >> 16) + (checksum & 0xFFFF))
}

func checkFileSize(p string) bool {
	f, err := os.Stat(p)
	checkError(err)
	return f.Size() == int64(SaveSize)
}

func checkPkSize(p string) bool {
	f, err := os.Stat(p)
	checkError(err)
	return f.Size() == int64(PkSize)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
