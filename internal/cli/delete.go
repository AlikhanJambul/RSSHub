package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"RSSHub/internal/domain"
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

	response, _ := io.ReadAll(resp.Body)

	var res struct {
		TextResponse string `json:"status"`
		ErrResponse  string `json:"error"`
	}

	err = json.Unmarshal(response, &res)

	if resp.StatusCode != http.StatusOK {

		fmt.Fprintln(os.Stderr, res.ErrResponse)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, res.TextResponse)
	os.Exit(0)
}
