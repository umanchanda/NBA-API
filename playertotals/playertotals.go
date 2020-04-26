package playertotals

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.basketball-reference.com"

var teamCodes = map[string]string{
	"Atlanta":       "ATL",
	"Boston":        "BOS",
	"Brooklyn":      "BKN",
	"Charlotte":     "CHO",
	"Chicago":       "CHI",
	"Cleveland":     "CLE",
	"Dallas":        "DAL",
	"Denver":        "DEN",
	"Detroit":       "DET",
	"Golden State":  "GSW",
	"Houston":       "HOU",
	"Indiana":       "IND",
	"LA Lakers":     "LAL",
	"LA Clippers":   "LAC",
	"Memphis":       "MEM",
	"Miami":         "MIA",
	"Milwaukee":     "MIL",
	"Minnesota":     "MIN",
	"New Orleans":   "NOP",
	"New York":      "NYK",
	"Oklahoma City": "OKC",
	"Orlando":       "ORL",
	"Philadelphia":  "PHI",
	"Phoenix":       "PHO",
	"Portland":      "POR",
	"Sacramento":    "SAC",
	"San Antonio":   "SAS",
	"Toronto":       "TOR",
	"Utah":          "UTA",
	"Washington":    "WAS",
}

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

// PlayerTotalsGame is the list of all players in a game for both both teams
type PlayerTotalsGame struct {
	Players []PlayerTotalsTeam `json:"players,omitempty"`
}

func GetGameSummaryHTML(month string, day string, year string, homeTeam string) []byte {
	var url = baseURL + "/boxscores/" + year + month + day + "0" + homeTeam + ".html"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}

func ExtractPlayerSummary(month string, day string, year string, awayTeam string, homeTeam string) string {
	var gameSummaryHTML = GetGameSummaryHTML(month, day, year, homeTeam)
	p := bytes.NewReader(gameSummaryHTML)
	doc, _ := goquery.NewDocumentFromReader(p)

	allPlayers := make([]PlayerTotalsTeam, 0)

	awayStarters := make([]PlayerTotals, 0)
	awayReserves := make([]PlayerTotals, 0)
	for i := 0; i < 5; i++ {
		player := doc.Find("#box-" + awayTeam + "-game-basic tbody tr").Eq(i)
		awayStarters = append(awayStarters, PlayerTotals{
			Team:                 awayTeam,
			Name:                 player.Find("a").Text(),
			MinutesPlayed:        player.Find("td").Eq(0).Text(),
			FieldGoals:           player.Find("td").Eq(1).Text(),
			FieldGoalsAttempted:  player.Find("td").Eq(2).Text(),
			FieldGoalPercentage:  player.Find("td").Eq(3).Text(),
			ThreePoint:           player.Find("td").Eq(4).Text(),
			ThreePointAttempted:  player.Find("td").Eq(5).Text(),
			ThreePointPercentage: player.Find("td").Eq(6).Text(),
			FreeThrows:           player.Find("td").Eq(7).Text(),
			FreeThrowsAttempted:  player.Find("td").Eq(8).Text(),
			FreeThrowPercentage:  player.Find("td").Eq(9).Text(),
			OffensiveRebounds:    player.Find("td").Eq(10).Text(),
			DefensiveRebounds:    player.Find("td").Eq(11).Text(),
			TotalRebounds:        player.Find("td").Eq(12).Text(),
			Assists:              player.Find("td").Eq(13).Text(),
			Steals:               player.Find("td").Eq(14).Text(),
			Blocks:               player.Find("td").Eq(15).Text(),
			Turnovers:            player.Find("td").Eq(16).Text(),
			PersonalFouls:        player.Find("td").Eq(17).Text(),
			Points:               player.Find("td").Eq(18).Text(),
			PlusMinus:            player.Find("td").Eq(19).Text(),
		})
	}

	for i := 6; i < 16; i++ {
		player := doc.Find("#box-" + awayTeam + "-game-basic tbody tr").Eq(i)
		awayReserves = append(awayReserves, PlayerTotals{
			Team:                 awayTeam,
			Name:                 player.Find("a").Text(),
			MinutesPlayed:        player.Find("td").Eq(0).Text(),
			FieldGoals:           player.Find("td").Eq(1).Text(),
			FieldGoalsAttempted:  player.Find("td").Eq(2).Text(),
			FieldGoalPercentage:  player.Find("td").Eq(3).Text(),
			ThreePoint:           player.Find("td").Eq(4).Text(),
			ThreePointAttempted:  player.Find("td").Eq(5).Text(),
			ThreePointPercentage: player.Find("td").Eq(6).Text(),
			FreeThrows:           player.Find("td").Eq(7).Text(),
			FreeThrowsAttempted:  player.Find("td").Eq(8).Text(),
			FreeThrowPercentage:  player.Find("td").Eq(9).Text(),
			OffensiveRebounds:    player.Find("td").Eq(10).Text(),
			DefensiveRebounds:    player.Find("td").Eq(11).Text(),
			TotalRebounds:        player.Find("td").Eq(12).Text(),
			Assists:              player.Find("td").Eq(13).Text(),
			Steals:               player.Find("td").Eq(14).Text(),
			Blocks:               player.Find("td").Eq(15).Text(),
			Turnovers:            player.Find("td").Eq(16).Text(),
			PersonalFouls:        player.Find("td").Eq(17).Text(),
			Points:               player.Find("td").Eq(18).Text(),
			PlusMinus:            player.Find("td").Eq(19).Text(),
		})
	}

	away := PlayerTotalsTeam{
		Starters: awayStarters,
		Reserves: awayReserves,
	}

	allPlayers = append(allPlayers, away)

	homeStarters := make([]PlayerTotals, 0)
	homeReserves := make([]PlayerTotals, 0)
	for i := 0; i < 5; i++ {
		player := doc.Find("#box-" + homeTeam + "-game-basic tbody tr").Eq(i)
		homeStarters = append(homeStarters, PlayerTotals{
			Team:                 homeTeam,
			Name:                 player.Find("a").Text(),
			MinutesPlayed:        player.Find("td").Eq(0).Text(),
			FieldGoals:           player.Find("td").Eq(1).Text(),
			FieldGoalsAttempted:  player.Find("td").Eq(2).Text(),
			FieldGoalPercentage:  player.Find("td").Eq(3).Text(),
			ThreePoint:           player.Find("td").Eq(4).Text(),
			ThreePointAttempted:  player.Find("td").Eq(5).Text(),
			ThreePointPercentage: player.Find("td").Eq(6).Text(),
			FreeThrows:           player.Find("td").Eq(7).Text(),
			FreeThrowsAttempted:  player.Find("td").Eq(8).Text(),
			FreeThrowPercentage:  player.Find("td").Eq(9).Text(),
			OffensiveRebounds:    player.Find("td").Eq(10).Text(),
			DefensiveRebounds:    player.Find("td").Eq(11).Text(),
			TotalRebounds:        player.Find("td").Eq(12).Text(),
			Assists:              player.Find("td").Eq(13).Text(),
			Steals:               player.Find("td").Eq(14).Text(),
			Blocks:               player.Find("td").Eq(15).Text(),
			Turnovers:            player.Find("td").Eq(16).Text(),
			PersonalFouls:        player.Find("td").Eq(17).Text(),
			Points:               player.Find("td").Eq(18).Text(),
			PlusMinus:            player.Find("td").Eq(19).Text(),
		})
	}

	for i := 6; i < 16; i++ {
		player := doc.Find("#box-" + homeTeam + "-game-basic tbody tr").Eq(i)
		homeReserves = append(homeReserves, PlayerTotals{
			Team:                 homeTeam,
			Name:                 player.Find("a").Text(),
			MinutesPlayed:        player.Find("td").Eq(0).Text(),
			FieldGoals:           player.Find("td").Eq(1).Text(),
			FieldGoalsAttempted:  player.Find("td").Eq(2).Text(),
			FieldGoalPercentage:  player.Find("td").Eq(3).Text(),
			ThreePoint:           player.Find("td").Eq(4).Text(),
			ThreePointAttempted:  player.Find("td").Eq(5).Text(),
			ThreePointPercentage: player.Find("td").Eq(6).Text(),
			FreeThrows:           player.Find("td").Eq(7).Text(),
			FreeThrowsAttempted:  player.Find("td").Eq(8).Text(),
			FreeThrowPercentage:  player.Find("td").Eq(9).Text(),
			OffensiveRebounds:    player.Find("td").Eq(10).Text(),
			DefensiveRebounds:    player.Find("td").Eq(11).Text(),
			TotalRebounds:        player.Find("td").Eq(12).Text(),
			Assists:              player.Find("td").Eq(13).Text(),
			Steals:               player.Find("td").Eq(14).Text(),
			Blocks:               player.Find("td").Eq(15).Text(),
			Turnovers:            player.Find("td").Eq(16).Text(),
			PersonalFouls:        player.Find("td").Eq(17).Text(),
			Points:               player.Find("td").Eq(18).Text(),
			PlusMinus:            player.Find("td").Eq(19).Text(),
		})
	}

	home := PlayerTotalsTeam{
		Starters: homeStarters,
		Reserves: homeReserves,
	}

	allPlayers = append(allPlayers, home)
	allPlayersJSON, _ := json.Marshal(allPlayers)

	return string(allPlayersJSON)
}