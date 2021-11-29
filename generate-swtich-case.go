package main

import (
	"strings"

	"github.com/cjlapao/common-go/helper"
)

func main() {
	file := helper.GetFlagValue("file", "")
	reverse := helper.GetFlagSwitch("reverse", false)

	if helper.FileExists(file) {
		println("File found reading it")
		result, err := helper.ReadFromFile(file)
		if err != nil {
			println(" there was an error reading the file")
		}

		lines := strings.Split(string(result), "\r\n")
		for _, line := range lines {
			columns := strings.Split(line, ",")
			if len(columns) >= 2 {
				if !reverse {
					println("case " + strings.Trim(columns[0], "") + ":")
					println("  return \"" + strings.Trim(columns[1], "") + "\"")
				} else {
					println("case \"" + strings.Trim(columns[1], "") + "\":")
					println("  return " + strings.Trim(columns[0], "") + "")

				}
			}
		}
	}

	println(file)
}
