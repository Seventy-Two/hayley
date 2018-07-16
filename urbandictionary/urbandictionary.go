package urbandictionary

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Seventy-Two/Cara/web"
)

const (
	urbanURL = "http://api.urbandictionary.com/v0/define?term=%s"
)

type DefinitionResults struct {
	Tags       []string `json:"tags"`
	ResultType string   `json:"result_type"`
	List       []struct {
		Defid       int    `json:"defid"`
		Word        string `json:"word"`
		Author      string `json:"author"`
		Permalink   string `json:"permalink"`
		Definition  string `json:"definition"`
		Example     string `json:"example"`
		ThumbsUp    int    `json:"thumbs_up"`
		ThumbsDown  int    `json:"thumbs_down"`
		CurrentVote string `json:"current_vote"`
	} `json:"list"`
	Sounds []interface{} `json:"sounds"`
}

func Urban(matches []string) (msg string, err error) {
	query := strings.Join(matches, " ")

	results := &DefinitionResults{}
	err = web.GetJSON(fmt.Sprintf(urbanURL, url.QueryEscape(query)), results)
	if err != nil {
		return fmt.Sprintf("%s | (No definition found)", query), nil
	}
	if results.ResultType == "no_results" {
		return fmt.Sprintf("%s | (No definition found)", query), nil
	}

	word := results.List[0].Word
	definition := results.List[0].Definition

	reg := regexp.MustCompile("\\s+")
	definition = reg.ReplaceAllString(definition, " ") // Strip tabs and newlines

	if len(definition) > 480 {
		definition = fmt.Sprintf("%s...", definition[0:480])
	}

	output := fmt.Sprintf("\n%s\n%s", strings.ToTitle(word), definition)

	return output, nil
}
