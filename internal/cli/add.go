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

func AddFeed(command *domain.Command) {
	body, err := json.Marshal(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	resp, err := http.Post("http://localhost:8080/add", "application/json", bytes.NewBuffer(body))
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
