package bencode

import (
	"bufio"
	"errors"
	"log"
	"regexp"
	"strconv"
)

// Bencode structure
type Bencode struct {
}

// BNode can be a string, integer, list or a dictionary
type BNode struct {
	Type BencodeType
	Node interface{}
}

// BString is a string
type BString string

// BInteger is an integer
type BInteger int

// BList is a list of BNodes
type BList []*BNode

// BDict is a dictionary of BNodes
type BDict map[string]*BNode

// BencodeType type
type BencodeType uint

// BencodeType Constants
const (
	BencodeUndefined BencodeType = 0
	BencodeString    BencodeType = 1
	BencodeInteger   BencodeType = 2
	BencodeList      BencodeType = 3
	BencodeDict      BencodeType = 4
)

// // NewBencode creates a new bencode reader
// func NewBencode(r io.Reader) *BencodeReader {
// 	b := &BencodeReader{
// 		reader: r,
// 	}

// 	return b
// }

// func (b *BencodeReader)

// BencodeRead ok
func BencodeRead(r *bufio.Reader) (*BNode, error) {
	return parse(r)
}

// GetString returns the string in the node
func (b *BNode) GetString() (BString, error) {
	if b.Type != BencodeString {
		return "", errors.New("not a string")
	}
	s, ok := b.Node.(*BString)
	if !ok {
		return "", errors.New("could not cast")
	}
	return *s, nil
}

// GetInteger returns the integer in the node
func (b *BNode) GetInteger() (BInteger, error) {
	if b.Type != BencodeInteger {
		return 0, errors.New("not an integer")
	}
	i, ok := b.Node.(*BInteger)
	if !ok {
		return 0, errors.New("could not cast")
	}
	return *i, nil
}

// GetList returns the list in the node
func (b *BNode) GetList() (BList, error) {
	if b.Type != BencodeList {
		return nil, errors.New("not a list")
	}
	l, ok := b.Node.(*BList)
	if !ok {
		return nil, errors.New("could not cast")
	}
	return *l, nil
}

// GetDict returns the dict in the node
func (b *BNode) GetDict() (BDict, error) {
	if b.Type != BencodeDict {
		return nil, errors.New("not a dict")
	}
	d, ok := b.Node.(*BDict)
	if !ok {
		return nil, errors.New("could not cast")
	}
	return *d, nil
}

// Print func
func (b *BNode) Print() {
	switch b.Type {
	case BencodeString:
		s, ok := b.Node.(*BString)
		if !ok {
			log.Panicln("could not cast node")
		}
		log.Println(*s)
		break
	case BencodeInteger:
		i, ok := b.Node.(*BInteger)
		if !ok {
			log.Panicln("could not cast node")
		}
		log.Println(*i)
		break
	case BencodeList:
		l, ok := b.Node.(*BList)
		if !ok {
			log.Panicln("could not cast node")
		}
		for _, v := range *l {
			v.Print()
		}
		break
	case BencodeDict:
		d, ok := b.Node.(*BDict)
		if !ok {
			log.Panicln("could not case node")
		}
		for k, v := range *d {
			log.Println(k)
			if v.Type == BencodeList || v.Type == BencodeDict {
				v.Print()
			}
		}
	}
}

func parseString(r *bufio.Reader) (*BString, error) {
	lstr, err := r.ReadString(byte(':'))
	if err != nil {
		return nil, err
	}
	l64, err := strconv.ParseUint(lstr[:len(lstr)-1], 10, 32)
	if err != nil {
		return nil, err
	}
	l := uint(l64)
	buffer := make([]byte, l)
	err = readUntilMax(r, int(l), buffer)
	if err != nil {
		return nil, err
	}
	s := BString(buffer)
	return &s, nil
}

func parseInteger(r *bufio.Reader) (*BInteger, error) {
	istr, err := r.ReadString('e')
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile("i-?00+e|i-0e")
	if err != nil {
		return nil, err
	}
	if reg.MatchString(istr) {
		return nil, errors.New("invalid integer")
	}
	i64, err := strconv.ParseInt(istr[1:len(istr)-1], 10, 64)
	i := BInteger(i64)
	return &i, nil
}

func parseList(r *bufio.Reader) (*BList, error) {
	err := discardSafely(r, 1)
	if err != nil {
		return nil, err
	}
	list := make(BList, 0)
	for {
		p, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if p[0] == byte('e') {
			err := discardSafely(r, 1)
			if err != nil {
				return nil, err
			}
			return &list, nil
		}
		node, err := parse(r)
		if err != nil {
			return nil, err
		}
		list = append(list, node)
	}
}

func parseDict(r *bufio.Reader) (*BDict, error) {
	err := discardSafely(r, 1)
	if err != nil {
		return nil, err
	}
	dict := make(BDict)
	for {
		p, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if p[0] == byte('e') {
			err := discardSafely(r, 1)
			if err != nil {
				return nil, err
			}
			return &dict, nil
		}
		key, err := parseString(r)
		if err != nil {
			return nil, err
		}
		val, err := parse(r)
		if err != nil {
			return nil, err
		}
		dict[string(*key)] = val
	}
}

func parse(r *bufio.Reader) (*BNode, error) {

	ttype, err := getType(r)

	if err != nil {
		return nil, err
	}

	ret := &BNode{}

	switch ttype {
	case BencodeString:
		s, err := parseString(r)
		if err != nil {
			return nil, err
		}
		ret.Type = BencodeString
		ret.Node = s
		break
	case BencodeInteger:
		i, err := parseInteger(r)
		if err != nil {
			return nil, err
		}
		ret.Type = BencodeInteger
		ret.Node = i
		break
	case BencodeList:
		l, err := parseList(r)
		if err != nil {
			return nil, err
		}
		ret.Type = BencodeList
		ret.Node = l
		break
	case BencodeDict:
		d, err := parseDict(r)
		if err != nil {
			return nil, err
		}
		ret.Type = BencodeDict
		ret.Node = d
		break
	default:
		return nil, errors.New("type undefined")
	}

	return ret, nil
}

func discardSafely(r *bufio.Reader, c int) error {
	iskip, err := r.Discard(c)
	if err != nil {
		return err
	}
	if iskip != c {
		return errors.New("discard: invalid number of bytes discarded")
	}
	return nil
}

func readUntilMax(r *bufio.Reader, c int, buffer []byte) error {
	i := 0
	for {
		n, err := r.Read(buffer[i:])
		if err != nil {
			return err
		}
		i += n
		if i == c {
			return nil
		}
	}
}

func getType(r *bufio.Reader) (BencodeType, error) {
	b, err := r.Peek(1)
	if err != nil {
		return BencodeUndefined, err
	}
	switch b[0] {
	case byte('i'):
		return BencodeInteger, nil
	case byte('l'):
		return BencodeList, nil
	case byte('d'):
		return BencodeDict, nil
	default:
		if b[0] >= byte('0') && b[0] <= byte('9') {
			return BencodeString, nil
		}
		return BencodeUndefined, errors.New("type undefined")
	}
}
