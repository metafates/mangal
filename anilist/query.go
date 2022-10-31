package anilist

import "fmt"

// mangaSubquery common manga query used for getting manga by id or searching it by name
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
	description
	rank
}
genres
coverImage {
	extraLarge
	large
	medium
	color
}
bannerImage
characters (page: 1, perPage: 10, role: MAIN) {
	nodes {
		id
		name {
			full
			native
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
staff {
	edges {
	  role
	  node {
		name {
		  full
		}
	  }
	}
}
status
synonyms
siteUrl
chapters
countryOfOrigin
externalLinks {
	url
}
`

// searchByNameQuery query used for searching manga by name
var searchByNameQuery = fmt.Sprintf(`
query ($query: String) {
	Page (page: 1, perPage: 30) {
		media (search: $query, type: MANGA) {
			%s
		}
	}
}
`, mangaSubquery)

// searchByIDQuery query used for searching manga by id
var searchByIDQuery = fmt.Sprintf(`
query ($id: Int) {
	Media (id: $id, type: MANGA) {
		%s
	}
}`, mangaSubquery)
