# NBA API

A Go web application for browsing NBA box scores and player season stats.

Live at: https://uman230-nba-api-2a9531748263.herokuapp.com

---

## How it works

The app has two main features:

### 1. Box Scores (scraped live from basketball-reference.com)

Select a date on the home page to see all games played that day. Each game card shows the final score, quarter-by-quarter breakdown, and a link to the full box score.

The full box score page shows each team's starters and reserves with per-game stats, plus a team totals row at the bottom.

**Routes:**

| Route | Description |
|---|---|
| `/scores/{year}/{month}/{day}` | All games for a given date |
| `/playerstats/{year}/{month}/{day}/{away}/{home}` | Full box score for a specific game |

Team codes are the standard three-letter abbreviations (e.g. `NYK`, `LAL`, `GSW`).

**Example:** Box score for the last day of the 2019-20 season before the COVID shutdown:
```
/scores/2020/03/11
```

### 2. Player Season Stats (served from a Neon PostgreSQL database)

Search for any player by name and optionally filter by season (year) and season type (Regular Season or Playoffs). Results show full season totals scraped from basketball-reference.com.

**Route:** `/searchPlayer`

---

## Running locally

```
git clone https://github.com/umanchanda/NBA-API
cd NBA-API
go mod tidy
go run .
```

Then open `http://localhost:8000`.

---

## Seeding the database

The player stats database is populated by scraping basketball-reference season totals pages. To seed all seasons from 1999-2000 to the current season:

```
export DATABASE_URL="postgres://..."
cd cmd/seed
go mod tidy
go run .
```

The seed command skips seasons already in the database, so it is safe to re-run.

---

## Environment variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | Neon PostgreSQL connection string |
| `PORT` | HTTP port (defaults to `8000`) |
