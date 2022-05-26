package version

import (
	"encoding/json"
	"fmt"

	"github.com/cjlapao/common-go/strcolor"

	"gopkg.in/yaml.v3"
)

// Version Entity
type Version struct {
	Name    string
	Author  string
	License string
	Major   int
	Minor   int
	Build   int
	Rev     int
}

// FormatedVersion Entity
type FormatedVersion struct {
	Version string `json:"version" yaml:"version"`
}

// OutputFormat Version output format enum
type OutputFormat int

// Version Output format enum definition
const (
	JSON OutputFormat = iota
	Yaml
)

var appVersion *Version

// Get Creates a new Version for the application
func Get(v ...int) *Version {
	if appVersion != nil {
		return appVersion
	}
	result := Version{
		Major: 0,
		Minor: 0,
		Build: 0,
		Rev:   0,
	}

	for i, versionSegment := range v {
		switch i {
		case 0:
			result.Major = versionSegment
		case 1:
			result.Minor = versionSegment
		case 2:
			result.Build = versionSegment
		case 3:
			result.Rev = versionSegment
		}
	}

	appVersion = &result
	return appVersion
}

func (v *Version) String() string {
	return fmt.Sprint(v.Major) + "." + fmt.Sprint(v.Minor) + "." + fmt.Sprint(v.Build) + "." + fmt.Sprint(v.Rev)
}

// PrintAnsiHeader Prints a Application Version Ansi Header
func (v *Version) PrintAnsiHeader() {
	fmt.Printf("********************************************************************************\n")
	fmt.Printf("*                                                                              *\n")
	if v.Name != "" {
		name := v.generateMiddle(fmt.Sprintf("%v %v", v.Name, v.String()))
		fmt.Printf("*%v*\n", name)
	}
	fmt.Printf("*                                                                              *\n")
	if v.Author != "" {
		fmt.Printf("*%v*\n", v.generateLeft(fmt.Sprintf(" Author: %v", v.Author)))
	}
	if v.License != "" {
		fmt.Printf("*%v*\n", v.generateLeft(fmt.Sprintf(" License: %v ", v.License)))
	}
	fmt.Printf("********************************************************************************\n")
	fmt.Println("")
}

// PrintHeader Prints an Application Version simple text header
func (v *Version) PrintHeader() {
	header := ""
	if v.Name != "" {
		header = fmt.Sprintf("%v %v", v.Name, strcolor.GetColorString(strcolor.BrightYellow, v.String()))
	} else {
		header = fmt.Sprintf("Version %v", strcolor.GetColorString(strcolor.BrightYellow, v.String()))
	}

	fmt.Printf("%v\n", header)
	fmt.Println("")
}

// PrintVersion Prints the application version in the desired format
func (v *Version) PrintVersion(format int) interface{} {
	formatedVersion := FormatedVersion{
		Version: v.String(),
	}

	switch format {
	case 0:
		jsonString, _ := json.MarshalIndent(formatedVersion, "", "  ")
		return string(jsonString)
	case 1:
		yamlString, _ := yaml.Marshal(v)
		return string(yamlString)
	default:
		jsonString, _ := json.MarshalIndent(formatedVersion, "", "  ")
		return string(jsonString)
	}
}

func (v *Version) generateEmpty(value int) string {
	result := ""
	for i := 0; i < value; i++ {
		result += " "
	}
	return result
}

func (v *Version) generateMiddle(value string) string {
	if value != "" {
		emptyCount := 0
		if len(value) < 78 {
			emptyCount = (78 - len(value)) / 2
		}

		if emptyCount > 0 {
			emptySpace := v.generateEmpty(emptyCount)
			value = fmt.Sprintf("%v%v%v", emptySpace, value, emptySpace)
			if len(value) < 78 {
				emptySpace = v.generateEmpty(78 - len(value))
				value += emptySpace
			}
			if len(value) > 78 {
				value = value[:78]
			}
		}

		return value
	}

	return ""
}

func (v *Version) generateLeft(value string) string {
	if value != "" {
		emptyCount := 0
		if len(value) < 78 {
			emptyCount = (78 - len(value))
		}

		if emptyCount > 0 {
			emptySpace := v.generateEmpty(emptyCount)
			value = fmt.Sprintf("%v%v", value, emptySpace)
			if len(value) < 78 {
				emptySpace = v.generateEmpty(78 - len(value))
				value += emptySpace
			}
			if len(value) > 78 {
				value = value[:78]
			}
		}

		return value
	}

	return ""
}

func (v *Version) generateRight(value string) string {
	if value != "" {
		emptyCount := 0
		if len(value) < 78 {
			emptyCount = (78 - len(value))
		}

		if emptyCount > 0 {
			emptySpace := v.generateEmpty(emptyCount)
			value = fmt.Sprintf("%v%v", emptySpace, value)
			if len(value) < 78 {
				emptySpace = v.generateEmpty(78 - len(value))
				value += emptySpace
			}
			if len(value) > 78 {
				value = value[:78]
			}
		}

		return value
	}

	return ""
}
