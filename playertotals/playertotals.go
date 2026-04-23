package playertotals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.basketball-reference.com"

// PlayerTotals represents a player box score from a single game
type PlayerTotals struct {
	Team                 string `json:"team,omitempty"`
	Name                 string `json:"name,omitempty"`
	MinutesPlayed        string `json:"minutes_played,omitempty"`
	FieldGoals           string `json:"field_goals,omitempty"`
	FieldGoalsAttempted  string `json:"field_goals_attempted,omitempty"`
	FieldGoalPercentage  string `json:"field_goal_percentage,omitempty"`
	ThreePoint           string `json:"three_point,omitempty"`
	ThreePointAttempted  string `json:"three_point_attempted,omitempty"`
	ThreePointPercentage string `json:"three_point_percentage,omitempty"`
	FreeThrows           string `json:"free_throws,omitempty"`
	FreeThrowsAttempted  string `json:"free_throws_attempted,omitempty"`
	FreeThrowPercentage  string `json:"free_throw_percentage,omitempty"`
	OffensiveRebounds    string `json:"offensive_rebounds,omitempty"`
	DefensiveRebounds    string `json:"defensive_rebounds,omitempty"`
	TotalRebounds        string `json:"total_rebounds,omitempty"`
	Assists              string `json:"assists,omitempty"`
	Steals               string `json:"steals,omitempty"`
	Blocks               string `json:"blocks,omitempty"`
	Turnovers            string `json:"turnovers,omitempty"`
	PersonalFouls        string `json:"personal_fouls,omitempty"`
	Points               string `json:"points,omitempty"`
	PlusMinus            string `json:"plus_minus,omitempty"`
}

// PlayerTotalsTeam is the list of all players in a game for a team and their score breakdown
type PlayerTotalsTeam struct {
	Starters []PlayerTotals `json:"starters,omitempty"`
	Reserves []PlayerTotals `json:"reserves,omitempty"`
}

// PlayerTotalsGame is the list of all players in a game for both teams
type PlayerTotalsGame struct {
	Players []PlayerTotalsTeam `json:"players,omitempty"`
}

func getGameSummaryHTML(month, day, year, homeTeam string) ([]byte, error) {
	url := baseURL + "/boxscores/" + year + month + day + "0" + homeTeam + ".html"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching boxscore page: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return body, nil
}

// extractPlayerRow builds a PlayerTotals from a single table row selection.
func extractPlayerRow(row *goquery.Selection, team string) PlayerTotals {
	td := func(i int) string { return row.Find("td").Eq(i).Text() }
	return PlayerTotals{
		Team:                 team,
		Name:                 row.Find("a").Text(),
		MinutesPlayed:        td(0),
		FieldGoals:           td(1),
		FieldGoalsAttempted:  td(2),
		FieldGoalPercentage:  td(3),
		ThreePoint:           td(4),
		ThreePointAttempted:  td(5),
		ThreePointPercentage: td(6),
		FreeThrows:           td(7),
		FreeThrowsAttempted:  td(8),
		FreeThrowPercentage:  td(9),
		OffensiveRebounds:    td(10),
		DefensiveRebounds:    td(11),
		TotalRebounds:        td(12),
		Assists:              td(13),
		Steals:               td(14),
		Blocks:               td(15),
		Turnovers:            td(16),
		PersonalFouls:        td(17),
		Points:               td(18),
		PlusMinus:            td(19),
	}
}

// extractTeamPlayers collects starters (rows 0-4) and reserves (rows 6-15) for a team.
func extractTeamPlayers(doc *goquery.Document, team string) PlayerTotalsTeam {
	rows := doc.Find("#box-" + team + "-game-basic tbody tr")

	starters := make([]PlayerTotals, 0, 5)
	for i := 0; i < 5; i++ {
		starters = append(starters, extractPlayerRow(rows.Eq(i), team))
	}

	reserves := make([]PlayerTotals, 0, 10)
	for i := 6; i < 16; i++ {
		reserves = append(reserves, extractPlayerRow(rows.Eq(i), team))
	}

	return PlayerTotalsTeam{Starters: starters, Reserves: reserves}
}

func ExtractPlayerSummary(month, day, year, awayTeam, homeTeam string) (string, error) {
	html, err := getGameSummaryHTML(month, day, year, homeTeam)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	allPlayers := []PlayerTotalsTeam{
		extractTeamPlayers(doc, awayTeam),
		extractTeamPlayers(doc, homeTeam),
	}

	result, err := json.Marshal(allPlayers)
	if err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}
	return string(result), nil
}
