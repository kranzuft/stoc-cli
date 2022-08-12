package main

import (
	"bufio"
	"fmt"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/search_error"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/types"
	"io"
	"math"
	"os"
	"runtime/debug"
	"strings"
)

func returnErrorMessageStr(rawData string, errorInfo search_error.PosError) {
	returnErrorMessageRune([]rune(rawData), errorInfo)
}

func returnErrorMessageRune(rawData []rune, errorInfo search_error.PosError) {
	fmt.Println("ERROR: " + errorInfo.Error())
	if len(rawData) > 0 {
		lenIndexDiff := len(rawData) - errorInfo.GetPos()
		start := int(math.Min(20, float64(errorInfo.GetPos())))
		end := int(math.Min(20, float64(lenIndexDiff)))
		fmt.Println(string(rawData[errorInfo.GetPos()-start : errorInfo.GetPos()+end]))
		fmt.Println(strings.Repeat(" ", start) + "^")
		debug.PrintStack()
	}
}

// searchPipeMode searches a condition command on piped input
// Consider for future what if a filename is piped?
func searchPipeMode(command string) {
	reader := bufio.NewReader(os.Stdin)
	prepped, err := stoc.LexIntoTokens(types.DefaultTokensDefinition, command)
	if err == nil {
		cont := true
		for cont {
			line, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				cont = false
			}

			res := stoc.SearchTokens(prepped, line)
			if res {
				fmt.Println(strings.Trim(line, "\n"))
			}
		}
	} else {
		returnErrorMessageStr(command, err)
	}
}

// searchBasicMode searches a condition command on a target.
// The target is split into lines and each line is searched for the condition.
// If the line matches the condition, it is output.
func searchBasicMode(command string, target string) {
	var fileErr error
	target, fileErr = loadIfFile(target)

	if fileErr == nil {
		prepped, err := stoc.LexIntoTokens(types.DefaultTokensDefinition, command)
		if err == nil {
			lines := strings.Split(target, "\n")
			for _, line := range lines {
				res := stoc.SearchTokens(prepped, line)
				if res {
					fmt.Println(target)
				}
			}
		} else {
			returnErrorMessageStr(command, err)
		}
	} else {
		fmt.Println("Error loading file: " + fileErr.Error())
	}
}

// loadIfFile check if the argument is a valid file that exists,
// and if so then load the file into memory as a string
// WIP needs fixing
func loadIfFile(potentialFile string) (string, error) {
	// open file
	file, err := os.Open(potentialFile)
	if err != nil {
		// ignore if the file does not exist
		// and simply return
		if err == os.ErrNotExist {
			return potentialFile, nil
		}
		return "", err
	}

	// make sure we close the file
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s", err)
		}
	}(file)

	// now read the file content into a string
	var fileContent string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent += scanner.Text() + "\n"
	}
	return fileContent, nil
}

// isPiping returns true if the program is being piped by another program
func isPiping() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return !(info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0)
}

func main() {
	if len(os.Args) > 1 {
		if isPiping() {
			command := strings.Join(os.Args[1:], " ")
			searchPipeMode(command)
		} else {
			// WIP subject to change
			if len(os.Args) == 5 && os.Args[len(os.Args)-4] == "-c" && os.Args[len(os.Args)-2] == "-t" {
				// Given a situation like -c "foo" -t "baz"
				// the argument between -c and -t is treated as the conditional and the argument after -t as the target to search
				// it requires quotes to be used, unless the target or conditional is a file
				// currently required to be at the end of the command
				// In future, perhaps check if flags are present, and ones that require arguments have one between them?
				// Otherwise, use simple mode. Or otherwise use a starter flag, such as -F (capitalised f) to indicate flag mode
				conditional := os.Args[len(os.Args)-3]
				target := os.Args[len(os.Args)-1]
				searchBasicMode(conditional, target)
			} else if len(os.Args) > 2 {
				// otherwise presume that the last argument is the target to search
				// and everything before is the conditional
				conditional := strings.Join(os.Args[1:len(os.Args)-1], " ")
				target := os.Args[len(os.Args)-1]
				searchBasicMode(conditional, target)
			}
		}
	}
}
