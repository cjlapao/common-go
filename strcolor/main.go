package strcolor

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ColorCode int

const (
	Black ColorCode = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack ColorCode = iota + 82
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func GetColorString(colorCode ColorCode, words ...string) string {
	agentID := os.Getenv("AGENT_ID")
	isPipeline := false
	if len(agentID) != 0 {
		isPipeline = true
	}

	var builder string
	for _, m := range words {
		if len(builder) > 0 {
			builder += " "
		}
		builder += m
	}

	if isPipeline {
		return fmt.Sprintf("\033[%vm%v\033[0m", fmt.Sprint(colorCode), builder)
	}

	return color.New(color.Attribute(colorCode)).Sprint(builder)
}
