package fileproc

import (
	"encoding/json"
	"strings"
)

// Variable Entity
type Variable struct {
	Name  string
	Value string
}

// InlineVariable Entity
type InlineVariable struct {
	Path  string
	Value string
}

// GetVariables Extract the variables from a file content
func GetVariables(content string) []InlineVariable {
	foundVariables := make([]InlineVariable, 0)

	startPos := -1
	endPos := -1
	substr := content
	for {
		startPos = strings.Index(substr, "{{")
		if startPos == -1 {
			break
		}
		if startPos > -1 {
			endPos = strings.Index(substr, "}}")
			if endPos > -1 {
				endPos += 2
			}
		}
		if startPos > -1 && endPos > -1 {
			variable := substr[startPos:endPos]
			result := InlineVariable{
				Path: variable,
			}
			variable = strings.ReplaceAll(variable, "{{", "")
			variable = strings.ReplaceAll(variable, "}}", "")
			variable = strings.TrimSpace(variable)
			result.Value = variable
			substr = substr[endPos:]
			foundVariables = append(foundVariables, result)
			logger.Debug("Variable %v found\n", variable)
		}
	}

	return foundVariables
}

// ReplaceAll Replaces all variables in a string content
func ReplaceAll(content []byte, variables ...Variable) []byte {
	stringContent := string(content)
	inlineVars := GetVariables(stringContent)
	// Replacing the variables first
	for _, inlineVar := range inlineVars {
		keys := strings.Split(inlineVar.Value, ".")
		replaced := false
		if len(keys) == 2 {
			switch keys[0] {
			case "variable", "var", "parameter":
				for _, v := range variables {
					if v.Name == keys[1] {
						inlineVar.Value = v.Value
						stringContent = strings.ReplaceAll(stringContent, inlineVar.Path, inlineVar.Value)
						replaced = true
						logger.Debug("Replaced %v variable with %v", keys[1], inlineVar.Value)
						break
					}
				}

				if !replaced {
					stringContent = strings.ReplaceAll(stringContent, inlineVar.Path, "")
					logger.Debug("Could not find a variable value for %v replacing with empty string", keys[1])
				}
			}
		}
	}

	// Replacing the file variables second to take any of the replace variables in
	var contentInterface map[string]interface{}
	err := json.Unmarshal([]byte(stringContent), &contentInterface)
	if err == nil {
		for _, inlineVar := range inlineVars {
			keys := strings.Split(inlineVar.Value, ".")
			replaced := false
			if len(keys) >= 2 {
				switch keys[0] {
				case "file":
					valueKey := strings.Join(keys[1:], ".")
					value, err := GetString(contentInterface, valueKey)
					if err == nil {
						logger.Debug("Replaced  %v file variable with %v", valueKey, value)
						stringContent = strings.ReplaceAll(stringContent, inlineVar.Path, value)
						replaced = true
					}

					if !replaced {
						logger.Debug("Could not find a file variable value for %v replacing with empty string", keys[1])
						stringContent = strings.ReplaceAll(stringContent, inlineVar.Path, "")
					}
				}

			}
		}
	}
	return []byte(stringContent)
}
