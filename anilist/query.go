package anilist

import "fmt"

var mangaSubquery = `
id
idMal
title {
	romaji
	english
	native
}
description(asHtml: false)
tags {
	name
}
genres
coverImage {
	extraLarge
	large
	medium
	color
}
characters (page: 1, perPage: 10, role: MAIN) {
	nodes {
		name {
			full
		}
	}
}
startDate {
	year
	month	
	day
}
endDate {
	year
	month	
	day
}
status
synonyms
siteUrl
countryOfOrigin
externalLinks {
	url
}
`

var searchByNameQuery = fmt.Sprintf(`
query ($query: String) {
	Page (page: 1, perPage: 30) {
		media (search: $query, type: MANGA) {
			%s
		}
	}
}
`, mangaSubquery)

var searchByIDQuery = fmt.Sprintf(`
query ($id: Int) {
	Media (id: $id, type: MANGA) {
		%s
	}
}`, mangaSubquery)
