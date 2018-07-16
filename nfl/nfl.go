package nfl

import (
	"fmt"
	"gopkg.in/xmlpath.v2"
	"net/http"
	"strings"
	"time"
)

const (
	URL = "http://www.nfl.com/liveupdate/scorestrip/ss.xml"
)

func Nfl() (msg []string, err error) {
	doc, _ := http.Get(fmt.Sprintf(URL))
	defer doc.Body.Close()
	root, err := xmlpath.Parse(doc.Body)
	var debug bool

	debug = true

	if err != nil {
		msg = append(msg, fmt.Sprintf("Could not retrieve matches."))
		return msg, nil
	}

	todaysdate := getToday()
	d := xmlpath.MustCompile("/ss/gms/g/@d")

	dateIter := d.Iter(root)
	var i int
	var timeStr string
	var awayScoreStr string
	var homeScoreStr string
	i = 1
	for dateIter.Next() {
		timeStr = ""
		awayScoreStr = ""
		homeScoreStr = ""

		if strings.EqualFold(todaysdate, dateIter.Node().String()) || debug {
			home := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@hnn", i))
			homeScore := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@hs", i))
			away := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@vnn", i))
			awayScore := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@vs", i))
			quarter := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@q", i))
			t := xmlpath.MustCompile(fmt.Sprintf("/ss/gms/g[%d]/@t", i))

			homeStr, _ := home.String(root)
			awayStr, _ := away.String(root)
			quarterStr, _ := quarter.String(root)
			if strings.EqualFold(quarterStr, "P") {
				timeStr, _ = t.String(root)
				timeStr = timeStr + " ET"
			} else {
				homeScoreStr, _ = homeScore.String(root)
				awayScoreStr, _ = awayScore.String(root)
			}

			homeStr = getTeamColour(homeStr)
			awayStr = getTeamColour(awayStr)
			quarterStr = fixQuarter(quarterStr)

			out := fmt.Sprintf(homeStr + " | " + homeScoreStr + " | - | " + awayScoreStr + " | " + awayStr + " | [" + quarterStr + timeStr + "]")
			msg = append(msg, out)
		}
		i++
	}

	return msg, nil
}

func getToday() (date string) {
	now := time.Now()
	nowUTC := now.UTC()
	loc, _ := time.LoadLocation("America/New_York")
	jst := nowUTC.In(loc)
	return jst.Format("Mon")
}

func getTeamColour(team string) (colouredTeam string) {
	switch team {
	case "cardinals":
		return "Arizona Cardinals"
	case "falcons":
		return "Atlanta Falcons"
	case "panthers":
		return "Carolina Panthers"
	case "bears":
		return "Chicago Bears"
	case "cowboys":
		return "Dallas Cowboys"
	case "lions":
		return "Detroit Lions"
	case "packers":
		return "Green Bay Packers"
	case "vikings":
		return "Minnesota Vikings"
	case "saints":
		return "New Orleans Saints"
	case "giants":
		return "New York Giants"
	case "eagles":
		return "Philadelphia Eagles"
	case "rams":
		return "Los Angeles Rams"
	case "49ers":
		return "San Fransisco 49ers"
	case "seahawks":
		return "Seattle Seahawks"
	case "buccaneers":
		return "Tampa Bay Buccaneers"
	case "redskins":
		return "Washington Redskins"
	case "ravens":
		return "Baltimore Ravens"
	case "bills":
		return "Buffalo Bills"
	case "bengals":
		return "Cincinnati Bengals"
	case "browns":
		return "Cleveland Browns"
	case "broncos":
		return "Denver Broncos"
	case "texans":
		return "Houston Texans"
	case "colts":
		return "Indianapolis Colts"
	case "jaguars":
		return "Jacksonville Jaguars"
	case "chiefs":
		return "Kansas City Chiefs"
	case "dolphins":
		return "Miami Dolphins"
	case "patriots":
		return "New England Patriots"
	case "jets":
		return "New York Jets"
	case "raiders":
		return "Oakland Raiders"
	case "steelers":
		return "Pittsburgh Steelers"
	case "chargers":
		return "Los Angeles Chargers"
	case "titans":
		return "Tennessee Titans"
	default:
		return team
	}
}

func fixQuarter(quarter string) (prettyQuarter string) {
	switch quarter {
	case "P":
		return ""
	case "1":
		return "Q1"
	case "2":
		return "Q2"
	case "3":
		return "Q3"
	case "4":
		return "Q4"
	case "5":
		return "OT"
	case "F":
		return "Final"
	case "FO":
		return "Final (OT)"
	default:
		return quarter
	}
}
