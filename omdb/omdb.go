package omdb

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Seventy-Two/Cara/web"
)

const (
	URL = "http://www.omdbapi.com/?t=%s&y=&plot=short&r=json&tomatoes=true&apikey="
)

func Omdb(matches []string) (msg string, err error) {
	data := &movie{}
	toQuery := strings.Join(matches, "+")
	err = web.GetJSON(fmt.Sprintf(URL, url.QueryEscape(toQuery)), data)

	if err != nil {
		return fmt.Sprintf("There was a problem with your request."), nil
	}
	if data.Title == "" {
		return fmt.Sprintf("Not found."), nil
	}
	return fmt.Sprintf("%s (%s)\n%s iMDb: %s\n%s\n%s", data.Title, data.Year, data.Genre, data.ImdbRating, data.Plot, data.Actors), nil
}
