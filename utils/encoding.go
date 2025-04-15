package utils

// https://bulbapedia.bulbagarden.net/wiki/Character_encoding_(Generation_III)

var westernCharsEnc = []rune{
	' ', 'À', 'Á', 'Â', 'Ç', 'È', 'É', 'Ê', 'Ë', 'Ì', ' ', 'Î', 'Ï', 'Ò', 'Ó', 'Ô',
	'Œ', 'Ù', 'Ú', 'Û', 'Ñ', 'ß', 'à', 'á', ' ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í',
	'î', 'ï', 'ò', 'ó', 'ô', 'œ', 'ù', 'ú', 'û', 'ñ', 'º', 'ª', ' ', '&', '+', ' ',
	' ', ' ', ' ', ' ', 'ˡ', '=', ';', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
	' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
	'▯', '¿', '¡', ' ', ' ', ' ', ' ', ' ', ' ', ' ', 'Í', '%', '(', ')', ' ', ' ',
	' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', 'â', ' ', ' ', ' ', ' ', ' ', ' ', 'í',
	' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '⬆', '⬇', '⬅', '⮕', '*', '*', '*',
	'*', '*', '*', '*', 'ᵉ', '<', '>', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
	' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
	' ', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!', '?', '.', '-', '・',
	'…', '“', '”', '‘', '’', '♂', '♀', '$', ',', '×', '/', 'A', 'B', 'C', 'D', 'E',
	'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
	'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
	'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '▶',
	':', 'Ä', 'Ö', 'Ü', 'ä', 'ö', 'ü', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
}

var jpnCharEnc = []rune{
	'　', 'あ', 'い', 'う', 'え', 'お', 'か', 'き', 'く', 'け', 'こ', 'さ', 'し', 'す', 'せ', 'そ',
	'た', 'ち', 'つ', 'て', 'と', 'な', 'に', 'ぬ', 'ね', 'の', 'は', 'ひ', 'ふ', 'へ', 'ほ', 'ま',
	'み', 'む', 'め', 'も', 'や', 'ゆ', 'よ', 'ら', 'り', 'る', 'れ', 'ろ', 'わ', 'を', 'ん', 'ぁ',
	'ぃ', 'ぅ', 'ぇ', 'ぉ', 'ゃ', 'ゅ', 'ょ', 'が', 'ぎ', 'ぐ', 'げ', 'ご', 'ざ', 'じ', 'ず', 'ぜ',
	'ぞ', 'だ', 'ぢ', 'づ', 'で', 'ど', 'ば', 'び', 'ぶ', 'べ', 'ぼ', 'ぱ', 'ぴ', 'ぷ', 'ぺ', 'ぽ',
	'っ', 'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ', 'サ', 'シ', 'ス', 'セ', 'ソ',
	'タ', 'チ', 'ツ', 'テ', 'ト', 'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ', 'ヒ', 'フ', 'ヘ', 'ホ', 'マ',
	'ミ', 'ム', 'メ', 'モ', 'ヤ', 'ユ', 'ヨ', 'ラ', 'リ', 'ル', 'レ', 'ロ', 'ワ', 'ヲ', 'ン', 'ァ',
	'ィ', 'ゥ', 'ェ', 'ォ', 'ャ', 'ュ', 'ョ', 'ガ', 'ギ', 'グ', 'ゲ', 'ゴ', 'ザ', 'ジ', 'ズ', 'ゼ',
	'ゾ', 'ダ', 'ヂ', 'ヅ', 'デ', 'ド', 'バ', 'ビ', 'ブ', 'ベ', 'ボ', 'パ', 'ピ', 'プ', 'ペ', 'ポ',
	'ッ', '０', '１', '２', '３', '４', '５', '６', '７', '８', '９', '！', '？', '。', '－', '・',
	'‥', '『', '』', '「', '」', '♂', '♀', '＄', '．', '×', '／', 'Ａ', 'Ｂ', 'Ｃ', 'Ｄ', 'Ｅ',
	'Ｆ', 'Ｇ', 'Ｈ', 'Ｉ', 'Ｊ', 'Ｋ', 'Ｌ', 'Ｍ', 'Ｎ', 'Ｏ', 'Ｐ', 'Ｑ', 'Ｒ', 'Ｓ', 'Ｔ', 'Ｕ',
	'Ｖ', 'Ｗ', 'Ｘ', 'Ｙ', 'Ｚ', 'ａ', 'ｂ', 'ｃ', 'ｄ', 'ｅ', 'ｆ', 'ｇ', 'ｈ', 'ｉ', 'ｊ', 'ｋ',
	'ｌ', 'ｍ', 'ｎ', 'ｏ', 'ｐ', 'ｑ', 'ｒ', 'ｓ', 'ｔ', 'ｕ', 'ｖ', 'ｗ', 'ｘ', 'ｙ', 'ｚ', '▶',
	'：', 'Ä', 'Ö', 'Ü', 'ä', 'ö', 'ü', '⬆', '⬇', '⬅', ' ', ' ', ' ', ' ', ' ', ' ',
}

func findRuneIndex(r rune) int {
	for i, rune := range westernCharsEnc {
		if r == rune {
			return i
		}
	}
	return 0
}

// This is pretty jank but it works for now
func GenerateNickNameString(bs []byte) string {
	if len(bs) > 10 {
		return ""
	}

	out := make([]rune, 10)

	for i, v := range bs {
		if v == 0xFF {
			out[i] = '\\'
			break
		} else {
			out[i] = westernCharsEnc[v]
		}
	}

	return string(out)
}

func GenerateOTNameString(bs []byte) string {
	if len(bs) > 7 {
		return ""
	}

	out := make([]rune, 7)

	for i, v := range bs {
		if v == 0xFF {
			out[i] = '\\'
			break
		} else {
			out[i] = westernCharsEnc[v]
		}

	}

	return string(out)
}

func GenerateNickNameSlice(s string) []byte {
	r := []rune(s)

	if len(r) > 10 {
		return nil
	}

	out := make([]byte, 10)

	for i, v := range r {
		if i == 10 {
			break
		}
		// TODO: Add const vals for special rune here and length consts
		if v == '\\' {
			out[i] = 0xFF
		} else {
			c := findRuneIndex(v)
			out[i] = byte(c)
		}
	}
	return out
}

func GenerateOTNameSlice(s string) []byte {
	r := []rune(s)

	if len(r) > 7 {
		return nil
	}

	out := make([]byte, 7)

	for i, v := range r {
		if v == '\\' {
			out[i] = 0xFF
		} else {
			c := findRuneIndex(v)
			out[i] = byte(c)
		}
	}
	return out
}

func WesternSliceToString(bs []byte) string {
	out := []rune{}
	for _, s := range bs {
		out = append(
			out,
			westernCharsEnc[s],
		)
	}

	return string(out)
}

func WesternStringToByteSlice(s string) []byte {
	out := []byte{}

	for _, c := range s {
		i := findRuneIndex(c)
		out = append(out, byte(i))
	}

	return out
}

func JPNSliceToString(bs []byte) string {
	out := []rune{}
	for _, s := range bs {
		out = append(
			out,
			jpnCharEnc[s],
		)
	}

	return string(out)
}

func JPNStringToByteSlice(s string) []byte {
	out := []byte{}

	for _, c := range s {
		out = append(
			out,
			byte(jpnCharEnc[c]),
		)
	}

	return out
}
