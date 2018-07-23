package radarr

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Search uses Radarr's method of online movie lookup
func (c Client) Search(title string) ([]Movie, error) {
	params := url.Values{}

	params.Set("term", title)

	resp, err := c.get("/api/movie/lookup", params)

	if err != nil {
		return []Movie{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Movie{}, errors.New(resp.Status)
	}

	var results []Movie

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return results, err
	}

	return results, nil
}

// SearchOffline searches for movies already in Radarr's library
func (c Client) SearchOffline(title string) (Movie, error) {
	return Movie{}, nil
}

// GetMovie returns a movie via the movie database id
func (c Client) GetMovie(tmdbID int) (Movie, error) {
	const endpoint = "/api/movies/lookup/tmdb"

	params := url.Values{}

	params.Set("tmdbId", strconv.Itoa(tmdbID))

	matchedMovie := Movie{}

	resp, err := c.get(endpoint, params)

	if err != nil {
		return matchedMovie, err
	}

	defer resp.Body.Close()

	// handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return matchedMovie, errors.New(resp.Status)
	}

	var result Movie

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return matchedMovie, err
	}

	return result, nil
}

// GetMovieIMDB returns a movie via the internet movie database id
func (c Client) GetMovieIMDB(imdbID int) (Movie, error) {
	const endpoint = "/api/movies/lookup/imdb"

	params := url.Values{}

	params.Set("imdbId", strconv.Itoa(imdbID))

	var result Movie

	resp, err := c.get(endpoint, params)

	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
