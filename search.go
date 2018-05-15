package radarr

import (
	"encoding/json"
	"errors"
	"net/http"
)

// SearchResults is the returned results of a search
type SearchResults struct {
	Added             string   `json:"added"`
	AlternativeTitles []string `json:"alternativeTitles"`
	Downloaded        bool     `json:"downloaded"`
	FolderName        string   `json:"folderName"`
	Genres            []string `json:"genres"`
	HasFile           bool     `json:"hasFile"`
	Images            []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	InCinemas           string `json:"inCinemas"`
	IsAvailable         bool   `json:"isAvailable"`
	MinimumAvailability string `json:"minimumAvailability"`
	Monitored           bool   `json:"monitored"`
	Overview            string `json:"overview"`
	PathState           string `json:"pathState"`
	ProfileID           int    `json:"profileId"`
	QualityProfileID    int    `json:"qualityProfileId"`
	Ratings             struct {
		Value float64 `json:"value"`
		Votes int     `json:"votes"`
	} `json:"ratings"`
	RemotePoster          string   `json:"remotePoster"`
	Runtime               int      `json:"runtime"`
	SecondaryYearSourceID int      `json:"secondaryYearSourceId"`
	SizeOnDisk            int      `json:"sizeOnDisk"`
	SortTitle             string   `json:"sortTitle"`
	Status                string   `json:"status"`
	Tags                  []string `json:"tags"`
	Title                 string   `json:"title"`
	TitleSlug             string   `json:"titleSlug"`
	TmdbID                int      `json:"tmdbId"`
	Year                  int      `json:"year"`
	ErrorMessage          string   `json:"error"`
}

// Search uses Radarr's method of online movie lookup
func (c Client) Search(title string) ([]SearchResults, error) {
	encodedTitle, err := encodeURL(title)

	if err != nil {
		return []SearchResults{}, err
	}

	resp, err := c.get(c.URL + "/api/movies/lookup?term=" + encodedTitle)

	if err != nil {
		return []SearchResults{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []SearchResults{}, errors.New(resp.Status)
	}

	var results []SearchResults

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return results, err
	}

	return results, nil
}

// SearchOffline searches for movies already in Radarr's library
func (c Client) SearchOffline(title string) (SearchResults, error) {
	return SearchResults{}, nil
}
