package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

const baseURL = "https://www.basketball-reference.com"

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
}

// AllTeamBoxScore is a struct that lists all the scores for a given day
type AllTeamBoxScore struct {
	BoxScores []TeamBoxScore `json:"box_scores,omitempty"`
}

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

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func getBoxScoreHTML(month string, day string, year string) []byte {
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

func getGameSummaryHTML(month string, day string, year string, homeTeam string) []byte {
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

func extractBoxScore(month string, day string, year string) string {
	var boxScoreHTML = getBoxScoreHTML(month, day, year)

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
		}
		scoresArray = append(scoresArray, score)
	}
	scores := AllTeamBoxScore{BoxScores: scoresArray}
	scoresJSON, _ := json.Marshal(scores)

	return string(scoresJSON)
}

func extractGameSummary(month string, day string, year string, awayTeam string, homeTeam string) string {
	var gameSummaryHTML = getGameSummaryHTML(month, day, year, homeTeam)
	p := bytes.NewReader(gameSummaryHTML)
	doc, _ := goquery.NewDocumentFromReader(p)

	boxScoresArray := make([]TeamTotals, 0)

	awayTeamBSBasic := doc.Find("#box-" + awayTeam + "-game-basic tfoot tr")

	awayStats := make([]string, 0)
	for i := 0; i < 19; i++ {
		awayStats = append(awayStats, awayTeamBSBasic.Find("td").Eq(i).Text())
	}

	ttAway := TeamTotals{
		Team:                 awayTeam,
		MinutesPlayed:        awayStats[0],
		FieldGoals:           awayStats[1],
		FieldGoalsAttempted:  awayStats[2],
		FieldGoalPercentage:  awayStats[3],
		ThreePoint:           awayStats[4],
		ThreePointAttempted:  awayStats[5],
		ThreePointPercentage: awayStats[6],
		FreeThrows:           awayStats[7],
		FreeThrowsAttempted:  awayStats[8],
		FreeThrowPercentage:  awayStats[9],
		OffensiveRebounds:    awayStats[10],
		DefensiveRebounds:    awayStats[11],
		TotalRebounds:        awayStats[12],
		Assists:              awayStats[13],
		Steals:               awayStats[14],
		Blocks:               awayStats[15],
		Turnovers:            awayStats[16],
		PersonalFouls:        awayStats[17],
		Points:               awayStats[18],
	}
	boxScoresArray = append(boxScoresArray, ttAway)

	homeTeamBSBasic := doc.Find("#box-" + homeTeam + "-game-basic tfoot")

	homeStats := make([]string, 0)
	for i := 0; i < 19; i++ {
		homeStats = append(homeStats, homeTeamBSBasic.Find("td").Eq(i).Text())
	}

	ttHome := TeamTotals{
		Team:                 homeTeam,
		MinutesPlayed:        homeStats[0],
		FieldGoals:           homeStats[1],
		FieldGoalsAttempted:  homeStats[2],
		FieldGoalPercentage:  homeStats[3],
		ThreePoint:           homeStats[4],
		ThreePointAttempted:  homeStats[5],
		ThreePointPercentage: homeStats[6],
		FreeThrows:           homeStats[7],
		FreeThrowsAttempted:  homeStats[8],
		FreeThrowPercentage:  homeStats[9],
		OffensiveRebounds:    homeStats[10],
		DefensiveRebounds:    homeStats[11],
		TotalRebounds:        homeStats[12],
		Assists:              homeStats[13],
		Steals:               homeStats[14],
		Blocks:               homeStats[15],
		Turnovers:            homeStats[16],
		PersonalFouls:        homeStats[17],
		Points:               homeStats[18],
	}
	boxScoresArray = append(boxScoresArray, ttHome)

	boxScoresArrayJSON, _ := json.Marshal(boxScoresArray)

	return string(boxScoresArrayJSON)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/boxscore/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		boxScore := extractBoxScore(month, day, year)
		fmt.Fprintf(w, boxScore)
	})
	r.HandleFunc("/boxscore/{year}/{month}/{day}/{awayteam}/{hometeam}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		hometeam := vars["hometeam"]
		awayteam := vars["awayteam"]
		gameSummary := extractGameSummary(month, day, year, awayteam, hometeam)
		fmt.Fprintf(w, gameSummary)
	})
	fmt.Println("listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
