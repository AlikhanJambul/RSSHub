package main

import (
	"RSSHub/internal/cli"
	"RSSHub/internal/domain"
	"flag"
	"fmt"
	"os"
)

func parseArgs() *domain.Command {
	if len(os.Args) < 2 {
		helpPrint()
		os.Exit(0)
	}

	cmd := domain.Command{Name: os.Args[1]}

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
		var count int

		workersCmd := flag.NewFlagSet("set-workers", flag.ExitOnError)
		workersCmd.IntVar(&count, "count", 0, "Number of workers")
		workersCmd.Parse(os.Args[2:])
		if count <= 0 {
			fmt.Println("Error: --count must be > 0")
			os.Exit(1)
		}
		cmd.Workers = int32(count)

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.IntVar(&cmd.Num, "num", -1, "Number of feeds to display")
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
	case "help":
		helpPrint()
		os.Exit(0)

	default:
		fmt.Printf("Unknown command: %s\n", cmd.Name)
		os.Exit(1)
	}

	return &cmd
}

func main() {
	cmd := parseArgs()

	switch cmd.Name {
	case "fetch":
		cli.StartServer()
	case "add":
		cli.AddFeed(cmd)
	case "set-interval":
		cli.SetInteval(cmd)
	case "set-workers":
		cli.SetWorkersCount(cmd)
	case "list":
		cli.ShowList(cmd)
	case "delete":
		cli.DeleteFeed(cmd)
	case "articles":
		fmt.Printf("Listing %d articles from feed: %s\n", cmd.Num, cmd.FeedName)
	}
}

func helpPrint() {
	helpText := `
  Usage:
    rsshub COMMAND [OPTIONS]

  Common Commands:
       add             add new RSS feed
       set-interval    set RSS fetch interval
       set-workers     set number of workers
       list            list available RSS feeds
       delete          delete RSS feed
       articles        show latest articles
       fetch           starts the background process that periodically fetches and processes RSS feeds using a worker pool`

	fmt.Println(helpText)
}
