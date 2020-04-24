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
	AwayTeam         string   `json:"away_team,omitempty"`
	HomeTeam         string   `json:"home_team,omitempty"`
	AwayQuarterScore []string `json:"away_quarter_score,omitempty"`
	HomeQuarterScore []string `json:"home_quarter_score,omitempty"`
}

// AllTeamBoxScore is a struct that lists all the scores for a given day
type AllTeamBoxScore struct {
	BoxScores []TeamBoxScore `json:"box_scores,omitempty"`
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func getHTML(month string, day string, year string) []byte {
	var url = baseURL + "/boxscores/?month=" + month + "&day=" + day + "&year=" + year
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return body
}

func extractGameSummary(month string, day string, year string) string {
	var boxScoreHTML = getHTML(month, day, year)

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

		score := TeamBoxScore{LosingTeam: losingTeam, WinningTeam: winningTeam, LosingTeamScore: losingTeamScore, WinningTeamScore: winningTeamScore, Status: status, AwayTeam: awayTeam, HomeTeam: homeTeam, AwayQuarterScore: awayScores, HomeQuarterScore: homeScores}
		scoresArray = append(scoresArray, score)
	}
	scores := AllTeamBoxScore{BoxScores: scoresArray}
	scoresJSON, _ := json.Marshal(scores)

	return string(scoresJSON)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/boxscore/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		gameSummary := extractGameSummary(month, day, year)
		fmt.Fprintf(w, gameSummary)
	})
	fmt.Println("listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
