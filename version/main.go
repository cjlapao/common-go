package version

import (
	"encoding/json"
	"fmt"

	"github.com/cjlapao/common-go/strcolor"

	"gopkg.in/yaml.v2"
)

// Version Entity
type Version struct {
	Major int
	Minor int
	Build int
	Rev   int
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
	fmt.Printf("*                     Carlos Http Load Tester Tool %v                     *\n", strcolor.GetColorString(strcolor.BrightYellow, v.String()))
	fmt.Printf("*                                                                              *\n")
	fmt.Printf("*  Author:  Carlos Lapao                                                       *\n")
	fmt.Printf("*  License: MIT                                                                *\n")
	fmt.Printf("********************************************************************************\n")
	fmt.Println("")
}

// PrintHeader Prints an Application Version simple text header
func (v *Version) PrintHeader() {
	fmt.Println("Ivanti HTTP Load Tester utility to test services capacity")
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
