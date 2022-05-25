package main

type URL struct {
	Relation *URL
	Scraper  *Scraper
	Address  string
	Info     string
	Index    int
}
