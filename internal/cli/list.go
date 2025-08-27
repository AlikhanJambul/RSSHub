package cli

import (
	"RSSHub/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ShowList(command *domain.Command) {
	body, err := json.Marshal(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	resp, err := http.Post("http://localhost:8080/list", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Something went wrong")
		os.Exit(1)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	var feeds []domain.Feed
	err = json.Unmarshal(responseBody, &feeds)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse response:", err)
		os.Exit(1)
	}

	for idx, feed := range feeds {
		text := fmt.Sprintf("%d Name: %s\n  URL: %s\n  Added: %s\n\n", idx+1, feed.Name, feed.URL, feed.CreatedAt)
		fmt.Fprintln(os.Stdout, text)
	}

	os.Exit(0)
}
