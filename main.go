package main

import (
	"fmt"
	"os"
	"strings"
)

var AllowedParams [2]string = [2]string{"reported", "unreported"}
var PredeclaredParams map[string]string = map[string]string{"unreportred": "", "reported": ""}

func main() {
	rslt := os.Args[1]
	fmt.Println(rslt)
	if len(rslt) > 1 {
		switch os.Args[1] {
		case "list":
			sanitzedParam, err := sanitaze(os.Args[2:]) // all args from position 2 until the end exlueing
			throwError(err)
			params, err := parse(sanitzedParam, AllowedParams)
			throwError(err)
			fmt.Println(params["reported"])
			_, isReported := params["reported"]
			_, isUnreported := params["unreported"]
		}
	}
}

// FIX: error handling not retuning value semantic copy of value
func parse(states []string, defaultPredefineds [2]string) (map[string]string, error) {
	var isOkay bool
	params := make(map[string]string)
	for _, state := range states {
		isOkay = false
		for _, defaultPredefined := range defaultPredefineds {
			if state != defaultPredefined {
				continue
			}
			isOkay = true
			params[state] = state
		}
		if !isOkay {
			return map[string]string{}, fmt.Errorf("Unlnown flag %s", state)
		}
	}
	return params, nil
}

func sanitaze(flags []string) ([]string, error) {
	var currentFlag string
	var sanitzedParams []string
	for _, flag := range flags {
		if strings.HasPrefix(flag, "--") {
			currentFlag = flag[2:] // 2 bytes of sting equels --
			sanitzedParams = append(sanitzedParams, currentFlag)
		} else if strings.HasPrefix(flag, "-") {
			currentFlag = flag[1:] // 1 bytes of sting equels --
			sanitzedParams = append(sanitzedParams, currentFlag)
		} else {
			return []string{}, fmt.Errorf("Flag has to be declared by - or --")
		}
	}
	return sanitzedParams, nil
}

func throwError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
