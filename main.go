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

	r.HandleFunc("/scores/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/boxscores.html")
	})

	r.HandleFunc("/boxscore/{year}/{month}/{day}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boxScore, err := teamboxscore.ExtractBoxScore(vars["month"], vars["day"], vars["year"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, boxScore)
	})

	r.HandleFunc("/teamstats/{year}/{month}/{day}/{awayteam}/{hometeam}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/teamstats.html")
	})

	r.HandleFunc("/playerstats/{year}/{month}/{day}/{awayteam}/{hometeam}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/playerstats.html")
	})

	r.HandleFunc("/boxscore/{year}/{month}/{day}/{awayteam}/{hometeam}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameSummary, err := teamtotals.ExtractGameSummary(vars["month"], vars["day"], vars["year"], vars["awayteam"], vars["hometeam"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, gameSummary)
	})

	r.HandleFunc("/boxscore/{year}/{month}/{day}/{awayteam}/{hometeam}/player", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameSummary, err := playertotals.ExtractPlayerSummary(vars["month"], vars["day"], vars["year"], vars["awayteam"], vars["hometeam"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, gameSummary)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
