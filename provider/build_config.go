package provider

import (
	"encoding/json"
	"strings"
)

type BuildConfig struct {
	JSON []byte
	Site string `json:"site"`
}

type BuildConfigWarning struct {
	Message string
	Line    *int
	Char    *int
}

func NewBuildConfig(buildConfigJSON []byte) (*BuildConfig, []BuildConfigWarning) {
	warnings := make([]BuildConfigWarning, 0)
	bc := NewDefaultBuildConfig()
	if buildConfigJSON == nil {
		warnings = append(warnings, BuildConfigWarning{Message: "missing config file"})
		return bc, warnings
	}

	bc.JSON = buildConfigJSON

	err := json.Unmarshal(buildConfigJSON, bc)
	if err != nil {
		line, char := jsonDecodeError(string(buildConfigJSON), err)
		warnings = append(warnings, BuildConfigWarning{
			Message: err.Error(),
			Line:    &line,
			Char:    &char,
		})
		return bc, warnings
	}

	return bc, nil
}

func NewDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		JSON: nil,
		Site: "/",
	}
}

func jsonDecodeError(js string, err error) (int, int) {
	syntax := err.(*json.SyntaxError)

	/*start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}*/

	start, _ := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	return line + 1, pos

	//// a more detailed error message:
	//fmt.Printf("Error in line %d: %s \n", line, err)
	//fmt.Printf("%s\n%s^", js[start:end], strings.Repeat(" ", pos))
}
