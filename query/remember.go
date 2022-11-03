package query

// Remember will add a query to the history.
// If query is already in the history, it will increment the rank by given weight
func Remember(query string, weight int) error {
	query = sanitize(query)

	cached, expired, err := cacher.Get()
	if expired || err != nil {
		cached = map[string]*queryRecord{}
	}

	if cached == nil {
		cached = make(map[string]*queryRecord)
	}

	// if the query is already in the cache
	// increment its rank
	if record, ok := cached[query]; ok {
		record.Rank += weight
	} else {
		cached[query] = &queryRecord{
			Rank:  weight,
			Query: query,
		}
	}

	return cacher.Set(cached)
}
