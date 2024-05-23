module kagome/examples/romaji_transliteration

go 1.19

require (
	github.com/gojp/kana v0.1.0
	github.com/ikawaha/kagome-dict v1.0.10
	github.com/ikawaha/kagome-dict/ipa v1.0.11
	github.com/ikawaha/kagome/v2 v2.0.0-00010101000000-000000000000
	golang.org/x/text v0.15.0
)

replace github.com/ikawaha/kagome/v2 => ../../
