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
	var g github.Github = github.Github{TOKEN: "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"}
	var lineNumber int
	var numberLineToScan []int
	fileFullContent, err := os.ReadFile(file)

	if err != nil {
		panic(err.Error())
	}

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}

	scanner := bufio.NewScanner(f)

	// iteration over scanjner
	for scanner.Scan() {
		line := scanner.Text()

		// takes line and check if exists todo comment
		if ok := findByReg(line); ok {

			// format todo line -> "// TODO" etc..
			todo := trimFormat(line)

			issue := map[string]interface{}{"title": "issue " + todo, "body": todo}
			_, err := g.SendReq("POST", "https://api.github.com/repos/Spuxy/snitchify/issues", issue)
			if err != nil {
				fmt.Println("Request Error: ", err.Error())
			}

			numberLineToScan = append(numberLineToScan, lineNumber)
			fileFullContent = replaceTodo(file, todo, fileFullContent)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// create file
	createdFile, err := os.Create(file)
	defer createdFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(fileFullContent))
	createdFile.WriteString(string(fileFullContent))

	// remove file
	os.Remove(file)
}

// return whole content file in string
func replaceTodo(fileName string, todo string, fileFullContentv2 []byte) []byte {
	// fileFullContent, _ := os.ReadFile(fileName)
	rg, err := regexp.Compile("^(.*)TODO(O*): " + todo)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	return rg.ReplaceAll(fileFullContentv2, []byte("DONE"))
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
