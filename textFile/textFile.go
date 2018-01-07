package textFile

import (
	"io/ioutil"
	"strings"

	"github.com/regomne/eutil/codec"
)

// ReadWin32TxtToLines read win32 txt file to lines (chinese os)
func ReadWin32TxtToLines(fname string) (ret []string, err error) {
	stm, err := ioutil.ReadFile(fname)
	if err != nil {
		return
	}
	if len(stm) >= 3 && stm[0] == 0xef && stm[1] == 0xbb && stm[2] == 0xbf {
		ret = strings.Split(string(stm[3:]), "\r\n")
	} else if len(stm) >= 2 && stm[0] == 0xff && stm[1] == 0xfe {
		ret = strings.Split(codec.Decode(stm, codec.UTF16LE), "\r\n")
	} else if len(stm) >= 2 && stm[0] == 0xfe && stm[1] == 0xff {
		ret = strings.Split(codec.Decode(stm, codec.UTF16BE), "\r\n")
	} else {
		ret = strings.Split(codec.Decode(stm, codec.C936), "\r\n")
	}
	return
}
