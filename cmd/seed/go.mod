module github.com/umanchanda/NBA-API/cmd/seed

go 1.21

require github.com/umanchanda/NBA-API/database v0.0.0

require (
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
)

replace github.com/umanchanda/NBA-API/database => ../../database
