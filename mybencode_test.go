package mybencode

import (
	"testing"
)

func TestDecodeMap(e *testing.T) {
	{
		a, l, _ := decodeMap([]byte("d2:abi32ee"))
		k, ok := a["ab"]
		if !ok || k.GetType() != IntValue || k.ToInt() != 32 || l != 10 {
			e.Error()
		}
	}
}

func TestDecodeList(e *testing.T) {
	{
		a, l, _ := decodeList([]byte("li3ee"))
		if len(a) != 1 || a[0].GetType() != IntValue || a[0].ToInt() != 3 || l != 5 {
			e.Error()
		}
	}
	{
		a, l, _ := decodeList([]byte("li-33e4:abcde"))
		if !(len(a) == 2 && a[0].GetType() == IntValue && a[0].ToInt() == -33 && a[1].GetType() == StringValue && a[1].ToString() == "abcd" && l == 13) {
			e.Error()
		}
	}
}

func TestDecodeString(e *testing.T) {
	tables := []struct {
		x []byte
		i string
		l int
	}{
		{[]byte("1:a"), "a", 3},
		{[]byte("9:abcdfsfgr"), "abcdfsfgr", 11},
		{[]byte("11:abcdfsfgraa"), "abcdfsfgraa", 14},
	}

	for _, t := range tables {
		value, len, err := decodeString(t.x)
		if value != t.i || t.l != len || err != nil {
			e.Errorf("error: %s, %d, %s, %d", t.i, t.l, value, len)
		}
	}
}

func TestDecodeInt(e *testing.T) {
	tables := []struct {
		x []byte
		i int64
		l int
	}{
		{[]byte("i0e"), 0, 3},
		{[]byte("i2e"), 2, 3},
		{[]byte("i10e"), 10, 4},
		{[]byte("i-1e"), -1, 4},
		{[]byte("i-33e"), -33, 5},
		{[]byte("i123456789e"), 123456789, 11},
	}

	for _, t := range tables {
		value, len, err := decodeInt(t.x)
		if value != t.i || t.l != len || err != nil {
			e.Errorf("error: %d, %d, %d, %d", t.i, t.l, value, len)
		}
	}
}

func TestAny(t *testing.T) {
	const cint int64 = 32
	i := wrapInt{cint}
	var a Any
	a = i
	if a.GetType() != IntValue {
		t.Errorf("a.GetType != IntValue")
	}
	if a.ToInt() != cint {
		t.Errorf("a.ToInt not equal to cint")
	}

	s := wrapString{"abcdefg"}
	var b Any
	b = s
	if b.GetType() != StringValue {
		t.Errorf("b.GetType != IntValue")
	}
	if b.ToInt() != 0 {
		t.Errorf("b.ToInt not equal to 0")
	}
	if b.ToString() != "abcdefg" {
		t.Errorf("b.ToInt not equal to abcdfg")
	}

	var c wrapList
	c.value = make([]Any, 0)
	c.value = append(c.value, a)
	c.value = append(c.value, b)
	c.value = append(c.value, a)

	if len(c.value) != 3 {
		t.Error()
	}

	var e Any = c
	if e.GetType() != ListValue {
		t.Error()
	}
}
