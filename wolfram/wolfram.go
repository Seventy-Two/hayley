package wolfram

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/xmlpath.v2"
)

const (
	wolframURL = "http://api.wolframalpha.com/v2/query?appid=WJKLJW-UJPL8RUJLH&input=%s"
)

func extractURL(text string) string {
	extractedURL := ""
	for _, value := range strings.Split(text, " ") {
		parsedURL, err := url.Parse(value)
		if err != nil {
			continue
		}
		if strings.HasPrefix(parsedURL.Scheme, "http") {
			extractedURL = parsedURL.String()
			break
		}
	}
	return extractedURL
}

func Wolfram(matches []string) (msg string, err error) {
	doc, _ := http.Get(fmt.Sprintf(wolframURL, url.QueryEscape(strings.Join(matches, " "))))
	defer doc.Body.Close()
	root, err := xmlpath.Parse(doc.Body)

	if err != nil {
		return "Wolfram | Stephen Wolfram doesn't know the answer to this", nil
	}

	success := xmlpath.MustCompile("//queryresult/@success")
	input := xmlpath.MustCompile("//pod[@position='100']//plaintext[1]")
	output := xmlpath.MustCompile("//pod[@position='200']/subpod[1]/plaintext[1]")

	suc, _ := success.String(root)

	if suc != "true" {
		return fmt.Sprintf("Wolfram | Stephen Wolfram doesn't know the answer to this"), nil
	}

	in, _ := input.String(root)
	out, _ := output.String(root)

	in = strings.Replace(in, `\:`, `\u`, -1)
	out = strings.Replace(out, `\:`, `\u`, -1)

	reg := regexp.MustCompile("\\s+")
	in = reg.ReplaceAllString(in, " ")
	out = reg.ReplaceAllString(out, " ")

	in, _ = strconv.Unquote(`"` + in + `"`)
	out, _ = strconv.Unquote(`"` + out + `"`)

	return fmt.Sprintf("Wolfram\n%s >>> %s", in, out), nil
}
