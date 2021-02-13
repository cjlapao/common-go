package shared

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Flag constants
const (
	FlagPrefix         = "--"
	FlagValueSeparator = "="
)

// GetModuleArgument Gets the module argument for the software
func GetModuleArgument() string {
	args := os.Args[1:]

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return ""
	}

	return strings.ToLower(args[0])
}

// GetCommandArgument Gets the command argument for the software
func GetCommandArgument() string {
	args := os.Args[2:]

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return ""
	}

	return strings.ToLower(args[0])
}

// GetFlagValue Gets a value of a flag from the command line arguments
func GetFlagValue(flag string, defaultValue string) string {
	result := defaultValue
	args := os.Args
	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasPrefix(arg, FlagPrefix) {
				keys := strings.SplitAfterN(arg, FlagValueSeparator, 2)
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

func readKey(input chan rune) {
	reader := bufio.NewReader(os.Stdin)

	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("There was an error reading the key")
	}

	fmt.Printf("char: %v", char)
	input <- char
}

func keyHandler() {
	// input := make(chan rune, 1)
	// go readKey(input)
	// select {
	// case i := <-input:
	// 	fmt.Printf("Input: %v\n", i)
	// }

	signs := make(chan os.Signal, 1)

	signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)

	msg := make(chan string, 1)

	go func() {
		for {
			var s string
			fmt.Scan(&s)
			msg <- s
		}
	}()

	for {
		select {
		case <-signs:
			fmt.Println("shutdown is requested")
			os.Exit(0)
			break
		case s := <-msg:
			fmt.Println("echoing", s)
		}
	}

}
