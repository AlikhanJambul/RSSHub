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

func DeleteFeed(command *domain.Command) {
	body, err := json.Marshal(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/delete", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response, _ := io.ReadAll(resp.Body)

		fmt.Fprintln(os.Stderr, string(response))
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Feed has been deleted")
	os.Exit(0)
}
