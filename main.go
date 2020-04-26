package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	
	"github.com/umanchanda/NBA-API/teamboxscore"
	"github.com/umanchanda/NBA-API/teamtotals"
	"github.com/umanchanda/NBA-API/playertotals"
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

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/boxscore/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		boxScore := teamboxscore.ExtractBoxScore(month, day, year)
		fmt.Fprintf(w, boxScore)
	})
	r.HandleFunc("/boxscore/{year}/{month}/{day}/{awayteam}/{hometeam}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		hometeam := vars["hometeam"]
		awayteam := vars["awayteam"]
		gameSummary := teamtotals.ExtractGameSummary(month, day, year, awayteam, hometeam)
		fmt.Fprintf(w, gameSummary)
	})
	r.HandleFunc("/boxscore/{year}/{month}/{day}/{awayteam}/{hometeam}/player", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		month := vars["month"]
		day := vars["day"]
		year := vars["year"]
		hometeam := vars["hometeam"]
		awayteam := vars["awayteam"]
		gameSummary := playertotals.ExtractPlayerSummary(month, day, year, awayteam, hometeam)
		fmt.Fprintf(w, gameSummary)
	})
	fmt.Println("listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
