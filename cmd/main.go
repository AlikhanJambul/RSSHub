package main

import (
	"RSSHub/internal/models"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func parseArgs() models.Command {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rsshub <command> [options]")
		os.Exit(1)
	}

	cmd := models.Command{Name: os.Args[1]}

	switch cmd.Name {
	case "fetch":
		// fetch без аргументов

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addCmd.StringVar(&cmd.NameArg, "name", "", "Feed name")
		addCmd.StringVar(&cmd.URL, "url", "", "RSS feed URL")
		addCmd.Parse(os.Args[2:])
		if cmd.NameArg == "" || cmd.URL == "" {
			fmt.Println("Error: --name and --url are required")
			os.Exit(1)
		}

	case "set-interval":
		intCmd := flag.NewFlagSet("set-interval", flag.ExitOnError)
		intCmd.StringVar(&cmd.Interval, "interval", "", "Interval duration (e.g. 3m, 10s)")
		intCmd.Parse(os.Args[2:])
		if cmd.Interval == "" {
			fmt.Println("Error: --interval is required")
			os.Exit(1)
		}

	case "set-workers":
		workersCmd := flag.NewFlagSet("set-workers", flag.ExitOnError)
		workersCmd.IntVar(&cmd.Workers, "count", 0, "Number of workers")
		workersCmd.Parse(os.Args[2:])
		if cmd.Workers <= 0 {
			fmt.Println("Error: --count must be > 0")
			os.Exit(1)
		}

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.IntVar(&cmd.Num, "num", 0, "Number of feeds to display")
		listCmd.Parse(os.Args[2:])

	case "delete":
		delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		delCmd.StringVar(&cmd.NameArg, "name", "", "Feed name to delete")
		delCmd.Parse(os.Args[2:])
		if cmd.NameArg == "" {
			fmt.Println("Error: --name is required")
			os.Exit(1)
		}

	case "articles":
		artCmd := flag.NewFlagSet("articles", flag.ExitOnError)
		artCmd.StringVar(&cmd.FeedName, "feed-name", "", "Feed name to display articles from")
		artCmd.IntVar(&cmd.Num, "num", 3, "Number of articles to display (default 3)")
		artCmd.Parse(os.Args[2:])
		if cmd.FeedName == "" {
			fmt.Println("Error: --feed-name is required")
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", cmd.Name)
		os.Exit(1)
	}

	return cmd
}

func main() {
	cmd := parseArgs()

	switch cmd.Name {
	case "fetch":
		if err := http.ListenAndServe("8080", nil); err != nil {
			log.Fatal(err)
		}
	case "add":
		fmt.Printf("Adding feed: %s (%s)\n", cmd.NameArg, cmd.URL)
	case "set-interval":
		fmt.Printf("Changing interval to: %s\n", cmd.Interval)
	case "set-workers":
		fmt.Printf("Changing workers to: %d\n", cmd.Workers)
	case "list":
		fmt.Printf("Listing feeds (limit: %d)\n", cmd.Num)
	case "delete":
		fmt.Printf("Deleting feed: %s\n", cmd.NameArg)
	case "articles":
		fmt.Printf("Listing %d articles from feed: %s\n", cmd.Num, cmd.FeedName)
	}
}
