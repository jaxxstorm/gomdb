// Package gomdb is a golang implementation of the OMDB API.
package gomdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseURL  = "http://www.omdbapi.com/?"
	plot     = "full"
	tomatoes = "true"

	MovieSearch   = "movie"
	SeriesSearch  = "series"
	EpisodeSearch = "episode"
)

// QueryData is the type to create the search query
type QueryData struct {
	Title      string
	Year       string
	ImdbId     string
	SearchType string // The type of search - "movie", "series", "episode"
	Page       string // In case of a paginated result, retrieve a particular page
}

// SearchResult is the type for the search results
type SearchResult struct {
	Title  string
	Year   string
	ImdbID string
	Type   string // The type of search - "movie", "series", "episode"
}

// SearchResponse is the struct of the response in a search
type SearchResponse struct {
	Search   []SearchResult
	Response string
	Error    string
	NumPages string `json:"totalResults"`
}

// MovieResult is the result struct of an specific movie search
type MovieResult struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Metascore  string
	ImdbRating string
	ImdbVotes  string
	ImdbID     string
	Ratings    []Ratings
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
	Error      string
}

type Ratings struct {
	Source string
	Value  string
}

// Search returns a SearchResponse struct. The search query is within the QueryData
// struct
func Search(query *QueryData) (*SearchResponse, error) {
	resp, err := requestAPI("search",
		query.Title,
		query.Year,
		query.SearchType,
		query.Page)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := new(SearchResponse)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}
	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// LookupByTitle returns a MovieResult
func LookupByTitle(query *QueryData) (*MovieResult, error) {
	resp, err := requestAPI("title", query.Title, query.Year, query.SearchType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := new(MovieResult)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}
	if r.Response == "False" {
		return r, errors.New(r.Error)
	}
	return r, nil
}

// LookupByImdbID returns a MovieResult given a ImdbID ex:"tt2015381"
func LookupByImdbID(id string) (*MovieResult, error) {
	resp, err := requestAPI("id", id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := new(MovieResult)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}
	if r.Response == "False" {
		return r, errors.New(r.Error)
	}
	return r, nil
}

// helper function to call the API
// param: apiCategory refers to which API we are calling. Can be "search", "title" or "id"
// Depending on that value, we will search by "t" or "s" or "i"
// param: params are the variadic list of params passed for that category
func requestAPI(apiCategory string, params ...string) (resp *http.Response, err error) {
	var URL *url.URL
	URL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// Checking for invalid category
	if len(params) > 1 && params[2] != "" {
		if params[2] != MovieSearch &&
			params[2] != SeriesSearch &&
			params[2] != EpisodeSearch {
			return nil, errors.New("Invalid search category- " + params[2])
		}
	}
	URL.Path += "/"
	parameters := url.Values{}
	switch apiCategory {
	case "search":
		parameters.Add("s", params[0])
		parameters.Add("y", params[1])
		parameters.Add("type", params[2])
		parameters.Add("page", params[3])
	case "title":
		parameters.Add("t", params[0])
		parameters.Add("y", params[1])
		parameters.Add("type", params[2])
		parameters.Add("plot", plot)
		parameters.Add("tomatoes", tomatoes)
	case "id":
		parameters.Add("i", params[0])
		parameters.Add("plot", plot)
		parameters.Add("tomatoes", tomatoes)
	}

	URL.RawQuery = parameters.Encode()
	res, err := http.Get(URL.String())
	err = checkErr(res.StatusCode)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func checkErr(status int) error {
	if status != 200 {
		return fmt.Errorf("Status Code %d received from IMDB", status)
	}
	return nil
}

// Stringer Interface for MovieResult
func (mr MovieResult) String() string {
	return fmt.Sprintf("#%s: %s (%s)", mr.ImdbID, mr.Title, mr.Year)
}

// Stringer Interface for SearchResult
func (sr SearchResult) String() string {
	return fmt.Sprintf("#%s: %s (%s) Type: %s", sr.ImdbID, sr.Title, sr.Year, sr.Type)
}
