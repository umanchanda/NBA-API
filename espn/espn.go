package espn

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/umanchanda/NBA-API/internal/fetch"
)

const scoreboardURL = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"

type Game struct {
	Name           string   `json:"name"`
	HomeTeam       string   `json:"home_team"`
	AwayTeam       string   `json:"away_team"`
	HomeScore      string   `json:"home_score"`
	AwayScore      string   `json:"away_score"`
	HomeLinescores []string `json:"home_linescores"`
	AwayLinescores []string `json:"away_linescores"`
	Status         string   `json:"status"`
	StatusDetail   string   `json:"status_detail"`
	Period         int      `json:"period"`
	Clock          string   `json:"clock"`
	SeriesSummary  string   `json:"series_summary,omitempty"`
	HomeWinner     bool     `json:"home_winner"`
}

type Scoreboard struct {
	Games []Game `json:"games"`
}

type espnResponse struct {
	Events []struct {
		Name         string `json:"name"`
		Competitions []struct {
			Status struct {
				DisplayClock string `json:"displayClock"`
				Period       int    `json:"period"`
				Type         struct {
					State       string `json:"state"`
					Detail      string `json:"detail"`
					Completed   bool   `json:"completed"`
				} `json:"type"`
			} `json:"status"`
			Competitors []struct {
				HomeAway   string `json:"homeAway"`
				Score      string `json:"score"`
				Winner     bool   `json:"winner"`
				Team       struct {
					DisplayName string `json:"displayName"`
				} `json:"team"`
				Linescores []struct {
					DisplayValue string `json:"displayValue"`
				} `json:"linescores"`
			} `json:"competitors"`
			Series *struct {
				Summary string `json:"summary"`
			} `json:"series"`
		} `json:"competitions"`
	} `json:"events"`
}

// FetchScoreboard returns today's NBA scoreboard from ESPN.
// date is optional in YYYYMMDD format; empty string fetches today.
func FetchScoreboard(date string) (string, error) {
	url := scoreboardURL
	if date != "" {
		url += "?dates=" + date
	}

	html, err := fetch.HTML(url)
	if err != nil {
		return "", fmt.Errorf("fetching ESPN scoreboard: %w", err)
	}

	var raw espnResponse
	if err := json.NewDecoder(bytes.NewReader(html)).Decode(&raw); err != nil {
		return "", fmt.Errorf("parsing ESPN response: %w", err)
	}

	games := make([]Game, 0, len(raw.Events))
	for _, event := range raw.Events {
		if len(event.Competitions) == 0 {
			continue
		}
		comp := event.Competitions[0]
		status := comp.Status

		var home, away struct {
			name       string
			score      string
			linescores []string
			winner     bool
		}

		for _, c := range comp.Competitors {
			ls := make([]string, len(c.Linescores))
			for i, s := range c.Linescores {
				ls[i] = s.DisplayValue
			}
			if c.HomeAway == "home" {
				home.name = c.Team.DisplayName
				home.score = c.Score
				home.linescores = ls
				home.winner = c.Winner
			} else {
				away.name = c.Team.DisplayName
				away.score = c.Score
				away.linescores = ls
				away.winner = c.Winner
			}
		}

		seriesSummary := ""
		if comp.Series != nil {
			seriesSummary = comp.Series.Summary
		}

		games = append(games, Game{
			Name:           event.Name,
			HomeTeam:       home.name,
			AwayTeam:       away.name,
			HomeScore:      home.score,
			AwayScore:      away.score,
			HomeLinescores: home.linescores,
			AwayLinescores: away.linescores,
			Status:         status.Type.State,
			StatusDetail:   status.Type.Detail,
			Period:         status.Period,
			Clock:          status.DisplayClock,
			SeriesSummary:  seriesSummary,
			HomeWinner:     home.winner,
		})
	}

	result, err := json.Marshal(Scoreboard{Games: games})
	if err != nil {
		return "", fmt.Errorf("encoding scoreboard: %w", err)
	}
	return string(result), nil
}
