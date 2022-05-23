package main

type URL struct {
	Address  string
	Info     string
	Relation *URL
	Scraper  *Scraper
}
