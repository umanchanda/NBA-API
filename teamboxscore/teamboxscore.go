package teamboxscore

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

// TeamBoxScore is a team box score
type TeamBoxScore struct {
	LosingTeam       string   `json:"losing_team,omitempty"`
	WinningTeam      string   `json:"winning_team,omitempty"`
	LosingTeamScore  string   `json:"losing_team_score,omitempty"`
	WinningTeamScore string   `json:"winning_team_score,omitempty"`
	Status           string   `json:"status,omitempty"`
	HomeTeam         string   `json:"home_team,omitempty"`
	AwayQuarterScore []string `json:"away_quarter_score,omitempty"`
	HomeQuarterScore []string `json:"home_quarter_score,omitempty"`
	ScoreBreakdown   string   `json:"score_breakdown,omitempty"`
	PlayerBreakdown  string   `json:"player_breakdown,omitempty`
}

// AllTeamBoxScore is a struct that lists all the scores for a given day
type AllTeamBoxScore struct {
	BoxScores []TeamBoxScore `json:"box_scores,omitempty"`
}

func GetBoxScoreHTML(month string, day string, year string) []byte {
	var url = baseURL + "/boxscores/?month=" + month + "&day=" + day + "&year=" + year
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

func ExtractBoxScore(month string, day string, year string) string {
	var boxScoreHTML = GetBoxScoreHTML(month, day, year)

	p := bytes.NewReader(boxScoreHTML)
	doc, _ := goquery.NewDocumentFromReader(p)

	gs := doc.Find(".game_summary")

	scoresArray := make([]TeamBoxScore, 0)

	for i := range gs.Nodes {
		game := gs.Eq(i)
		table0 := game.Find("table").Eq(0)
		table1 := game.Find("table").Eq(1)

		losingTeam, _ := table0.Find("tbody .loser td a").Html()
		losingTeamScore, _ := table0.Find("tbody .loser .right").Html()
		winningTeam, _ := table0.Find("tbody .winner td a").Html()
		winningTeamScore, _ := table0.Find("tbody .winner .right").Html()
		status, _ := doc.Find("tbody .loser .gamelink a").Html()

		awayScores := make([]string, 0)
		homeScores := make([]string, 0)

		awayTeam, _ := table1.Find("tbody tr").Eq(0).Find("td a").Html()
		homeTeam, _ := table1.Find("tbody tr").Eq(1).Find("td a").Html()

		periods := table1.Find("tbody tr").Eq(0).Find(".center")

		for i := range periods.Nodes {
			awayScore, _ := table1.Find("tbody tr").Eq(0).Find(".center").Eq(i).Html()
			homeScore, _ := table1.Find("tbody tr").Eq(1).Find(".center").Eq(i).Html()
			awayScores = append(awayScores, awayScore)
			homeScores = append(homeScores, homeScore)
		}

		awayTeamCode := teamCodes[awayTeam]
		homeTeamCode := teamCodes[homeTeam]
		scoreBreakdown := "/boxscore/" + year + "/" + month + "/" + day + "/" + awayTeamCode + "/" + homeTeamCode
		playerBreakdown := "/boxscore/" + year + "/" + month + "/" + day + "/" + awayTeamCode + "/" + homeTeamCode + "/player"

		score := TeamBoxScore{
			LosingTeam:       losingTeam,
			WinningTeam:      winningTeam,
			LosingTeamScore:  losingTeamScore,
			WinningTeamScore: winningTeamScore,
			Status:           status,
			HomeTeam:         homeTeam,
			AwayQuarterScore: awayScores,
			HomeQuarterScore: homeScores,
			ScoreBreakdown:   scoreBreakdown,
			PlayerBreakdown:  playerBreakdown,
		}
		scoresArray = append(scoresArray, score)
	}
	scores := AllTeamBoxScore{BoxScores: scoresArray}
	scoresJSON, _ := json.Marshal(scores)

	return string(scoresJSON)
}
