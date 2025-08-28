package cli

import (
	"RSSHub/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ShowList(command *domain.Command) {
	url := fmt.Sprintf("http://localhost:8080/list?count=%d", command.Num)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	var res struct {
		ErrResponse string        `json:"error"`
		Feeds       []domain.Feed `json:"feeds"`
	}

	err = json.Unmarshal(response, &res)

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, res.ErrResponse)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse response:", err)
		os.Exit(1)
	}

	for idx, feed := range res.Feeds {
		text := fmt.Sprintf("%d Name: %s\n  URL: %s\n  Added: %s\n\n", idx+1, feed.Name, feed.URL, feed.CreatedAt)
		fmt.Fprintln(os.Stdout, text)
	}

	os.Exit(0)
}
