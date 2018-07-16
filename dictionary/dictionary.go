package dictionary

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Seventy-Two/Cara/web"
)

const (
	nikURL  = "http://api.wordnik.com/v4/word.json/%s/definitions?limit=3&includeRelated=true&sourceDictionaries=all&useCanonical=true&includeTags=false&api_key="
	wotdURL = "http://api.wordnik.com:80/v4/words.json/wordOfTheDay?api_key="
)

func Dict(matches []string) (msg string, err error) {
	text := url.QueryEscape(strings.Join(matches, "+"))
	var data []Wordnik
	var result []string

	err = web.GetJSON(fmt.Sprintf(nikURL, text), &data)
	if err != nil {
		return fmt.Sprintf("There was a problem with your request."), err
	}
	if len(data) == 0 {
		return fmt.Sprintf("Word/phrase not found."), nil
	}
	cap := len(data) // never >3 because limit=3 in URL
	for i := 0; i < cap; i++ {
		result = append(result, fmt.Sprintf("%s - %s\n%s", data[i].Word, data[i].PartOfSpeech, data[i].Text))
	}
	out := ""
	for _, res := range result {
		out += res
		out += "\n"
	}

	return out, nil
}

func Wotd() (msg string, err error) {
	data := &wotd{}
	err = web.GetJSON(wotdURL, data)
	if err != nil {
		return fmt.Sprintf("There was a problem with your request."), nil
	}
	return fmt.Sprintf("Word of the day: %s\n%s - %s", data.Word, data.Note, data.Definitions[0].Text), nil // I really hate doing [0] but we only want one definition. I hate comments that cause horizontal scroll also.
}
