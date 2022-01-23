package todo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Spuxy/snitchify/github"
)

type Todo struct {
	line        int
	description string
	filename    string
}

func New(d string, f string, l int) Todo {
	return Todo{l, d, f}
}

func GetTodosFromFile(file string) {
	g := github.Github{TOKEN: "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"}

	f, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNumber := 0
	var numberLineToScan []int
	for scanner.Scan() {
		line := scanner.Text()
		if ok := findByReg(line); ok {
			todo := trimFormat(line)
			t := New(todo, file, lineNumber)
			fmt.Println(t)
			issue := map[string]interface{}{"title": "issue " + todo, "body": todo}
			_, err := g.SendReq("POST", "https://api.github.com/repos/Spuxy/snitchify/issues", issue)
			if err != nil {
				fmt.Println("Request Error: ", err.Error())
			}
			numberLineToScan = append(numberLineToScan, lineNumber)
		}
		lineNumber++
		fmt.Println("line:", lineNumber, line)
		fmt.Println(numberLineToScan)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func findByReg(line string) bool {
	if ok, _ := regexp.Match("^(.*)TODO(O*): (.*)$", []byte(line)); ok {
		return true
	}
	return false
}

func trimFormat(l string) string {
	s := strings.TrimPrefix(l, "// TODO:")
	s = strings.TrimPrefix(s, "//TODO:")
	s = strings.TrimSpace(s)
	return s
}
