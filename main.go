package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
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

// BoxScore is a box score
type BoxScore struct {
	AwayTeam  string
	HomeTeam  string
	AwayScore string
	HomeScore string
	AwayQuarterbyQuarterScore
	HomeQuarterbyQuarterScore
}

// AwayQuarterbyQuarterScore is the quarter score for the away team
type AwayQuarterbyQuarterScore struct {
	FirstQuarter  string
	SecondQuarter string
	ThirdQuarter  string
	FourthQuarter string
}

// HomeQuarterbyQuarterScore is the quarter score for the home team
type HomeQuarterbyQuarterScore struct {
	FirstQuarter  string
	SecondQuarter string
	ThirdQuarter  string
	FourthQuarter string
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func readCSVFile(filename string) string {
	NBACSVFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(bufio.NewReader(NBACSVFile))
	var players []NBAPlayer
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		players = append(players, NBAPlayer{
			Name:     line[0],
			Height:   line[1],
			Weight:   line[2],
			Team:     line[3],
			Age:      line[4],
			Salary:   line[5],
			Points:   line[6],
			Blocks:   line[7],
			Steals:   line[8],
			Assists:  line[9],
			Rebounds: line[10],
			FT:       line[11],
			FTA:      line[12],
			FG3:      line[13],
			FG3A:     line[14],
			FG:       line[15],
			FGA:      line[16],
			MP:       line[17],
			G:        line[18],
			PER:      line[19],
			OWS:      line[20],
			DWS:      line[21],
			WS:       line[22],
			WS48:     line[23],
			USG:      line[24],
			BPM:      line[25],
			VORP:     line[26],
		})
	}

	playersJSON, _ := json.Marshal(players)
	return string(playersJSON)
}

func searchPlayer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/search.html")
	case "POST":
		player := r.FormValue("player")
		resp, err := http.Get("localhost:8000/players")
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		var players []NBAPlayer
		err = json.Unmarshal(body, &players)
		fmt.Fprintf(w, string(player))

	}
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

func getBoxScoreHTML() string {

	var boxScoreHTML = getHTML("2", "6", "2020")

	p := bytes.NewReader(boxScoreHTML)
	doc, _ := goquery.NewDocumentFromReader(p)

	gameSummary, err := doc.Find(".game_summaries").Html()
	if err != nil {
		fmt.Println(err)
	}

	return gameSummary
}

func main() {

	// playerData := readCSVFile("nbastats2018-2019.csv")

	// r := mux.NewRouter()
	// r.HandleFunc("/", index)
	// r.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, playerData)
	// })
	// http.ListenAndServe(":8000", r)
	fmt.Println(getBoxScoreHTML())
}
