package codec

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type codecSupport struct {
	encode func(s string, method int) []byte
	decode func(b []byte) string
}

var codecTable = [...]*codecSupport{
	&codecSupport{toGBK, fromGBK},
	&codecSupport{toShiftJIS, fromShiftJIS},
	&codecSupport{toUtf16, fromUtf16},
	&codecSupport{toUtf16NoBom, fromUtf16},
	&codecSupport{toUtf16BENoBom, fromUtf16},
	&codecSupport{toUtf8, fromUtf8},
	&codecSupport{toUtf8Sig, fromUtf8},
}

//Codecs
const (
	Unknown  = -1
	GBK      = 0
	C936     = 0
	ShiftJIS = 1
	C932     = 1
	UTF16    = 2
	UTF16LE  = 3
	UTF16BE  = 4
	UTF8     = 5
	UTF8Sig  = 6
)

//Replace method
const (
	Replace     = 1
	ReplaceHTML = 2
)

var codecMap = map[string]int{
	"gbk":       0,
	"936":       0,
	"shiftjis":  1,
	"932":       1,
	"u16":       2,
	"utf16":     2,
	"utf-16-le": 3,
	"utf-16-be": 4,
}

//Encode encode string to []byte
func Encode(s string, codec int, method int) []byte {
	return codecTable[codec].encode(s, method)
}

//Decode decode []byte to string
func Decode(b []byte, codec int) string {
	return codecTable[codec].decode(b)
}

func utf8ToCodec(enc *encoding.Encoder, s string, method int) (b []byte) {
	var reader io.Reader
	switch method {
	case Replace:
		reader = transform.NewReader(bytes.NewReader([]byte(s)),
			encoding.ReplaceUnsupported(enc))
	case ReplaceHTML:
		reader = transform.NewReader(bytes.NewReader([]byte(s)),
			encoding.HTMLEscapeUnsupported(enc))
	case 0:
		reader = transform.NewReader(bytes.NewReader([]byte(s)), enc)
	default:
		panic(errors.New("unknown replace method"))
	}

	d, e := ioutil.ReadAll(reader)
	if e != nil {
		panic(e)
	}
	return d
}

func codecToUtf8(dec *encoding.Decoder, b []byte) (s string) {
	reader := transform.NewReader(bytes.NewReader(b), dec)
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		panic(e)
	}
	return string(d)
}

func toGBK(s string, method int) []byte {
	return utf8ToCodec(simplifiedchinese.GBK.NewEncoder(), s, method)
}

func fromGBK(b []byte) string {
	return codecToUtf8(simplifiedchinese.GBK.NewDecoder(), b)
}

func toShiftJIS(s string, method int) []byte {
	return utf8ToCodec(japanese.ShiftJIS.NewEncoder(), s, method)
}

func fromShiftJIS(b []byte) string {
	return codecToUtf8(japanese.ShiftJIS.NewDecoder(), b)
}

func toUtf16(s string, method int) []byte {
	return utf8ToCodec(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder(), s, method)
}

func toUtf16NoBom(s string, method int) []byte {
	return utf8ToCodec(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), s, method)
}

func fromUtf16(b []byte) string {
	return codecToUtf8(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder(), b)
}

func toUtf16BE(s string, method int) []byte {
	return utf8ToCodec(unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder(), s, method)
}

func toUtf16BENoBom(s string, method int) []byte {
	return utf8ToCodec(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder(), s, method)
}

func fromUtf16BE(b []byte) string {
	return codecToUtf8(unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder(), b)
}

func toUtf8(s string, method int) []byte {
	return []byte(s)
}

func toUtf8Sig(s string, method int) []byte {
	return append([]byte{0xef, 0xbb, 0xbf}, []byte(s)...)
}

func fromUtf8(b []byte) string {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return string(b[3:])
	}
	return string(b)
}
