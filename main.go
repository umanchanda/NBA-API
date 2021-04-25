package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/umanchanda/NBA-API/playertotals"
	"github.com/umanchanda/NBA-API/teamboxscore"
	"github.com/umanchanda/NBA-API/teamtotals"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	r.HandleFunc("/searchPlayer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/search.html")
	})
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
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}
	fmt.Println("listening on :8000")
	log.Fatal(http.ListenAndServe(port, r))
}
