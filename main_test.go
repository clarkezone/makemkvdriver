package makemkvdriver

import (
	"bufio"
	"strings"
	"testing"
)

func TestMySplitter(t *testing.T) {
	const input = "hello\rhello\n"
	scanner := bufio.NewScanner(strings.NewReader(input))

	scanner.Split(mySplitter)

	for i := 0; i < 2; i++ {
		res := scanner.Scan()
		if !res {
			t.Error("No output")
		}
		result := scanner.Text()

		if result != "hello" {
			t.Error("result:" + result)
		}
	}
}
