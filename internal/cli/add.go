package cli

import (
	"RSSHub/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func AddFeed(command models.Command) {
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

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		fmt.Fprintln(os.Stderr, "Something went wrong")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Success")
	os.Exit(0)
}
