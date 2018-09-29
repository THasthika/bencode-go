package bencode

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_Main(t *testing.T) {
	reader, err := os.Open("./test2.torrent")

	if err != nil {
		t.Error(err)
	}

	// // ret, err := BencodeRead(bufio.NewReader(strings.NewReader("li324ei12412e5:hello8:tharindue")))
	// ret, err := BencodeRead(bufio.NewReader(reader))

	// if err != nil {
	// 	t.Error(err)
	// }

	// dict, err := ret.GetDict()

	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(dict)

	// dict["announce"].Print()
	// dict["announce-list"].Print()
	// dict["info"].Print()

	reader.Close()
}

func Test_Tracker(t *testing.T) {
	reader, err := os.Open("./test2.torrent")

	if err != nil {
		t.Error(err)
	}

	// ret, err := BencodeRead(bufio.NewReader(strings.NewReader("li324ei12412e5:hello8:tharindue")))
	ret, err := BencodeRead(bufio.NewReader(reader))

	if err != nil {
		t.Error(err)
	}

	dict, err := ret.GetDict()

	if err != nil {
		t.Error(err)
	}

	// dict["announce"].Print()
	// dict["announce-list"].Print()
	dict.Get("info").Print()

	infodict, err := dict.Get("info").GetDict()

	if err != nil {
		t.Error(err)
	}

	log.Println(infodict.Get("name").GetString())
	log.Println(infodict.Get("piece length").GetInteger())

	s, err := dict.Get("info").GetBencode()
	if err != nil {
		t.Error(err)
	}

	h := sha1.New()

	io.WriteString(h, s)

	fmt.Printf("%x\n", h.Sum(nil))

	reader.Close()
}

func Test_Bencodeding(t *testing.T) {

	ret, err := BencodeRead(bufio.NewReader(strings.NewReader("d5:hello5:hello5:helloi124124ee")))

	if err != nil {
		t.Error(err)
	}

	log.Println(ret.GetBencode())

}
