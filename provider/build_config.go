package provider

import (
	"encoding/json"
	"strings"
)

type BuildConfigProperties struct {
	Site string `json:"site"`
}

type BuildConfig struct {
	JSON       []byte
	Warnings   []BuildConfigWarning
	Properties *BuildConfigProperties
}

type BuildConfigWarning struct {
	Message string
	Line    *int
	Char    *int
}

func NewBuildConfig(buildConfigJSON []byte) *BuildConfig {
	bc := NewDefaultBuildConfig()
	if buildConfigJSON == nil {
		bc.Warnings = append(bc.Warnings, BuildConfigWarning{Message: "missing config file"})
		return bc
	}

	bc.JSON = buildConfigJSON

	err := json.Unmarshal(buildConfigJSON, bc.Properties)
	if err != nil {
		line, char := jsonDecodeError(string(buildConfigJSON), err)
		bc.Warnings = append(bc.Warnings, BuildConfigWarning{
			Message: err.Error(),
			Line:    &line,
			Char:    &char,
		})
		return bc
	}

	return bc
}

func NewDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		JSON:     nil,
		Warnings: make([]BuildConfigWarning, 0),
		Properties: &BuildConfigProperties{
			Site: "/",
		},
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
