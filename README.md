# NBA-API

To get started:

1. `go get github.com/umanchanda/NBA-API`

2. run: `go run ${GOPATH}/src/github.com/umanchanda/NBA-API/main.go` in your terminal

3. Head to `localhost:8000` in your browser

To run as a Docker build:

1. Clone the repo

2. `docker build -t nba .`

3. `docker run -p 8000:8000 nba`

4. Head to `localhost:8000` in your browser

How to use this repo:

To view all the boxscores for a given date, head over to `/boxscore/${year}/${month}/${date}`

To view team totals for a given game, head over to `/boxscore/${year}/${month}/${date}/${awayTeam}/${homeTeam}`
Note that `awayTeam` and `homeTeam` must be the three letter code for the team
So, for example, the `New York Knicks` are `NYK`, the `Cleveland Cavaliers`, are `CLE`, and so on so forth

To view player box scores for a given game, head over to `/boxscore/${year}/${month}/${date}/${awayTeam}/${homeTeam}/player`
Same caveats apply

For example, if I wanted to view the boxscores for March 11th, 2020 (coicidentally the last day there were NBA games before the league shut the season down), I would head to `localhost:8000/boxscore/2020/03/11`

If I wanted to check out the team boxscore for the Dallas v. Denver game, I would header over to `localhost:8000/boxscore/2020/03/11/DEN/DAL`

If I wanted to check out the player boxscore for that game, I would head over to `localhost:8000/boxscore/2020/03/11/DEN/DAL/player`
