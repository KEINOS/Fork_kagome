package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/gojp/kana"
	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/width"
)

func main() {
	// User input text
	input := `
ローマ字変換プログラム作ってみた。
五街道のひとつである、東海道五十三次の品川宿などを変換してみると面白いかもしれない。
こぼれたままの流星群  一秒  一秒
流転 lights 消せないコナゴナ銀河。
駆逐艦"そよかぜ"
光る 雲を突き抜け Fly Away! (Fly Away) フライ アウェイ
東海道五十三次の品川宿を訪問する。
`

	// Built-in user dictionary
	usrDict := `
東海道五十三次,東海道 五十三 次,トウカイドウ ゴジュウサン ツギ,カスタム名詞
品川宿,品川 宿,シナガワ ジュク,カスタム名詞
`

	// Convert user dictionary string to tokenizer.Option
	usrDictOpt, err := newUserDictOpt(usrDict)
	if err != nil {
		log.Fatal(err)
	}

	// Create IPA-dict-based tokenizer with user dictionary
	tkn, err := tokenizer.New(ipa.Dict(), usrDictOpt, tokenizer.OmitBosEos())
	if err != nil {
		log.Fatal(err)
	}

	// Split input text by line
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		// Get Yomi (pronunciation/reading) from the line in Katakana
		yomi := getYomi(tkn, line)
		if yomi == "" {
			continue // ignore empty lines
		}

		// Transliterate to Romaji.
		romaji := getRomantic(yomi)

		fmt.Println(romaji)
	}
	//
	// Output:
	// ro-maji henkan puroguramu tsukutte mita.
	// go kaido- no hitotsudearu, toukaidou gojuusan tsugi no shinagawa juku nado wo henkan shite miruto omoshiroi kamo shirenai.
}

var conversionMap = map[rune]rune{
	'、': ',',
	'。': '.',
	'！': '!',
	'？': '?',
	'「': '"',
	'」': '"',
	'『': '"',
	'』': '"',
}

// asHalfWidth converts full-width alpha-numeric characters and symbols to
// half-width characters according to the conversionMap.
// Note that half-width katakana characters are converted to full-width.
func asHalfWidth(input string) string {
	// Convert half-width katakana characters to full-width and full-width
	// alphanumeric characters to half-width.
	input = width.Fold.String(input)

	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Han, r) {
			return r
		}

		if converted, ok := conversionMap[r]; ok {
			return converted
		}

		return r
	}, input)
}

// getRomantic returns the Romaji transliteration of the input in Katakana.
func getRomantic(line string) (yomi string) {
	defer func() {
		// Finally, remove extra spaces
		if yomi != "" {
			yomi = strings.Join(strings.Fields(yomi), " ")
		}
	}()

	// In this example we use the github.com/gojp/kana package for Katakana
	// transliteration to Romaji. However, other packages are available.
	// Such as:
	// - github.com/robpike/nihongo
	// - github.com/kotaroooo0/gojaconv
	// - github.com/yosida95/romaji
	// - github.com/goark/krconv
	romaji := kana.KanaToRomaji(line)

	// Barely normalize full-width symbols to half-width
	// ('、' -> ',', '。' -> '.', etc.)
	romaji = asHalfWidth(romaji)

	// Capitalize the first letter of each sentence
	sentences := strings.SplitAfter(romaji, ".")

	for index, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		isFirst := true

		// Capitalize the first letter of each sentence
		sentence = strings.Map(func(r rune) rune {
			if isFirst {
				isFirst = false

				return unicode.ToUpper(r)
			}

			return r
		}, sentence)

		//sentences[index] = cases.Title(language.English).String(sentence)
		sentences[index] = sentence
	}

	return strings.Join(sentences, " ")
}

// getYomi returns the pronunciation/reading (Yomi) of the input in Katakana
// using the given tokenizer.
func getYomi(tkn *tokenizer.Tokenizer, line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	if isASCII(line) {
		return line
	}

	tokens := tkn.Tokenize(line)
	chunks := []string{}
	tmpChunk := ""
	isPrevASCII := false

	// Evaluate each token to retrieve the pronunciation or reading as a
	// slice of chunks to join them later. It is similar to Wakachi, but
	// with a bit more complex logic.
	for _, token := range tokens {
		prevKey := len(chunks) - 1

		// Detect ASCII words
		if isASCII(token.Surface) {
			if isPrevASCII {
				chunks[prevKey] += token.Surface
			} else {
				chunks = append(chunks, token.Surface)
				isPrevASCII = true
			}

			continue
		} else if isPrevASCII {
			// Capitalize the previous chunk if it was all in ASCII
			chunks[prevKey] = strings.ToUpper(chunks[prevKey])
		}

		isPrevASCII = false

		// Retrieve the pronunciation/reading from the token in katakana
		if usrExtra := token.UserExtra(); usrExtra != nil {
			tmpChunk = strings.Join(usrExtra.Readings, " ")
		} else if p, ok := token.Pronunciation(); ok {
			tmpChunk = p
		} else if r, ok := token.Reading(); ok {
			tmpChunk = r // fallback to reading if pronunciation is not available
		} else {
			tmpChunk = token.Surface
		}

		tmpChunk = strings.TrimSpace(tmpChunk)
		//fmt.Println("Log:", tmpChunk, token.POS())

		if isPartOfPrev(token) {
			chunks[prevKey] += tmpChunk // Append to the previous chunk
		} else {
			chunks = append(chunks, tmpChunk) // Append to the slice of chunks
		}
	}

	return strings.Join(chunks, " ")
}

// isASCII returns true if the string is all in ASCII.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// isPartOfPrev returns true if the token prefers to be part of the previous chunk.
//
// e.g. tsuku te mi ta。--> tsukutte mita。
func isPartOfPrev(token tokenizer.Token) bool {
	// Not "助詞" "助動詞" nor "記号"
	if !strings.ContainsAny(token.POS()[0], "助"+"記") {
		return false
	}

	switch token.POS()[1] {
	// Ignore below particles, conjunctions, and auxiliary verbs
	case "副助詞", "連体化", "格助詞":
		return false
	// Else, consider as part of the previous chunk
	default:
		return true
	}
}

// newUserDictOpt creates a tokenizer.Option from a user dictionary string.
func newUserDictOpt(rec string) (tokenizer.Option, error) {
	// Read user dictionary records from the string.
	usrDictRec, err := dict.NewUserDicRecords(strings.NewReader(rec))
	if err != nil {
		return nil, err
	}

	// Create a dict.UserDict from the records.
	usrDict, err := usrDictRec.NewUserDict()
	if err != nil {
		return nil, err
	}

	// Cast the UserDict to tokenizer.Option.
	return tokenizer.UserDict(usrDict), nil
}
