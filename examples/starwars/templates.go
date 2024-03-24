package main

var (
	formats = map[int][]string{
		BrowseState:          []string{"BROWSE 📖", "Browse the results using the arrow keys - Item %d/%d", "API 🟢"},
		CheckConnectionState: []string{"CONNECTING", "Checking the connection with the API endpoint...", "API 🔴"},
		DownloadingState:     []string{"DOWNLOAD %s", "Fetching results from the endpoint", "API 🟢"},
	}
)
