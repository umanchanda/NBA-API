package teamtotals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.basketball-reference.com"

// TeamTotals gives the totals in the box score for a team in a game
type TeamTotals struct {
	Team                 string `json:"team,omitempty"`
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
}

// TeamTotalsGame represents score breakdown for a game (both teams)
type TeamTotalsGame struct {
	TeamTotals []TeamTotals
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

// extractTeamTotals builds a TeamTotals from a tfoot selection for one team.
func extractTeamTotals(sel *goquery.Selection, team string) TeamTotals {
	td := func(i int) string { return sel.Find("td").Eq(i).Text() }
	return TeamTotals{
		Team:                 team,
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
	}
}

func ExtractGameSummary(month, day, year, awayTeam, homeTeam string) (string, error) {
	html, err := getGameSummaryHTML(month, day, year, homeTeam)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	boxScores := []TeamTotals{
		extractTeamTotals(doc.Find("#box-"+awayTeam+"-game-basic tfoot tr"), awayTeam),
		extractTeamTotals(doc.Find("#box-"+homeTeam+"-game-basic tfoot"), homeTeam),
	}

	result, err := json.Marshal(boxScores)
	if err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}
	return string(result), nil
}
