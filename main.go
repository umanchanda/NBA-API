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

// NBAPlayer contains fields for various nba stats
type NBAPlayer struct {
	Name     string `json:"name,omitempty"`
	Height   string `json:"height,omitempty"`
	Weight   string `json:"weight,omitempty"`
	Team     string `json:"team,omitempty"`
	Age      string `json:"age,omitempty"`
	Salary   string `json:"salary,omitempty"`
	Points   string `json:"points,omitempty"`
	Blocks   string `json:"blocks,omitempty"`
	Steals   string `json:"steals,omitempty"`
	Assists  string `json:"assists,omitempty"`
	Rebounds string `json:"rebounds,omitempty"`
	FT       string `json:"ft,omitempty"`
	FTA      string `json:"fta,omitempty"`
	FG3      string `json:"fg3,omitempty"`
	FG3A     string `json:"fg3a,omitempty"`
	FG       string `json:"fg,omitempty"`
	FGA      string `json:"fga,omitempty"`
	MP       string `json:"mp,omitempty"`
	G        string `json:"g,omitempty"`
	PER      string `json:"per,omitempty"`
	OWS      string `json:"ows,omitempty"`
	DWS      string `json:"dws,omitempty"`
	WS       string `json:"ws,omitempty"`
	WS48     string `json:"ws48,omitempty"`
	USG      string `json:"usg,omitempty"`
	BPM      string `json:"bpm,omitempty"`
	VORP     string `json:"vorp,omitempty"`
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
