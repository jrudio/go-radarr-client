package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Movie ...
type Movie struct {
	Added      string `json:"added"`
	AddOptions struct {
		IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
		IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
		SearchForMovie             bool `json:"searchForMovie"`
	} `json:"addOptions"`
	AlternativeTitles []struct {
		Language   string `json:"language"`
		MovieID    int    `json:"movieId"`
		SourceID   int    `json:"sourceId"`
		SourceType string `json:"sourceType"`
		Title      string `json:"title"`
		VoteCount  int    `json:"voteCount"`
		Votes      int    `json:"votes"`
	} `json:"alternativeTitles"`
	CleanTitle       string   `json:"cleanTitle"`
	Deleted          bool     `json:"deleted"`
	Downloaded       bool     `json:"downloaded"`
	ErrorMessage     string   `json:"error"`
	EpisodeCount     int      `json:"episodeCount"`
	EpisodeFileCount int      `json:"episodeFileCount"`
	FolderName       string   `json:"folderName"`
	Genres           []string `json:"genres"`
	HasFile          bool     `json:"hasFile"`
	ID               int      `json:"id"`
	ImdbID           string   `json:"imdbId"`
	Images           []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	InCinemas           string `json:"inCinemas"`
	IsAvailable         bool   `json:"isAvailable"`
	IsExisting          bool   `json:"isExisting"`
	MinimumAvailability string `json:"minimumAvailability"`
	Monitored           bool   `json:"monitored"`
	Overview            string `json:"overview"`
	Path                string `json:"path"`
	PathState           string `json:"pathState"`
	PhysicalRelease     string `json:"physicalRelease"`

	ProfileID        int `json:"profileId"`
	QualityProfileID int `json:"qualityProfileId"`
	Ratings          struct {
		Value float64 `json:"value"`
		Votes int     `json:"votes"`
	} `json:"ratings"`
	RemotePoster          string   `json:"remotePoster"`
	RootFolderPath        string   `json:"rootFolderPath"`
	Runtime               int      `json:"runtime"`
	SecondaryYearSourceID int      `json:"secondaryYearSourceId"`
	SizeOnDisk            int      `json:"sizeOnDisk"`
	SortTitle             string   `json:"sortTitle"`
	Saved                 bool     `json:"saved"`
	Status                string   `json:"status"`
	Studio                string   `json:"studio"`
	Tags                  []string `json:"tags"`
	Title                 string   `json:"title"`
	TitleSlug             string   `json:"titleSlug"`
	TmdbID                int      `json:"tmdbId"`
	Year                  int      `json:"year"`
	YouTubeTrailerID      string   `json:"youTubeTrailerId"`
	Website               string   `json:"website"`
}

// AddMovie adds a movie to your wanted list
func (c Client) AddMovie(movie Movie) error {
	const endpoint = "/api/movie"

	// check required fields
	if movie.Title == "" {
		return errors.New("title is required")
	}

	if movie.QualityProfileID == 0 {
		return errors.New("quality profile id needs to be set")
	}

	if movie.TitleSlug == "" {
		return errors.New("title slug is required")
	}

	if len(movie.Images) == 0 {
		return errors.New("an array of images is required")
	}

	if movie.TmdbID == 0 {
		return errors.New("tmdbid is required")
	}

	if movie.Path == "" && movie.RootFolderPath == "" {
		return errors.New("either a path or rootFolderPath is required")
	}

	requestPayload, err := json.Marshal(movie)

	if err != nil {
		return err
	}

	resp, err := c.post(endpoint, requestPayload)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}

	return nil
}

// DeleteMovie removes a movie from your wanted list and/or local disk
// id is the id for the movie in the radarr library
func (c Client) DeleteMovie(id string, deleteFiles, addExclusion bool) error {
	const endpoint = "/api/movie/%s"

	params := make(url.Values, 1)

	params.Set("deleteFiles", strconv.FormatBool(deleteFiles))
	params.Set("addExclusion", strconv.FormatBool(addExclusion))

	resp, err := c.delete(fmt.Sprintf(endpoint, id), params)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
