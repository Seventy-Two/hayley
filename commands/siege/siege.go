package siege

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	AuthURL      string
	AuthUser     string
	AuthPassword string
	AuthToken    string
	AuthExpiry   *time.Time
	ProfileURL   string
	LevelURL     string
	RankURL      string
}

var serviceConfig *Service

func siege(user string) (string, error) {
	client := &http.Client{}

	if serviceConfig.AuthExpiry == nil || serviceConfig.AuthExpiry.Before(time.Now()) {
		tok, expiry, err := authenticate(client, serviceConfig.AuthURL, serviceConfig.AuthUser, serviceConfig.AuthPassword)
		if err != nil {
			log.Print(err)
			return "", nil
		}
		serviceConfig.AuthToken = tok
		serviceConfig.AuthExpiry = expiry
	}

	id, err := retrieveUserID(client, user, serviceConfig.ProfileURL, serviceConfig.AuthToken)
	if err != nil {
		log.Print(err)
		return "", nil
	}
	level, err := retrieveUserLevel(client, id, serviceConfig.LevelURL, serviceConfig.AuthToken)
	if err != nil {
		log.Print(err)
		return "", nil
	}

	rank, mmr, season, err := retrieveUserRank(client, id, serviceConfig.RankURL, serviceConfig.AuthToken)
	if err != nil {
		log.Print(err)
		return "", nil
	}

	cleanRank := convertRank(rank)
	cleanSeaons := convertSeason(season)

	out := fmt.Sprintf("%s - Level %s - Rank %s - %s MMR - %s", user, level, cleanRank, mmr, cleanSeaons)
	return out, nil
}

func convertSeason(season int) string {
	if season == 0 {
		return "Season Unknown"
	}
	year := int(season / 4)
	sub := (season % 4)
	if sub == 0 {
		sub = 4
	}
	return fmt.Sprintf("Year %d Season %d", year, sub)
}

func convertRank(rank int) string {
	switch rank {
	case 0:
		return "Copper IV"
	case 1:
		return "Copper III"
	case 2:
		return "Copper II"
	case 3:
		return "Copper I"
	case 4:
		return "Bronze IV"
	case 5:
		return "Bronze III"
	case 6:
		return "Bronze II"
	case 7:
		return "Bronze I"
	case 8:
		return "Silver IV"
	case 9:
		return "Silver III"
	case 10:
		return "Silver II"
	case 11:
		return "Silver I"
	case 12:
		return "Gold IV"
	case 13:
		return "Gold III"
	case 14:
		return "Gold II"
	case 15:
		return "Gold I"
	case 16:
		return "Platinum III"
	case 17:
		return "Platinum II"
	case 18:
		return "Platinum I"
	case 19:
		return "Diamond"
	default:
		return "Coward"
	}
}

func retrieveUserRank(client *http.Client, id, rankURL, auth string) (int, string, int, error) {
	body, err := makeUbiRequest(client, fmt.Sprintf(rankURL, id), auth)
	r := &rankResponse{}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Print(err)
		return 0, "", 0, err
	}
	for _, player := range r.Players {
		return player.Rank, strconv.Itoa(int(player.Mmr)), player.Season, nil
	}
	return 0, "", 0, nil
}

type rankResponse struct {
	Players map[string]player `json:"players"`
}

type player struct {
	BoardID             string    `json:"board_id"`
	PastSeasonsAbandons int       `json:"past_seasons_abandons"`
	UpdateTime          time.Time `json:"update_time"`
	SkillMean           float64   `json:"skill_mean"`
	Abandons            int       `json:"abandons"`
	Season              int       `json:"season"`
	Region              string    `json:"region"`
	ProfileID           string    `json:"profile_id"`
	PastSeasonsLosses   int       `json:"past_seasons_losses"`
	MaxMmr              float64   `json:"max_mmr"`
	Mmr                 float64   `json:"mmr"`
	Wins                int       `json:"wins"`
	SkillStdev          float64   `json:"skill_stdev"`
	Rank                int       `json:"rank"`
	Losses              int       `json:"losses"`
	NextRankMmr         float64   `json:"next_rank_mmr"`
	PastSeasonsWins     int       `json:"past_seasons_wins"`
	PreviousRankMmr     float64   `json:"previous_rank_mmr"`
	MaxRank             int       `json:"max_rank"`
}

func retrieveUserLevel(client *http.Client, id, levelURL, auth string) (string, error) {
	body, err := makeUbiRequest(client, fmt.Sprintf(levelURL, id), auth)
	r := &levelResponse{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return "", err
	}

	for _, prof := range r.PlayerProfiles {
		return strconv.Itoa(prof.Level), nil
	}

	return "", nil
}

type levelResponse struct {
	PlayerProfiles []struct {
		Xp                 int    `json:"xp"`
		ProfileID          string `json:"profile_id"`
		LootboxProbability int    `json:"lootbox_probability"`
		Level              int    `json:"level"`
	} `json:"player_profiles"`
}

func retrieveUserID(client *http.Client, user, profURL, auth string) (string, error) {
	body, err := makeUbiRequest(client, fmt.Sprintf(profURL, user), auth)
	r := &userResponse{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return "", err
	}
	for _, prof := range r.Profiles {
		return prof.IDOnPlatform, nil
	}
	return "", nil
}

type userResponse struct {
	Profiles []struct {
		ProfileID      string `json:"profileId"`
		UserID         string `json:"userId"`
		PlatformType   string `json:"platformType"`
		IDOnPlatform   string `json:"idOnPlatform"`
		NameOnPlatform string `json:"nameOnPlatform"`
	} `json:"profiles"`
}

func authenticate(client *http.Client, authURL, user, pass string) (string, *time.Time, error) {
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Ubi-AppId", "39baebad-39e5-4552-8c25-2c9b919064e2")

	// eUser := base64.StdEncoding.EncodeToString([]byte(user))
	// ePass := base64.StdEncoding.EncodeToString([]byte(pass))
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	r := &authResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(body, r)
	if err != nil {
		return "", nil, err
	}

	return r.Ticket, &r.Expiration, nil
}

func makeUbiRequest(client *http.Client, ubiURL, tok string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ubiURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Ubi-AppId", "39baebad-39e5-4552-8c25-2c9b919064e2")
	req.Header.Set("Authorization", "Ubi_v1 t="+tok)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type authResponse struct {
	PlatformType                  string      `json:"platformType"`
	Ticket                        string      `json:"ticket"`
	TwoFactorAuthenticationTicket interface{} `json:"twoFactorAuthenticationTicket"`
	ProfileID                     string      `json:"profileId"`
	UserID                        string      `json:"userId"`
	NameOnPlatform                string      `json:"nameOnPlatform"`
	Environment                   string      `json:"environment"`
	Expiration                    time.Time   `json:"expiration"`
	SpaceID                       string      `json:"spaceId"`
	ClientIP                      string      `json:"clientIp"`
	ClientIPCountry               string      `json:"clientIpCountry"`
	ServerTime                    time.Time   `json:"serverTime"`
	SessionID                     string      `json:"sessionId"`
	SessionKey                    string      `json:"sessionKey"`
	RememberMeTicket              interface{} `json:"rememberMeTicket"`
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
	var str string

	switch matches[0] {
	case "!r6":
		res, err := siege(matches[1])

		if err != nil {
			str = fmt.Sprintf("an error occured (%s)", err)
		} else {
			str = res
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
