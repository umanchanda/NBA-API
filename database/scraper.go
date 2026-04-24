package database

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.basketball-reference.com"

// coalesce returns the text of the first data-stat that has a value.
func coalesce(row *goquery.Selection, stats ...string) string {
	for _, s := range stats {
		if v := strings.TrimSpace(row.Find("td[data-stat='" + s + "'] a").Text()); v != "" {
			return v
		}
		if v := strings.TrimSpace(row.Find("td[data-stat='" + s + "']").Text()); v != "" {
			return v
		}
	}
	return ""
}

// ScrapeTotals fetches player totals for the given year and season type.
func ScrapeTotals(year, seasonType string) ([]NBAPlayer, error) {
	var url string
	if seasonType == SeasonTypePlayoffs {
		url = baseURL + "/playoffs/NBA_" + year + "_totals.html"
	} else {
		url = baseURL + "/leagues/NBA_" + year + "_totals.html"
	}
	log.Printf("fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching totals: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("HTTP %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	rows := doc.Find("#totals_stats tbody tr")
	log.Printf("found %d rows in #totals_stats", rows.Length())

	var players []NBAPlayer

	rows.Each(func(_ int, row *goquery.Selection) {
		if row.HasClass("thead") {
			return
		}

		stat := func(name string) string {
			return strings.TrimSpace(row.Find("td[data-stat='" + name + "']").Text())
		}

		name := strings.TrimSpace(row.Find("td[data-stat='name_display'] a").Text())
		if name == "" {
			name = strings.TrimSpace(row.Find("td[data-stat='player'] a").Text())
		}
		if name == "" {
			return
		}

		players = append(players, NBAPlayer{
			SeasonType: seasonType,
			Name:       name,
			Team:       coalesce(row, "team_name_abbr", "team_id"),
			Pos:        stat("pos"),
			Age:        stat("age"),
			G:          coalesce(row, "games", "g"),
			GS:         coalesce(row, "games_started", "gs"),
			MP:         stat("mp"),
			FG:         stat("fg"),
			FGA:        stat("fga"),
			FGPct:      stat("fg_pct"),
			FG3:        stat("fg3"),
			FG3A:       stat("fg3a"),
			FG3Pct:     stat("fg3_pct"),
			FT:         stat("ft"),
			FTA:        stat("fta"),
			FTPct:      stat("ft_pct"),
			ORB:        stat("orb"),
			DRB:        stat("drb"),
			TRB:        stat("trb"),
			AST:        stat("ast"),
			STL:        stat("stl"),
			BLK:        stat("blk"),
			TOV:        stat("tov"),
			PF:         stat("pf"),
			PTS:        stat("pts"),
		})
	})

	log.Printf("scraped %d players", len(players))

	if len(players) == 0 {
		return nil, fmt.Errorf("no players found — table selector may have changed or no playoff data for %s", year)
	}

	return players, nil
}
