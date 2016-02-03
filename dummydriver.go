package makemkvdriver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

type dummyCommand struct {
}

func (r dummyCommand) commandLine() string {
	//go r.readCommandAsync("c:\\users\\james\\desktop\\GoWorkspace\\bin\\dummyoutput.exe")
	return `c:\users\james\desktop\GoWorkspace\bin\dummyoutput.exe`
}

func (ro dummyCommand) parseStdOut(outStream io.Reader, r PercentageSetter) {
	//	regexList := genRegex()

	inputReader := bufio.NewScanner(outStream)

	inputReader.Split(mySplitter)

	percFinder := "([0-9]+)%"
	exr, err := regexp.Compile(percFinder)
	if err != nil {
		log.Fatal(err)
	}
	for inputReader.Scan() {
		line := inputReader.Text()
		result := exr.FindStringSubmatch(line)
		if len(result) == 2 {
			perc, _ := strconv.Atoi(string(result[1])) // first result in slice is whole match not subgroup
			r(perc)
		}
		if err == io.EOF {
			break
		}
		fmt.Printf("%s \n", line)
	}
}

//func genRegex() []regexp.Regexp {
//	exprList := []string{
//		"a",
//		"b",
//		"c"}
//
//	for i := range exprList {
//		percFinder := i
//		exr, err := regexp.Compile(percFinder)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}

func (ro dummyCommand) old_parseStdOut(outStream io.Reader, r PercentageSetter) {
	inputReader := bufio.NewScanner(outStream)

	inputReader.Split(mySplitter)

	percFinder := "([0-9]+)%"
	exr, err := regexp.Compile(percFinder)
	if err != nil {
		log.Fatal(err)
	}
	for inputReader.Scan() {
		line := inputReader.Text()
		result := exr.FindStringSubmatch(line)
		if len(result) == 2 {
			perc, _ := strconv.Atoi(string(result[1])) // first result in slice is whole match not subgroup
			r(perc)
		}
		if err == io.EOF {
			break
		}
		fmt.Printf("%s \n", line)
	}
}
