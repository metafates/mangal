package anilist

var searchQuery = `
query ($query: String) {
	Page (page: 1, perPage: 30) {
		media (search: $query, type: MANGA) {
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
		}
	}
}
`
