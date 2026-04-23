module github.com/umanchanda/NBA-API

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	github.com/umanchanda/NBA-API/database v0.0.0
	github.com/umanchanda/NBA-API/playertotals v0.0.0-20200428185557-565476734505
	github.com/umanchanda/NBA-API/teamboxscore v0.0.0-20200428185557-565476734505
	github.com/umanchanda/NBA-API/teamtotals v0.0.0-20200428185557-565476734505
)

replace (
	github.com/umanchanda/NBA-API/database => ./database
	github.com/umanchanda/NBA-API/playertotals => ./playertotals
	github.com/umanchanda/NBA-API/teamboxscore => ./teamboxscore
	github.com/umanchanda/NBA-API/teamtotals => ./teamtotals
)
