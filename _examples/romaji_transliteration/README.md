# Romaji Transliteration (ローマ字変換)

This example demonstrates:

1. Convert a string of user dictionary to `dict.UserDict` type as a `tokenizer.Option`.
1. Retrieve the readings as much as possible from the token information:
    - `tokenizer.Token.Pronunciation`, `tokenizer.Token.Reading` and `tokenizer.UserExtra.Readings`
1. Adjust the readings by referring the POS (Part of Speech) information.
1. Transliterate Japanese readings to Romaji using the `github.com/gojp/kana` package.

```shellsession
$ cd /path/to/kagome/_examples/romaji_transliteration
$ # Build the binary
$ go build -o yomi .
$ # Run the binary
$ ./yomi
ro-maji henkan puroguramu tsukutte mita.
go kaido- no hitotsudearu, toukaidou gojuusan tsugi no shinagawa juku nado wo henkan shite miruto omoshiroi kamo shirenai.
```

- This example is based on the [issue #308](https://github.com/ikawaha/kagome/issues/308).
