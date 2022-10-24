package helper

import (
	"os"
	"strings"
)

const (
	FlagPrefix         = "--"
	FlagValueSeparator = "="
)

func GetArgumentAt(index int) string {
	args := os.Args
	index += 1
	argsLen := len(args) - 1
	if len(args) > 0 && index <= argsLen {
		return args[index]
	}
	return ""
}

func GetCommandAt(index int) string {
	args := os.Args
	index += 1
	argsLen := len(args) - 1
	if len(args) > 0 && index <= argsLen {
		result := args[index]
		if strings.HasPrefix(result, FlagPrefix) {
			return ""
		}

		return result
	}

	return ""
}

// GetFlagValue Gets a value of a flag from the command line arguments
func GetFlagValue(flag string, defaultValue string) string {
	result := defaultValue
	args := os.Args
	if len(args) > 0 {
		for i, arg := range args {
			if strings.HasPrefix(arg, FlagPrefix) {
				keys := make([]string, 0)
				if strings.Index(arg, FlagValueSeparator) == -1 {
					if len(args) > i+1 {
						nextArg := GetArgumentAt(i + 1)
						if nextArg != "" && !strings.HasPrefix(nextArg, FlagPrefix) {
							keys = append(keys, arg)
							keys = append(keys, nextArg)
						}
					}
				} else {
					keys = strings.SplitAfterN(arg, FlagValueSeparator, 2)
				}
				if len(keys) == 2 {
					flagName := strings.ReplaceAll(keys[0], FlagPrefix, "")
					flagName = strings.ReplaceAll(flagName, FlagValueSeparator, "")
					if strings.ToLower(strings.TrimSpace(flagName)) == strings.ToLower(strings.TrimSpace(flag)) {
						value := keys[1]
						if value != "" {
							if value[0] == []byte("'")[0] {
								value = strings.Trim(value, "'")
							}
							if value[0] == []byte("\"")[0] {
								value = strings.Trim(value, "\"")
							}

						}
						result = value
					}
				}
			}
		}
	}
	return result
}

// GetFlagSwitch Gets a switch value of a flag from the command line arguments
func GetFlagSwitch(flag string, defaultValue bool) bool {
	result := defaultValue
	args := os.Args
	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasPrefix(arg, FlagPrefix) {
				keys := strings.SplitAfterN(arg, FlagValueSeparator, 2)
				if len(keys) == 1 {
					flagName := strings.ReplaceAll(keys[0], FlagPrefix, "")
					flagName = strings.ReplaceAll(flagName, FlagValueSeparator, "")
					if strings.ToLower(strings.TrimSpace(flagName)) == strings.ToLower(strings.TrimSpace(flag)) {
						result = true
					}
				}
			}
		}
	}
	return result
}

// GetFlagArrayValue Gets a array of values of a flag from the command line arguments
func GetFlagArrayValue(flag string) []string {
	var result []string
	args := os.Args
	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasPrefix(arg, FlagPrefix) {
				keys := strings.SplitAfterN(arg, FlagValueSeparator, 2)
				if len(keys) == 2 {
					flagName := strings.ReplaceAll(keys[0], FlagPrefix, "")
					flagName = strings.ReplaceAll(flagName, FlagValueSeparator, "")
					flagName = strings.Trim(flagName, "\"")
					flagName = strings.Trim(flagName, "'")
					if strings.ToLower(strings.TrimSpace(flagName)) == strings.ToLower(strings.TrimSpace(flag)) {
						value := keys[1]
						if value != "" {
							if value[0] == []byte("'")[0] {
								value = strings.Trim(value, "'")
							}
							if value[0] == []byte("\"")[0] {
								value = strings.Trim(value, "\"")
							}

						}
						result = append(result, value)
					}
				}
			}
		}
	}
	return result
}

// MapFlagValue Maps Keypair values from the flag
func MapFlagValue(flag string) (string, string) {
	keys := strings.Split(flag, ":")
	if len(keys) == 2 {
		return keys[0], keys[1]
	}
	return "", ""

}
