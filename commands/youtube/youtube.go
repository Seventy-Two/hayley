package youtube

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/seventy-two/Cara/web"
)

type Service struct {
	SearchURL string
	StatsURL  string
	APIKey    string
}

var serviceConfig *Service

type youtubeSearch struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	Nextpagetoken string `json:"nextPageToken"`
	Pageinfo      struct {
		Totalresults   int `json:"totalResults"`
		Resultsperpage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			Videoid string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Publishedat time.Time `json:"publishedAt"`
			Channelid   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
			Channeltitle         string `json:"channelTitle"`
			Livebroadcastcontent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

type youtubeVideo struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	Pageinfo struct {
		Totalresults   int `json:"totalResults"`
		Resultsperpage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			Publishedat time.Time `json:"publishedAt"`
			Channelid   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			Channeltitle         string   `json:"channelTitle"`
			Tags                 []string `json:"tags"`
			Categoryid           string   `json:"categoryId"`
			Livebroadcastcontent string   `json:"liveBroadcastContent"`
			Localized            struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
		} `json:"snippet"`
		Contentdetails struct {
			Duration          string `json:"duration"`
			Dimension         string `json:"dimension"`
			Definition        string `json:"definition"`
			Caption           string `json:"caption"`
			Licensedcontent   bool   `json:"licensedContent"`
			Regionrestriction struct {
				Blocked []string `json:"blocked"`
			} `json:"regionRestriction"`
		} `json:"contentDetails"`
		Statistics struct {
			Viewcount     string `json:"viewCount"`
			Likecount     string `json:"likeCount"`
			Dislikecount  string `json:"dislikeCount"`
			Favoritecount string `json:"favoriteCount"`
			Commentcount  string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
}

func parseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

func youtube(matches []string) (msg string, err error) {
	req := strings.Join(matches, " ")
	search := &youtubeSearch{}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.SearchURL, url.QueryEscape(req), serviceConfig.APIKey), search)
	if err != nil {
		return fmt.Sprintf("YouTube | Could not find video for: %s", req), nil
	}

	if search.Pageinfo.Totalresults == 0 {
		return fmt.Sprintf("YouTube | Could not find video for: %s", req), nil
	}

	id := search.Items[0].ID.Videoid

	video := &youtubeVideo{}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.StatsURL, id, serviceConfig.APIKey), video)
	if err != nil {
		return fmt.Sprintf("YouTube | Could not find video for: %s", req), nil
	}

	reg := regexp.MustCompile("\\s+")
	title := video.Items[0].Snippet.Title
	title = reg.ReplaceAllString(title, " ") // Strip excessive spaces

	duration := parseDuration(video.Items[0].Contentdetails.Duration)

	output := fmt.Sprintf("https://youtu.be/%s , %s", id, duration)

	return output, nil
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
	case "!yt":
		str, err := youtube(matches[1:])
		if err != nil {
			str = fmt.Sprintf("an error occured (%s)", err)
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
