package mybencode

type ValueType int

const (
	IntValue ValueType = iota
	StringValue
	ListValue
	MapValu
)

type Any interface {
	GetType() ValueType
	ToInt() int64
	ToString() string
	ToList() []Any
	ToMap() map[string]Any
}

type wrapInt struct {
	value int64
}

func (v wrapInt) GetType() ValueType {
	return IntValue
}

func (v wrapInt) ToInt() int64 {
	return v.value
}

func (v wrapInt) ToString() string {
	return ""
}

func (v wrapInt) ToList() []Any {
	return nil
}

func (v wrapInt) ToMap() map[string]Any {
	return nil
}

type wrapString struct {
	value string
}

func (v wrapString) GetType() ValueType {
	return StringValue
}

func (v wrapString) ToInt() int64 {
	return 0
}

func (v wrapString) ToString() string {
	return v.value
}

func (v wrapString) ToList() []Any {
	return nil
}

func (v wrapString) ToMap() map[string]Any {
	return nil
}

type wrapList struct {
	value []Any
}

func (v wrapList) GetType() ValueType {
	return ListValue
}

func (v wrapList) ToInt() int64 {
	return 0
}

func (v wrapList) ToString() string {
	return ""
}

func (v wrapList) ToList() []Any {
	return v.value
}

func (v wrapList) ToMap() map[string]Any {
	return nil
}

type wrapMap struct {
	value map[string]Any
}

func (v wrapMap) GetType() ValueType {
	return MapValu
}

func (v wrapMap) ToInt() int64 {
	return 0
}

func (v wrapMap) ToString() string {
	return ""
}

func (v wrapMap) ToList() []Any {
	return nil
}

func (v wrapMap) ToMap() map[string]Any {
	return v.value
}

type parseError struct {
	info string
}

func (p parseError) Error() string {
	return p.info
}

func decodeItem(b []byte) (Any, int, error) {
	switch {
	case b[0] == byte('i'):
		i, l, err := decodeInt(b)
		return wrapInt{i}, l, err
	case b[0] == byte('l'):
		v, l, err := decodeList(b)
		return wrapList{v}, l, err
	case b[0] == byte('d'):
		v, l, err := decodeMap(b)
		return wrapMap{v}, l, err
	case b[0] >= byte('0') && b[0] <= byte('9'):
		s, l, err := decodeString(b)
		return wrapString{s}, l, err
	}

	return nil, 0, parseError{""}
}

func decodeList(b []byte) ([]Any, int, error) {
	listLen := 1
	r := make([]Any, 0)
	for listLen < len(b) && b[listLen] != byte('e') {
		a, l, err := decodeItem(b[listLen:])
		if err != nil {
			return nil, l, err
		}
		r = append(r, a)
		listLen += l
	}
	return r, listLen + 1, nil
}

func decodeMap(b []byte) (map[string]Any, int, error) {
	mapLen := 1
	r := make(map[string]Any)
	for mapLen < len(b) && b[mapLen] != byte('e') {
		a, l, err := decodeItem(b[mapLen:])
		if err != nil || a.GetType() != StringValue {
			return nil, l, err
		}
		mapLen += l
		v, l, err := decodeItem(b[mapLen:])
		if err != nil {
			return nil, l, err
		}
		mapLen += l
		r[a.ToString()] = v
	}
	return r, mapLen + 1, nil
}

func decodeString(b []byte) (string, int, error) {
	if len(b) < 3 {
		return "", 0, parseError{""}
	}

	var strLen int
	var numLen int
	for i, v := range b {
		if v == byte(':') {
			b = b[i:]
			numLen = i
			break
		} else if v >= byte('0') && v <= byte('9') {
			strLen *= 10
			strLen += int(v - byte('0'))
		} else {
			return "", 0, parseError{""}
		}
	}

	if len(b) < 1+strLen {
		return "", 0, parseError{""}
	}

	return string(b[1 : strLen+1]), strLen + numLen + 1, nil
}

func decodeInt(b []byte) (int64, int, error) {
	if len(b) < 3 {
		return 0, 0, parseError{""}
	}
	var value int64
	b = b[1:]
	isNegtive := false
	for i, v := range b {
		if i == 0 && v == byte('-') {
			isNegtive = true
		} else if v == byte('e') {
			if isNegtive {
				value = 0 - value
			}
			return value, i + 2, nil
		} else if v >= byte('0') && v <= byte('9') {
			value *= 10
			value += int64(v - byte('0'))
		} else {
			break
		}
	}
	return 0, 0, parseError{"unexpected value "}
}
