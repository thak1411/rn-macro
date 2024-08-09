package main

import (
	"fmt"
	"os"

	"github.com/go-vgo/robotgo"
	"golang.org/x/term"
)

type Pos struct {
	x, y int
}

func Menu() {
	fmt.Printf("0. Run\r\n")
	fmt.Printf("1. Record\r\n")
	fmt.Printf("2. Test Record\r\n")
	fmt.Printf("3. Check Recorded Data\r\n")
	fmt.Printf("4. Saving Data & Exit\r\n")
}

func Move(x, y int) {
	robotgo.Move(x, y)
}

func Click(x, y int) {
	robotgo.Move(x, y)
	robotgo.Click()
}

func ChangeInterval(ms int) {
	robotgo.MouseSleep = ms
}

func GetMousePos() (int, int) {
	return robotgo.Location()
}

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func GetChar() byte {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		fmt.Printf("Error reading input: %v\r\n", err)
		return 0
	}
	return buf[0]
}

func Run(data []Pos) {
	ClearTerminal()
	fmt.Printf("Running...\r\n")
	ChangeInterval(50)
	for _, pos := range data {
		fmt.Printf("Click(%d, %d)\r\n", pos.x, pos.y)
		Click(pos.x, pos.y)
	}
	fmt.Printf("Press any key to continue...\r\n")
	GetChar()
}

func Record() []Pos {
	ClearTerminal()
	fmt.Printf("Recording...\r\n - If you want to record to click current Mouse Position, press any key\r\n - If you want to stop recording, press 'q'\r\n")
	ret := make([]Pos, 0)
	for {
		if GetChar() == 'q' {
			break
		}
		x, y := GetMousePos()
		fmt.Printf("Click(%d, %d)\r\n", x, y)
		ret = append(ret, Pos{x, y})
	}
	return ret
}

func TestRecord(data []Pos) {
	ClearTerminal()
	fmt.Printf("Testing...\r\n")
	ChangeInterval(700)
	for _, pos := range data {
		fmt.Printf("Move(%d, %d)\r\n", pos.x, pos.y)
		Move(pos.x, pos.y)
	}
	fmt.Printf("Press any key to continue...\r\n")
	GetChar()
}

func CheckRecordedData(data []Pos) {
	ClearTerminal()
	fmt.Printf("Recorded Data...\r\n")
	for _, pos := range data {
		fmt.Printf("Pos(%d, %d)\r\n", pos.x, pos.y)
	}
	fmt.Printf("Press any key to continue...\r\n")
	GetChar()
}

func Exit(data []Pos) {
	ClearTerminal()
	fmt.Printf("Saving Data & Exit...\r\n")

	// Save Data to text
	file, err := os.Create("recorded_data.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\r\n", err)
		return
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("%d\r\n", len(data)))
	for _, pos := range data {
		file.WriteString(fmt.Sprintf("%d %d\r\n", pos.x, pos.y))
	}
}

func Init() []Pos {
	ret := make([]Pos, 0)
	file, err := os.Open("recorded_data.txt")
	if err != nil {
		fmt.Printf("No Recorded File Exists\r\n")
		return ret
	}
	defer file.Close()

	var n int
	_, err = fmt.Fscanf(file, "%d\r\n", &n)
	if err != nil {
		fmt.Printf("Error reading file: %v\r\n", err)
		return ret
	}

	for i := 0; i < n; i++ {
		var x, y int
		_, err = fmt.Fscanf(file, "%d %d\r\n", &x, &y)
		if err != nil {
			fmt.Printf("Error reading file: %v\r\n", err)
			return ret
		}
		ret = append(ret, Pos{x, y})
	}

	fmt.Printf("Recorded Data Loaded\r\n")
	for _, pos := range ret {
		fmt.Printf("Pos(%d, %d)\r\n", pos.x, pos.y)
	}
	return ret
}

func main() {
	ChangeInterval(50)

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("Error setting terminal to raw mode: %v\r\n", err)
		return
	}
	defer term.Restore(fd, oldState)

	state := 0
	var _ = state
	recordedData := Init()

	for {
		Menu()
		state = 0

		switch GetChar() {
		case '0':
			Run(recordedData)
		case '1':
			recordedData = Record()
		case '2':
			TestRecord(recordedData)
		case '3':
			CheckRecordedData(recordedData)
		case '4':
			Exit(recordedData)
			return
		}
		ClearTerminal()
	}
}
