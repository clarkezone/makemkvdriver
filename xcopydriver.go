package makemkvdriver

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strconv"
)

//"c:\\users\\james\\desktop\\GoWorkspace\\bin\\dummyoutput.exe"
type roboCommand struct {
}

func (r roboCommand) commandLine() string {
	//go r.readCommandAsync("robocopy C:\\Users\\james\\Desktop\\iPlayerRecordings C:\\Users\\james\\Desktop\\iPlayerRecordings2")
	return `robocopy C:\Users\james\Desktop\iPlayerRecordings C:\Users\james\Desktop\iPlayerRecordings2`
}

func (ro roboCommand) parseStdOut(outStream io.Reader, r PercentageSetter, s StatusSetter) {
	inputReader := bufio.NewScanner(outStream)

	inputReader.Split(mySplitter)

	st := getSteps()
	var cs consoleParseStep
	si := 0 // step index

	for inputReader.Scan() {
		line := inputReader.Text()

		// if current step is nil try and find another one
		if cs == nil {
			if si < len(st) {
				cs = st[si]
				cs.prepare()
			}
		}
		if cs != nil {
			res := cs.process(line, r, s)

			// if the step completed, reset accumulated input and current step
			if res {
				cs = nil
				si++
			}
		}
	}
}

func getSteps() []consoleParseStep {
	steps := []consoleParseStep{&preambleStep{}}
	return steps
}

////////////////////////////////////////////////////////////////
// Preamble Step
////////////////////////////////////////////////////////////////

type preambleStep struct {
	search     string
	stepRegex  *regexp.Regexp
	regexList  []searchLineDef
	regexIndex int
}

type searchLineDef struct {
	regex    string
	matchDef string
}

func (ps *preambleStep) prepare() {
	//	ps.search = "([0-9]+)%"
	//	var err error
	//	ps.theRegex, err = regexp.Compile(ps.search)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	ps.regexList = []searchLineDef{
		searchLineDef{"-{27}", ""},
		searchLineDef{`\s*ROBOCOPY`, ""},
		searchLineDef{"-{27}", ""},
		searchLineDef{"", ""},
		searchLineDef{".*", ""},
		searchLineDef{`\s*Source\s*:\s+(.*)`, ""},
		searchLineDef{`\s*Dest\s*:\s+(.*)`, ""},
	}
}

func (ps *preambleStep) process(line string, p PercentageSetter, s StatusSetter) bool {
	searchFor := ps.regexList[ps.regexIndex].regex
	if ps.stepRegex == nil {
		ps.stepRegex, _ = regexp.Compile(searchFor)
	}
	//format := ps.regexList[ps.regexIndex].matchDef
	var match bool
	var captures []string
	if line == "" {
		match = line == searchFor
	} else {
		captures = ps.stepRegex.FindStringSubmatch(line)
		if len(captures) > 0 {
			match = true
		}
	}
	if match {
		s("Found:" + line)
		ps.regexIndex++
		ps.stepRegex = nil
		if ps.regexIndex < len(ps.regexList) {
			return false
		} else {
			return true
		}
	} else {
		s("Expecting:" + ps.regexList[ps.regexIndex].regex + " found:" + line + ":")
		return false
	}
	return false
}

//func (ps *preambleStep) _old_process(line string, p PercentageSetter, s StatusSetter) bool {
//	searchFor := ps.regexList[ps.regexIndex].regex
//	if ps.stepRegex == nil {
//		ps.stepRegex, _ = regexp.Compile(searchFor)
//	}
//	format := ps.regexList[ps.regexIndex].matchDef
//	var match bool
//	var captures []string
//	if line == "" {
//		match = line == searchFor
//	} else {
//		captures = ps.stepRegex.FindStringSubmatch(line)
//	}
//	if match {
//		s("Found:" + line)
//		ps.regexIndex++
//		ps.stepRegex = nil
//		if ps.regexIndex < len(ps.regexList) {
//			return false
//		} else {
//			return true
//		}
//	} else {
//		s("Expecting:" + ps.regexList[ps.regexIndex].regex + " found:" + line + ":")
//		return false
//	}
//	return false
//}

////////////////////////////////////////////////////////////////
// Progress Step
////////////////////////////////////////////////////////////////

type progressStep struct {
	search   string
	theRegex *regexp.Regexp
}

func (ps *progressStep) prepare() {
	ps.search = "([0-9]+)%"
	var err error
	ps.theRegex, err = regexp.Compile(ps.search)
	if err != nil {
		log.Fatal(err)
	}
}

func (ps *progressStep) process(line string, p PercentageSetter, s StatusSetter) bool {
	s(line)
	result := ps.theRegex.FindStringSubmatch(line)
	if len(result) == 2 {
		perc, _ := strconv.Atoi(string(result[1])) // first result in slice is whole match not subgroup
		p(perc)
		if perc == 100 {
			return true
		}
	}
	return false
}
