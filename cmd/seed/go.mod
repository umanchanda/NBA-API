module github.com/umanchanda/NBA-API/cmd/seed

go 1.25.0

require github.com/umanchanda/NBA-API/database v0.0.0

require (
	github.com/PuerkitoBio/goquery v1.12.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/lib/pq v1.12.3 // indirect
	golang.org/x/net v0.53.0 // indirect
)

replace github.com/umanchanda/NBA-API/database => ../../database
