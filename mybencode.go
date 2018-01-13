package Bencode

import (
	"strconv"
)

type Handler interface {
	GetByPos(int) Handler
	GetByKey(string) Handler
	ToInt64() (int64, error)
	ToString() (string, error)
	ToList() ([]Any, error)
	ToMap() (map[string]Any, error)
}

type handler struct {
	any Any
	err error
}

//NewHandler make new Handler
func NewHandler(a Any) Handler {
	return handler{a, nil}
}

func (h handler) ToList() ([]Any, error) {
	if h.any == nil {
		return nil, &bencodeError{"Any is nil"}
	}

	if h.any.GetType() != ListValue {
		return nil, &bencodeError{"Any is not list"}
	}

	return h.any.ToList(), nil
}

func (h handler) ToMap() (map[string]Any, error) {
	if h.any == nil {
		return nil, &bencodeError{"Any is nil"}
	}

	if h.any.GetType() != MapValue {
		return nil, &bencodeError{"Any is not map"}
	}

	return h.any.ToMap(), nil
}

func (h handler) GetByPos(pos int) Handler {
	if h.any == nil {
		return handler{nil, h.err}
	}

	if h.any.GetType() != ListValue {
		return handler{nil, &bencodeError{"Any is not list"}}
	}

	list := h.any.ToList()
	if len(list) <= pos || pos < 0 {
		return handler{nil, &bencodeError{"pos out of range"}}
	}

	return handler{list[pos], nil}
}

func (h handler) GetByKey(key string) Handler {
	if h.any == nil {
		return handler{nil, h.err}
	}

	if h.any.GetType() != MapValue {
		return handler{nil, &bencodeError{"Any is not list"}}
	}

	dic := h.any.ToMap()
	if v, ok := dic[key]; ok {
		return handler{v, nil}
	}

	return handler{nil, &bencodeError{"key not found"}}
}

func (h handler) ToInt64() (int64, error) {
	if h.any == nil {
		return 0, h.err
	}

	if h.any.GetType() != IntValue {
		return 0, &bencodeError{"Any is not int"}
	}

	return h.any.ToInt(), nil
}

func (h handler) ToString() (string, error) {
	if h.any == nil {
		return "", h.err
	}

	if h.any.GetType() != StringValue {
		return "", &bencodeError{"Any is not string"}
	}

	return h.any.ToString(), nil
}

type ValueType int

const (
	IntValue ValueType = iota
	StringValue
	ListValue
	MapValue
)

type Any interface {
	GetType() ValueType
	ToInt() int64
	ToString() string
	ToList() []Any
	ToMap() map[string]Any
}

type Writer interface {
	StartDic()
	EndDic()
	StartList()
	EndList()
	AppendInt64(int64)
	AppendString(string)
	GetBytes() []byte
}

func NewWriter() Writer {
	var w writer
	w.buf = make([]byte, 0)
	return &w
}

type writer struct {
	buf []byte
}

func (w *writer) GetBytes() []byte {
	return w.buf
}
func (w *writer) StartDic() {
	w.buf = append(w.buf, byte('d'))
}

func (w *writer) EndDic() {
	w.buf = append(w.buf, byte('e'))
}

func (w *writer) StartList() {
	w.buf = append(w.buf, byte('l'))
}

func (w *writer) EndList() {
	w.buf = append(w.buf, byte('e'))
}

func (w *writer) AppendInt64(i int64) {
	s := strconv.FormatInt(i, 10)
	w.buf = append(w.buf, byte('i'))
	w.buf = append(w.buf, s...)
	w.buf = append(w.buf, byte('e'))
}

func (w *writer) AppendString(s string) {
	strLen := len(s)
	byteStrLen := []byte(strconv.Itoa(strLen))
	w.buf = append(w.buf, byteStrLen...)
	w.buf = append(w.buf, byte(':'))
	w.buf = append(w.buf, []byte(s)...)
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
	return MapValue
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

type bencodeError struct {
	info string
}

func (p *bencodeError) Error() string {
	return p.info
}

func DecodeItem(b []byte) (Any, int, error) {
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

	return nil, 0, &bencodeError{""}
}

func decodeList(b []byte) ([]Any, int, error) {
	listLen := 1
	r := make([]Any, 0)
	for listLen < len(b) && b[listLen] != byte('e') {
		a, l, err := DecodeItem(b[listLen:])
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
		a, l, err := DecodeItem(b[mapLen:])
		if err != nil || a.GetType() != StringValue {
			return nil, l, err
		}
		mapLen += l
		v, l, err := DecodeItem(b[mapLen:])
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
		return "", 0, &bencodeError{""}
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
			return "", 0, &bencodeError{""}
		}
	}

	if len(b) < 1+strLen {
		return "", 0, &bencodeError{""}
	}

	return string(b[1 : strLen+1]), strLen + numLen + 1, nil
}

func decodeInt(b []byte) (int64, int, error) {
	if len(b) < 3 {
		return 0, 0, &bencodeError{""}
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
	return 0, 0, &bencodeError{"unexpected value "}
}
