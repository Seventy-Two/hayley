package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/Seventy-Two/Cara/web"
	"github.com/ryanuber/columnize"
)

func Forecast(matches []string) (msg string, err error) {
	if len(matches) < 1 {
		return "Fuck off", nil
	}

	location := strings.Title(strings.Join(matches, " "))
	coords := getCoords(location)
	if coords == "" {
		return fmt.Sprintf("Could not find %s", location), nil
	}

	data := &forecast{}
	err = web.GetJSON(fmt.Sprintf(DarkSkyURL, coords), data)
	if err != nil {
		return fmt.Sprintf("Could not get weather for: %s", location), nil
	}

	units := "°C"
	if data.Flags.Units == "us" {
		units = "°F"
	}

	output := fmt.Sprintf("Forecast: %s (%s)\n", location, coords)
	var forecasts []string
	for i := range data.Daily.Data[0:4] {
		tm := time.Unix(data.Daily.Data[i].Time, 0)
		loc, _ := time.LoadLocation(data.Timezone)
		day := tm.In(loc).Weekday()
		forecasts = append(forecasts, fmt.Sprintf("\n%s | %s | %v%s|/|%v%s ",
			day,
			Emoji(data.Daily.Data[i].Icon),
			Round(data.Daily.Data[i].TemperatureMax),
			units,
			Round(data.Daily.Data[i].TemperatureMin),
			units,
		))
	}
	output += columnize.SimpleFormat(forecasts)
	return output, nil
}
