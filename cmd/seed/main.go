package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/umanchanda/NBA-API/database"
)

const firstSeason = 2000

func seedType(db *sql.DB, year int, seasonType string) {
	season := fmt.Sprintf("%d", year)

	exists, err := database.SeasonExists(db, season, seasonType)
	if err != nil {
		log.Printf("[%s/%s] skipping — could not check existence: %v", season, seasonType, err)
		return
	}
	if exists {
		log.Printf("[%s/%s] already seeded, skipping", season, seasonType)
		return
	}

	players, err := database.ScrapeTotals(season, seasonType)
	if err != nil {
		log.Printf("[%s/%s] scrape failed: %v", season, seasonType, err)
		return
	}

	for i := range players {
		players[i].Season = season
	}

	if err := database.InsertPlayers(db, players); err != nil {
		log.Printf("[%s/%s] insert failed: %v", season, seasonType, err)
		return
	}

	log.Printf("[%s/%s] done (%d players)", season, seasonType, len(players))
}

func main() {
	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	if err := database.CreateTable(db); err != nil {
		log.Fatalf("create table failed: %v", err)
	}

	currentYear := time.Now().Year()

	for year := firstSeason; year <= currentYear; year++ {
		seedType(db, year, database.SeasonTypeRegular)
		time.Sleep(2 * time.Second)

		seedType(db, year, database.SeasonTypePlayoffs)
		time.Sleep(2 * time.Second)
	}

	log.Println("all seasons seeded")
}
