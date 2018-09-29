package torrentclient

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Test_Main(t *testing.T) {
	reader, err := os.Open("./test.torrent")

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

	fmt.Println(dict)

	dict["announce"].Print()
	dict["announce-list"].Print()
	dict["info"].Print()

	reader.Close()
}
