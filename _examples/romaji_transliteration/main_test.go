package main

import (
	"log"
	"testing"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func Test_asHalfWidth(t *testing.T) {
	for index, test := range []struct {
		input  string
		expect string
	}{
		{input: "、。！？", expect: ",.!?"},
		{input: "ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ", expect: "abcdefghijklmnopqrstuvwxyz"},
		{input: "ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ", expect: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{input: "１２３４５６７８９０", expect: "1234567890"},
		{input: "漢字", expect: "漢字"},
		{input: "ｱｲｳｴｵ", expect: "アイウエオ"},
	} {
		result := asHalfWidth(test.input)

		if result != test.expect {
			t.Errorf("test #%d failed. Expected: %v, got: %v", index+1, test.expect, result)
		}
	}
}

func Test_getYomi(t *testing.T) {
	tkn, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		log.Fatal(err)
	}

	for index, test := range []struct {
		input  string
		expect string
	}{
		{input: "ローマ字変換プログラム作ってみた。", expect: "ローマジ ヘンカン プログラム ツクッテ ミタ。"},
		{input: "すもももももももものうち", expect: "スモモモ モモモ モモ ノ ウチ"},
		{input: "This is an English example.", expect: "This is an English example."},
		{input: "例えばstringsと言う用語がある。", expect: "タトエバ STRINGS ト イウ ヨーゴ ガ アル。"},
	} {
		result := getYomi(tkn, test.input)

		if result != test.expect {
			t.Errorf("test #%d failed. Expected: %v, got: %v", index+1, test.expect, result)
		}
	}
}

func Test_isASCII(t *testing.T) {
	for index, test := range []struct {
		input  string
		expect bool
	}{
		{input: " ", expect: true},
		{input: "A", expect: true},
		{input: "　", expect: false},
		{input: "Ａ", expect: false},
	} {
		result := isASCII(test.input)

		if result != test.expect {
			t.Errorf("test #%d failed.\n  Input: %#v Expected: %v, got: %v",
				index+1, test.input, test.expect, result)
		}
	}
}
