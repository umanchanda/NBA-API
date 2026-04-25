package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/umanchanda/NBA-API/espn"
	"github.com/umanchanda/NBA-API/playertotals"
	"github.com/umanchanda/NBA-API/teamboxscore"
	"github.com/umanchanda/NBA-API/teamtotals"

	"database/sql"
)

func dbConn() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	}
	return sql.Open("postgres", connStr)
}

type PlayerStat struct {
	Season     string `json:"season"`
	SeasonType string `json:"season_type"`
	Name       string `json:"name"`
	Team       string `json:"team"`
	Pos        string `json:"pos"`
	Age        string `json:"age"`
	G          string `json:"g"`
	GS         string `json:"gs"`
	MP         string `json:"mp"`
	FG         string `json:"fg"`
	FGA        string `json:"fga"`
	FGPct      string `json:"fg_pct"`
	FG3        string `json:"fg3"`
	FG3A       string `json:"fg3a"`
	FG3Pct     string `json:"fg3_pct"`
	FT         string `json:"ft"`
	FTA        string `json:"fta"`
	FTPct      string `json:"ft_pct"`
	ORB        string `json:"orb"`
	DRB        string `json:"drb"`
	TRB        string `json:"trb"`
	AST        string `json:"ast"`
	STL        string `json:"stl"`
	BLK        string `json:"blk"`
	TOV        string `json:"tov"`
	PF         string `json:"pf"`
	PTS        string `json:"pts"`
}

func main() {
	db, err := dbConn()
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	r.HandleFunc("/searchPlayer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/search.html")
	})

	r.HandleFunc("/api/player", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		season := r.URL.Query().Get("season")
		seasonType := r.URL.Query().Get("season_type")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		query := `SELECT season, season_type, name, team, pos, age, g, gs, mp,
			fg, fga, fg_pct, fg3, fg3a, fg3_pct,
			ft, fta, ft_pct, orb, drb, trb,
			ast, stl, blk, tov, pf, pts
			FROM playerstats
			WHERE LOWER(name) LIKE LOWER($1)`
		args := []interface{}{"%" + name + "%"}

		if season != "" {
			args = append(args, season)
			query += fmt.Sprintf(" AND season = $%d", len(args))
		}
		if seasonType != "" {
			args = append(args, seasonType)
			query += fmt.Sprintf(" AND season_type = $%d", len(args))
		}
		query += " ORDER BY season DESC, season_type ASC, name ASC"

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []PlayerStat
		for rows.Next() {
			var p PlayerStat
			if err := rows.Scan(
				&p.Season, &p.SeasonType, &p.Name, &p.Team, &p.Pos, &p.Age, &p.G, &p.GS, &p.MP,
				&p.FG, &p.FGA, &p.FGPct, &p.FG3, &p.FG3A, &p.FG3Pct,
				&p.FT, &p.FTA, &p.FTPct, &p.ORB, &p.DRB, &p.TRB,
				&p.AST, &p.STL, &p.BLK, &p.TOV, &p.PF, &p.PTS,
			); err != nil {
				http.Error(w, "scan failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, p)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Printf("encoding /api/player response: %v", err)
		}
	})

	r.HandleFunc("/today", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/today.html")
	})

	r.HandleFunc("/api/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		scoreboard, err := espn.FetchScoreboard(date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, scoreboard)
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
