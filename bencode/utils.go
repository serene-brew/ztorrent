package bencode

import (
	"fmt"
	"io"
)

// Decode parses Bencoded data into an interface{}
func (d *BencodeDecoder) Decode() (interface{}, error) {
	ch, err := d.readByte()
	if err != nil {
		return nil, err
	}

	switch {
	case ch == 'i': // Integer
		return d.readInt()
	case ch >= '0' && ch <= '9': // String
		return d.readString(ch)
	case ch == 'l': // List
		return d.readList()
	case ch == 'd': // Dictionary
		return d.readDict()
	default:
		return nil, fmt.Errorf("invalid bencode type: %c", ch)
	}
}

func (d *BencodeDecoder) readInt() (int64, error) {
	var result int64
	var isNegative bool

	ch, err := d.readByte()
	if err != nil {
		return 0, err
	}
	if ch == '-' {
		isNegative = true
		ch, err = d.readByte()
		if err != nil {
			return 0, err
		}
	}
	for ch != 'e' {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("invalid integer")
		}
		result = result*10 + int64(ch-'0')
		ch, err = d.readByte()
		if err != nil {
			return 0, err
		}
	}
	if isNegative {
		result = -result
	}
	return result, nil
}

func (d *BencodeDecoder) readString(lengthByte byte) (string, error) {
	length := int(lengthByte - '0')
	for {
		ch, err := d.readByte()
		if err != nil {
			return "", err
		}
		if ch == ':' {
			break
		}
		length = length*10 + int(ch-'0')
	}
	str := make([]byte, length)
	if _, err := io.ReadFull(d.reader, str); err != nil {
		return "", err
	}
	return string(str), nil
}

func (d *BencodeDecoder) readList() ([]interface{}, error) {
	var list []interface{}
	for {
		ch, err := d.readByte()
		if err != nil {
			return nil, err
		}
		if ch == 'e' {
			break
		}
		d.unreadByte()
		item, err := d.Decode()
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (d *BencodeDecoder) readDict() (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	for {
		ch, err := d.readByte()
		if err != nil {
			return nil, err
		}
		if ch == 'e' {
			break
		}
		d.unreadByte()
		key, err := d.Decode()
		if err != nil {
			return nil, err
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		dict[key.(string)] = value
	}
	return dict, nil
}
