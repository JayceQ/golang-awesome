package engine

import "golang-awesome/crawler/fetcher"

func Worker(r Request)(ParserResult, error){
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		return ParserResult{},err
	}

	return r.ParserFunc
}