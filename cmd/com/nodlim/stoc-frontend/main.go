package main

import (
	"bufio"
	"fmt"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc"
	"io"
	"os"
	"strings"
)

func searchPipeMode() {
	reader := bufio.NewReader(os.Stdin)
	command := strings.Join(os.Args[1:], " ")
	cont := true
	for cont {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			cont = false
		}
		res, searchErr := stoc.SearchString(command, line)
		if searchErr == nil && res {
			fmt.Println(strings.Trim(line, "\n"))
		} else if searchErr != nil {
			fmt.Println(searchErr)
		}
	}
}

func isPiping() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return !(info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0)
}

func main() {
	if isPiping() && len(os.Args) > 1 {
		searchPipeMode()
	}
}
