package yaml

import (
	"testing"
	"bytes"
	"fmt"

	"runtime/debug"
)

var parseTests = []struct {
	Input  string
	Output string
}{
	{
		Input:  "key1: val1\n",
		Output: "key1: val1\n",
	},
	{
		Input: "key: nest: val\n",
		Output: "key:\n" +
			"  nest: val\n",
	},
	{
		Input: "a: b: c: d\n" +
			"      e: f\n" +
			"   g: h: i\n" +
			"      j: k\n" +
			"   l: m\n" +
			"n: o\n" +
			"",
		Output: "n: o\n" +
			"a:\n" +
			"  l: m\n" +
			"  b:\n" +
			"    c: d\n" +
			"    e: f\n" +
			"  g:\n" +
			"    h: i\n" +
			"    j: k\n" +
			"",
	},
	{
		Input: "- item\n" +
			"",
		Output: "- item\n" +
			"",
	},
	{
		Input: "- item2\n" +
			"- item1\n" +
			"",
		Output: "- item2\n" +
			"- item1\n" +
			"",
	},
	{
		Input: "- - list1a\n" +
			"  - list1b\n" +
			"- - list2a\n" +
			"  - list2b\n" +
			"",
		Output: "- - list1a\n" +
			"  - list1b\n" +
			"- - list2a\n" +
			"  - list2b\n" +
			"",
	},
	{
		Input: "-   \n" +
			"  - - listA1a\n" +
			"    - listA1b\n" +
			"  - - listA2a\n" +
			"    - listA2b\n" +
			"-\n" +
			"  - - listB1a\n" +
			"    - listB1b\n" +
			"  - - listB2a\n" +
			"    - listB2b\n" +
			"",
		Output: "- - - listA1a\n" +
			"    - listA1b\n" +
			"  - - listA2a\n" +
			"    - listA2b\n" +
			"- - - listB1a\n" +
			"    - listB1b\n" +
			"  - - listB2a\n" +
			"    - listB2b\n" +
			"",
	},
	{
		Input: "  - keyA1a: aaa\n" +
			"    keyA1b: bbb\n" +
			"  - keyA2a: ccc\n" +
			"    keyA2b: ddd\n" +
			"  - keyB1a: eee\n" +
			"    keyB1b: fff\n" +
			"  - keyB2a: ggg\n" +
			"    keyB2b: hhh\n" +
			"",
		Output: "- keyA1a: aaa\n" +
			"  keyA1b: bbb\n" +
			"- keyA2a: ccc\n" +
			"  keyA2b: ddd\n" +
			"- keyB1a: eee\n" +
			"  keyB1b: fff\n" +
			"- keyB2a: ggg\n" +
			"  keyB2b: hhh\n" +
			"",
	},
	{
		Input: "japanese:\n" +
			" - ichi\n" +
			" - ni\n" +
			" - san\n" +
			"french:\n" +
			" - un\n" +
			" - deux\n" +
			" - troix\n" +
			"english:\n" +
			" - one\n" +
			" - two\n" +
			" - three\n" +
			"",
		Output: "english:\n" +
			"  - one\n" +
			"  - two\n" +
			"  - three\n" +
			"french:\n" +
			"  - un\n" +
			"  - deux\n" +
			"  - troix\n" +
			"japanese:\n" +
			"  - ichi\n" +
			"  - ni\n" +
			"  - san\n" +
			"",
	},
}

func TestParse(t *testing.T) {
	/*
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
			}
		}()
	*/
	_ = debug.PrintStack

	for idx, test := range parseTests {
		buf := bytes.NewBufferString(test.Input)
		node, err := Parse(buf)
		if err != nil {
			t.Errorf("parse: %s", err)
		}
		buf.Truncate(0)
		fmt.Fprintf(buf, "%s", node)
		if got, want := buf.String(), test.Output; got != want {
			t.Errorf("---%d---", idx)
			t.Errorf("got: %q:\n%s", got, got)
			t.Errorf("want: %q:\n%s", want, want)
		}
	}
}

var getTypeTests = []struct {
	Value string
	Type  int
	Split int
}{
	{
		Value: "a: b",
		Type:  typMapping,
		Split: 1,
	},
	{
		Value: "- b",
		Type:  typSequence,
		Split: 1,
	},
}

func TestGetType(t *testing.T) {
	for idx, test := range getTypeTests {
		v, s := getType([]byte(test.Value))
		if got, want := v, test.Type; got != want {
			t.Errorf("%d. type(%q) = %s, want %s", idx, test.Value,
				typNames[got], typNames[want])
		}
		if got, want := s, test.Split; got != want {
			got0, got1 := test.Value[:got], test.Value[got:]
			want0, want1 := test.Value[:want], test.Value[want:]
			t.Errorf("%d. split is %s|%s, want %s|%s", idx,
				got0, got1, want0, want1)
		}
	}
}