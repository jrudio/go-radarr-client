package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jrudio/go-radarr-client"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func startDB() (store, error) {
	// create persistent key store in user home directory
	storeDirectory, err := homedir.Dir()

	if err != nil {
		return store{}, err
	}

	storeDirectory = filepath.Join(storeDirectory, homeFolderName)

	return initDataStore(storeDirectory)
}

func unlock(c *cli.Context) error {
	storeDirectory, err := homedir.Dir()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	storeDirectory = filepath.Join(storeDirectory, homeFolderName)
	lockFilePath := filepath.Join(storeDirectory, "LOCK")

	if err := os.Remove(lockFilePath); err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to remove file: %v", err), 1)
	}

	fmt.Println("removed LOCK file")

	return nil
}

func save(c *cli.Context) error {
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// prompt to save url to radarr application
	fmt.Println("enter the url that points to radarr...")

	var radarrURL string

	fmt.Scanln(&radarrURL)

	if radarrURL == "" {
		return cli.NewExitError("url is required", 1)
	}

	// prompt to save api key
	fmt.Println("enter your api key...")

	var key string

	fmt.Scanln(&key)

	if key == "" {
		return cli.NewExitError("api key is required", 1)
	}

	// confirm
	fmt.Printf("URL: %s\nAPI Key: %s\n", radarrURL, key)
	fmt.Println("Are you sure you want to save?")

	// show success/error

	if err := db.saveRadarrURL(radarrURL); err != nil {
		return cli.NewExitError(fmt.Sprintf("save url failed: %v", err), 1)
	}

	if err := db.saveRadarrKey(key); err != nil {
		// revert url save
		db.saveRadarrURL("")

		return cli.NewExitError(fmt.Sprintf("save api key failed: %v", err), 1)
	}

	return nil
}

func getCredentials(c *cli.Context) error {
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	radarrURL, err := db.getRadarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	key, err := db.getRadarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("URL: %s\nAPI Key: %s\n", radarrURL, key)

	return nil
}

func search(c *cli.Context) error {
	title := strings.Join(c.Args(), " ")

	if title == "" {
		return cli.NewExitError("a title is required", 1)
	}

	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	radarrKey, err := db.getRadarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	radarrURL, err := db.getRadarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	client, err := radarr.New(radarrURL, radarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	results, err := client.Search(title)

	for _, movie := range results {
		fmt.Printf("%s (%d) - %d\n", movie.Title, movie.Year, movie.TmdbID)
	}

	return nil
}

func showMovieInfo(c *cli.Context) error {
	tmdbID := c.Args().First()

	if tmdbID == "" {
		return cli.NewExitError("a tmdb id is required", 1)
	}

	// fire up store
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// grab credentials
	radarrKey, err := db.getRadarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	radarrURL, err := db.getRadarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// create radarr client to interface with radarr
	client, err := radarr.New(radarrURL, radarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	tmdbIDStr, err := strconv.Atoi(tmdbID)

	movie, err := client.GetMovie(tmdbIDStr)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// title (year) - tmdbid
	// 		summary
	const output = "%s (%d) - %d\n\t%s\n"

	fmt.Printf(output, movie.Title, movie.Year, movie.TmdbID, movie.Overview)

	return nil
}

func showLibrary(c *cli.Context) error {
	// fire up store
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// grab credentials
	radarrKey, err := db.getRadarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	radarrURL, err := db.getRadarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// create radarr client to interface with radarr
	client, err := radarr.New(radarrURL, radarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	movies, err := client.GetMovies(radarr.GetMovieOptions{})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for _, movie := range movies {
		status := "missing"

		if movie.Downloaded {
			status = "downloaded"
		}

		fmt.Printf("%s (%d) - %s\n", movie.Title, movie.Year, status)
	}

	return nil
}

func addMovie(c *cli.Context) error {
	tmdbID := c.Args().First()

	if tmdbID == "" {
		return cli.NewExitError("a tmdb id is required", 1)
	}

	// fire up store
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// grab credentials
	radarrKey, err := db.getRadarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	radarrURL, err := db.getRadarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// create radarr client to interface with radarr
	client, err := radarr.New(radarrURL, radarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	tmdbIDStr, err := strconv.Atoi(tmdbID)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	movie, err := client.GetMovie(tmdbIDStr)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// show available profiles
	profiles, err := client.GetProfiles()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	profileCount := len(profiles)

	if profileCount == 0 {
		fmt.Println("aborting...")
		return cli.NewExitError("no profiles found", 1)
	}

	fmt.Print("available quality profiles:\n\n")

	for i, profile := range profiles {
		fmt.Printf("[%d] - %s\n", i, profile.Name)
	}

	fmt.Print("\nplease choose a profile: ")

	// ask user for requested quality
	var requestedQualityIndex int
	fmt.Scanln(&requestedQualityIndex)

	// bound-check user input
	if requestedQualityIndex < 0 || requestedQualityIndex > profileCount {
		return cli.NewExitError("invalid selection", 1)
	}

	profile := profiles[requestedQualityIndex].ID

	// display available root folders
	folders, err := client.GetRootFolders()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(folders) == 0 {
		fmt.Println("aborting...")
		return cli.NewExitError("failed to find root folders", 1)
	}

	fmt.Println("\navailable root folders:")

	for i, folder := range folders {
		fmt.Printf("[%d] - %s\n", i, folder.Path)
	}

	fmt.Print("\nchoose a folder to download this movie to: ")

	// ask user where we should download this movie to
	var rootFolderPathIndex int
	fmt.Scanln(&rootFolderPathIndex)

	fmt.Println()

	rootFolder := folders[rootFolderPathIndex].Path

	// set movie path and profile quality to user preference
	movie.AddOptions.SearchForMovie = true
	movie.QualityProfileID = profile
	movie.RootFolderPath = rootFolder
	movie.Monitored = true

	if errors := client.AddMovie(movie); errors != nil {
		output := ""

		for _, err := range errors {
			output += err.Error() + "\n"
		}

		fmt.Printf(output)

		return cli.NewExitError(fmt.Errorf(""), 1)
	}

	fmt.Printf("added %s (%d) successfully\n", movie.Title, movie.Year)

	return nil
}

func deleteMovie(c *cli.Context) error {
	return nil
}
