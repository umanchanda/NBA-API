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

// NBAPlayer is an NBA Player
type NBAPlayer struct {
	Name     string
	Height   string
	Weight   string
	Team     string
	Age      string
	Salary   string
	Points   string
	Blocks   string
	Steals   string
	Assists  string
	Rebounds string
	FT       string
	FTA      string
	FG3      string
	FG3A     string
	FG       string
	FGA      string
	MP       string
	G        string
	PER      string
	OWS      string
	DWS      string
	WS       string
	WS48     string
	USG      string
	BPM      string
	VORP     string
}

// TeamBoxScore is a team box score
type TeamBoxScore struct {
	LosingTeam       string           `json:"losing_team,omitempty"`
	WinningTeam      string           `json:"winning_team,omitempty"`
	LosingTeamScore  string           `json:"losing_team_score,omitempty"`
	WinningTeamScore string           `json:"winning_team_score,omitempty"`
	Status           string           `json:"status,omitempty"`
	AwayQuarterScore AwayQuarterScore `json:"away_quarter_score,omitempty"`
	HomeQuarterScore HomeQuarterScore `json:"home_quarter_score,omitempty"`
}

// AwayQuarterScore is the quarter score for the away team
type AwayQuarterScore struct {
	AwayTeam      string `json:"away_team,omitempty"`
	FirstQuarter  string `json:"first_quarter,omitempty"`
	SecondQuarter string `json:"second_quarter,omitempty"`
	ThirdQuarter  string `json:"third_quarter,omitempty"`
	FourthQuarter string `json:"fourth_quarter,omitempty"`
}

// HomeQuarterScore is the quarter score for the home team
type HomeQuarterScore struct {
	HomeTeam      string `json:"home_team,omitempty"`
	FirstQuarter  string `json:"first_quarter,omitempty"`
	SecondQuarter string `json:"second_quarter,omitempty"`
	ThirdQuarter  string `json:"third_quarter,omitempty"`
	FourthQuarter string `json:"fourth_quarter,omitempty"`
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

	losingTeam, _ := doc.Find(".game_summary .teams tbody .loser td a").Html()
	losingTeamScore, _ := doc.Find(".game_summary .teams tbody .loser .right").Html()
	winningTeam, _ := doc.Find(".game_summary .teams tbody .winner td a").Html()
	winningTeamScore, _ := doc.Find(".game_summary table tbody .winner .right").Html()
	status, _ := doc.Find(".game_summary .teams tbody .loser .gamelink a").Html()

	table := doc.Find(".game_summary table").Eq(1)
	tbody1 := table.Find("tbody tr").Eq(0)
	tbody2 := table.Find("tbody tr").Eq(1)

	awayTeam, _ := tbody1.Find("td a").Html()
	homeTeam, _ := tbody2.Find("td a").Html()

	awayFirstQuarter, _ := tbody1.Find(".center").Eq(0).Html()
	awaySecondQuarter, _ := tbody1.Find(".center").Eq(1).Html()
	awayThirdQuarter, _ := tbody1.Find(".center").Eq(2).Html()
	awayFourthQuarter, _ := tbody1.Find(".center").Eq(3).Html()

	homeFirstQuarter, _ := tbody2.Find(".center").Eq(0).Html()
	homeSecondQuarter, _ := tbody2.Find(".center").Eq(1).Html()
	homeThirdQuarter, _ := tbody2.Find(".center").Eq(2).Html()
	homeFourthQuarter, _ := tbody2.Find(".center").Eq(3).Html()

	awayQuarterScore := AwayQuarterScore{AwayTeam: awayTeam, FirstQuarter: awayFirstQuarter, SecondQuarter: awaySecondQuarter, ThirdQuarter: awayThirdQuarter, FourthQuarter: awayFourthQuarter}
	homeQuarterScore := HomeQuarterScore{HomeTeam: homeTeam, FirstQuarter: homeFirstQuarter, SecondQuarter: homeSecondQuarter, ThirdQuarter: homeThirdQuarter, FourthQuarter: homeFourthQuarter}

	score := TeamBoxScore{LosingTeam: losingTeam, WinningTeam: winningTeam, LosingTeamScore: losingTeamScore, WinningTeamScore: winningTeamScore, Status: status, AwayQuarterScore: awayQuarterScore, HomeQuarterScore: homeQuarterScore}
	scoreJSON, err := json.Marshal(score)
	if err != nil {
		panic(err)
	}

	return string(scoreJSON)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	// r.HandleFunc("/searchPlayer", searchPlayer)
	r.HandleFunc("/boxscore/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		gameSummary := extractGameSummary(month, day, year)
		fmt.Fprintf(w, gameSummary)
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}
