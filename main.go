package main

import (
	"os"

	"github.com/seventy-two/Hayley/commands/dictionary"
	"github.com/seventy-two/Hayley/commands/stocks"
	"github.com/seventy-two/Hayley/commands/weather"
	"github.com/seventy-two/Hayley/commands/youtube"
	"github.com/seventy-two/Hayley/service"

	"github.com/seventy-two/Hayley/commands/dota"

	cli "github.com/jawher/mow.cli"
)

type appLink struct {
	name string
	url  string
}

type serviceConfig struct {
	discordAPI    *service.Service
	dotaAPI       *dota.Service
	weatherAPI    *weather.Service
	divegrassAPI  *service.Service
	dictionaryAPI *dictionary.Service
	movieAPI      *service.Service
	mathAPI       *service.Service
	tvAPI         *service.Service
	urbanAPI      *service.Service
	youtubeAPI    *youtube.Service
	nflAPI        *service.Service
	stocksAPI     *stocks.Service
}

var appMeta = struct {
	name        string
	description string
	discord     string
	maintainers string
	links       []appLink
}{
	name:        "Hayley",
	description: "Discord assistant with various commands",
	discord:     "https://discord.gg/F2cD4cN",
	maintainers: "github.com/seventy-two",
	links: []appLink{
		{name: "vcs", url: "https://github.com/seventy-two/hayley"},
	},
}

func main() {

	app := cli.App(appMeta.name, appMeta.description)

	Services := &serviceConfig{
		discordAPI: &service.Service{
			APIKey: *app.String(cli.StringOpt{
				Name:   "DiscordAPIKey",
				Value:  "",
				EnvVar: "DISCORD_API_KEY",
			}),
		},
		dotaAPI: &dota.Service{
			APIKey: *app.String(cli.StringOpt{
				Name:   "DotaAPIKey",
				Value:  "",
				EnvVar: "DOTA_API_KEY",
			}),
			DotaLeagueURL: *app.String(cli.StringOpt{
				Name:   "DotaLeagueURL",
				Value:  "http://api.steampowered.com/IDOTA2Match_570/GetLiveLeagueGames/v1/?key=%s",
				EnvVar: "DOTA_LEAGUE_URL",
			}),
			DotaListingURL: *app.String(cli.StringOpt{
				Name:   "DotaListingURL",
				Value:  "http://api.steampowered.com/IDOTA2Match_570/GetLeagueListing/v1/?key=%s",
				EnvVar: "DOTA_LISTING_URL",
			}),
			DotaMatchURL: *app.String(cli.StringOpt{
				Name:   "DotaMatchURL",
				Value:  "http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v1/?key=%s",
				EnvVar: "DOTA_MATCH_URL",
			}),
			DotaHeroesURL: *app.String(cli.StringOpt{
				Name:   "DotaMatchURL",
				Value:  "http://api.steampowered.com/IEconDOTA2_570/GetHeroes/v1/?language=en_gb&key=%s",
				EnvVar: "DOTA_MATCH_URL",
			}),
		},
		weatherAPI: &weather.Service{
			GeoCodeURL: *app.String(cli.StringOpt{
				Name:   "GeoCodeURL",
				Value:  "https://maps.googleapis.com/maps/api/geocode/json?address=%s&region=UK&key=%s",
				EnvVar: "GEOCODE_URL",
			}),
			GeoCodeAPIKey: *app.String(cli.StringOpt{
				Name:   "GeoCodeAPIKey",
				Value:  "",
				EnvVar: "GEOCODE_API_KEY",
			}),
			DarkSkyURL: *app.String(cli.StringOpt{
				Name:   "WeatherURL",
				Value:  "https://api.forecast.io/forecast/%s/%s?units=auto&exclude=minutely,hourly,alerts",
				EnvVar: "WEATHER_URL",
			}),
			DarkSkyAPIKey: *app.String(cli.StringOpt{
				Name:   "WeatherAPIKey",
				Value:  "",
				EnvVar: "WEATHER_API_KEY",
			}),
		},
		divegrassAPI: &service.Service{
			TargetURL: *app.String(cli.StringOpt{
				Name:   "DivegrassURL",
				Value:  "http://api.football-data.org/v1/competitions/467/fixtures?timeFrame=%s",
				EnvVar: "DIVEGRASS_URL",
			}),
			APIKey: *app.String(cli.StringOpt{
				Name:   "DivegrassAPIKey",
				Value:  "",
				EnvVar: "DIVEGRASS_API_KEY",
			}),
		},
		dictionaryAPI: &dictionary.Service{
			WordnikURL: *app.String(cli.StringOpt{
				Name:   "WordnikURL",
				Value:  "http://api.wordnik.com/v4/word.json/%s/definitions?limit=3&includeRelated=true&sourceDictionaries=all&useCanonical=true&includeTags=false&api_key=%s",
				EnvVar: "WORDNIK_URL",
			}),
			WOTDURL: *app.String(cli.StringOpt{
				Name:   "WOTDURL",
				Value:  "http://api.wordnik.com:80/v4/words.json/wordOfTheDay?api_key=%s",
				EnvVar: "WOTD_URL",
			}),
			WordnikAPIKey: *app.String(cli.StringOpt{
				Name:   "WordnikAPIKey",
				Value:  "",
				EnvVar: "WORDNIK_API_KEY",
			}),
		},
		movieAPI: &service.Service{
			TargetURL: *app.String(cli.StringOpt{
				Name:   "MovieURL",
				Value:  "http://www.omdbapi.com/?t=%s&y=&plot=short&r=json&tomatoes=true&apikey=%s",
				EnvVar: "MOVIE_URL",
			}),
			APIKey: *app.String(cli.StringOpt{
				Name:   "MovieAPIKey",
				Value:  "",
				EnvVar: "MOVIE_API_KEY",
			}),
		},
		mathAPI: &service.Service{
			TargetURL: *app.String(cli.StringOpt{
				Name:   "mathURL",
				Value:  "http://api.wolframalpha.com/v2/query?appid=%s&input=%s",
				EnvVar: "MATH_URL",
			}),
			APIKey: *app.String(cli.StringOpt{
				Name:   "mathAPIKey",
				Value:  "",
				EnvVar: "MATH_API_KEY",
			}),
		},
		tvAPI: &service.Service{
			TargetURL: *app.String(cli.StringOpt{
				Name:   "TvURL",
				Value:  "http://api.tvmaze.com/singlesearch/shows?q=%s",
				EnvVar: "TV_URL",
			}),
		},
		youtubeAPI: &youtube.Service{
			SearchURL: *app.String(cli.StringOpt{
				Name:   "YoutubeSearchURL",
				Value:  "https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&q=%s&key=%s",
				EnvVar: "YOUTUBE_SEARCH_URL",
			}),
			StatsURL: *app.String(cli.StringOpt{
				Name:   "YoutubeStatsURL",
				Value:  "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&id=%s&key=%s",
				EnvVar: "YOUTUBE_STATS_URL",
			}),
			APIKey: *app.String(cli.StringOpt{
				Name:   "YoutubeAPIKey",
				Value:  "",
				EnvVar: "YOUTUBE_API_KEY",
			}),
		},
		nflAPI: &service.Service{
			TargetURL: *app.String(cli.StringOpt{
				Name:   "NFLURL",
				Value:  "http://www.nfl.com/liveupdate/scores/scores.json",
				EnvVar: "NFL_API_URL",
			}),
		},
		stocksAPI: &stocks.Service{
			QuoteURL: *app.String(cli.StringOpt{
				Name:   "StocksQuoteURL",
				Value:  "https://api.iextrading.com/1.0/stock/%s/quote",
				EnvVar: "STOCKS_QUOTE_URL",
			}),
			LookupURL: *app.String(cli.StringOpt{
				Name:   "StocksLookupURL",
				Value:  "http://autoc.finance.yahoo.com/autoc?query=%s&region=EU&lang=en-GB",
				EnvVar: "STOCKS_LOOKUP_URL",
			}),
		},
	}

	app.Action = func() {
		start(app, Services)
	}

	app.Run(os.Args)
}
