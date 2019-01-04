package stocks

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seventy-two/Cara/web"
)

type Service struct {
	QuoteURL  string
	LookupURL string
	APIKey    string
}

var serviceConfig *Service

func getStock(matches []string) (msg string, err error) {
	if len(matches) == 0 {
		return "No search terms", nil
	}
	lookup := &Lookup{}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.LookupURL, strings.Join(matches, "+")), lookup)
	if err != nil {
		return fmt.Sprintf("There was a problem with your request. %s", err), nil
	}
	if len(lookup.ResultSet.Result) == 0 {
		return fmt.Sprintf("No results found."), nil
	}
	data := &IEXStocks{}

	var symbol string

	for _, res := range lookup.ResultSet.Result {
		if !strings.Contains(res.Symbol, ".") {
			symbol = res.Symbol
			break
		}
	}

	if symbol == "" {
		symbol = strings.Split(lookup.ResultSet.Result[0].Symbol, ".")[0]
	}

	err = web.GetJSON(fmt.Sprintf(serviceConfig.QuoteURL, symbol), data)
	if err != nil {
		return fmt.Sprintf("No data for stock symbol %s", symbol), nil
	}

	if data.CompanyName == "" {
		return fmt.Sprintf("No results found."), nil
	}

	change := data.LatestPrice - data.PreviousClose
	perChange := (change / data.PreviousClose) * 100

	sign := ""
	if change > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s - %s (%s) | %.2f ( %s%.2f %s%.2f%s )", data.CompanyName,
		data.PrimaryExchange,
		data.LatestSource,
		data.LatestPrice,
		sign,
		change,
		sign,
		perChange,
		"%",
	), nil
}

func RegisterService(dg *discordgo.Session, config *Service) {
	serviceConfig = config
	dg.AddHandler(invokeCommand)
}

func invokeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	matches := strings.Split(m.Content, " ")

	switch matches[0] {
	case "!quote":
		str, err := getStock(matches[1:])
		if err != nil {
			str = fmt.Sprintf("an error occured (%s)", err)
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
