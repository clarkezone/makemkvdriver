package makemkvdriver

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	t.Log("TestParser test")
	const input = "10%\n20%"
	stream := strings.NewReader(input)
	dc := dummyCommand{}

	reader := func(value int) {
		t.Log(value)
	}
	dc.parseStdOut(stream, reader)

	//	for i := 0; i < 2; i++ {
	//		res := scanner.Scan()
	//		if !res {
	//			t.Error("No output")
	//		}
	//		result := scanner.Text()
	//
	//		if result != "hello" {
	//			t.Error("result:" + result)
	//		}
	//	}
}
