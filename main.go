package makemkvdriver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	//"strings"
)

type runningState int

type PercentageSetter func(int)
type StatusSetter func(string)

const (
	ready   runningState = iota
	running              = iota
)

type engineState struct {
	perc          int
	currentState  runningState
	mylock        sync.Mutex
	currentDriver consoleDriver
}

type Ripper struct {
	engineState //anonymous struct member
}

func Create() *Ripper {
	var instance Ripper
	instance.currentState = ready
	return &instance
}

type consoleDriver interface {
	commandLine() string
	parseStdOut(inStream io.Reader, r PercentageSetter)
}

type consoleParseStep interface {
	prepare()
	process(source string, p PercentageSetter, s StatusSetter) bool
}

func (r *Ripper) Start() bool {
	var canStart bool
	r.mylock.Lock()
	if r.currentState == ready {
		r.currentState = running
		canStart = true
		r.currentDriver = dummyCommand{}
		//r.currentDriver = roboCommand{}
	}
	r.mylock.Unlock()

	if canStart {
		fmt.Println("staring makemkv engine.")
		//readCommand("C:\\Program Files (x86)\\MakeMKV\\makemkvcon.exe")
		go r.readCommandAsync(r.currentDriver.commandLine())
		return true
	} else {
		return false
	}
}

func (r *Ripper) setPercentage(perc int) {
	r.mylock.Lock()
	r.perc = perc
	r.mylock.Unlock()
}

func (r *Ripper) GetStatus() string {
	return fmt.Sprintf("State: %s Percentage %v", r.currentState, r.perc)
}

func (r *Ripper) readCommandAsync(myCommand ...string) string {
	fmt.Println(myCommand[0])
	//cmd := exec.Command("cmd", "/C", myCommand[0], "-r")
	cmd := exec.Command("cmd", "/C", myCommand[0])
	outStream, readError := cmd.StdoutPipe()
	if readError != nil {
		log.Fatal(readError)
	}
	err := cmd.Start()
	if err != nil {
		//log.Fatal(err)
	}
	waitchan := make(chan bool)
	go waitForCommand(cmd, waitchan)

	theFunc := func(value int) {
		r.setPercentage(value)
	}
	r.currentDriver.parseStdOut(outStream, theFunc)
	//consoleRecorder(outStream, r)
	log.Println("UI thread wait start")
	<-waitchan
	fmt.Println("UI thread wait complete")

	r.mylock.Lock()
	r.currentState = ready
	r.mylock.Unlock()

	return ""
}

func consoleRecorder(inStream io.Reader, r *Ripper) {
	fmt.Println("Recording")
	out, err := os.Create("foo")
	defer out.Close()
	if err != nil {
		fmt.Println("Error opening out")
		log.Fatal(err)
	}
	_, err2 := io.Copy(out, inStream)
	if err2 != nil {
		fmt.Println("Error with copy")
	} else {
		fmt.Println("Copy succeeded")
	}
}

func waitForCommand(c *exec.Cmd, w chan bool) {
	c.Wait()
	fmt.Println("Done waiting")
	w <- true
	fmt.Println("signaled to UIthread")
}

func readCommand(myCommand ...string) string {
	fmt.Println(myCommand[0])
	//cmd := exec.Command("cmd", myCommand...)
	cmd := exec.Command("cmd", "/C", myCommand[0], "-r")
	//cmd := exec.Command("cmd", "/C", "dir")
	//cmd := exec.Command(myCommand)
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var errorBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errorBuf
	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String()+errorBuf.String())

	return ""
}

func mySplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	j := bytes.IndexByte(data[0:], '\n')
	k := bytes.IndexByte(data[0:], '\r')
	// -1, 3; 2, -1; 2, 3

	if j >= 0 && k >= 0 {
		value := min(j, k)
		return value + 1, data[0:value], nil
	} else {
		value := max(j, k)
		if value >= 0 {
			return value + 1, data[0:value], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
