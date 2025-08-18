package models

type Command struct {
	Name     string `json:"name"`
	NameArg  string `json:"name_arg"`
	URL      string `json:"url"`
	Num      int    `json:"num"`
	Interval string `json:"interval"`
	Workers  int    `json:"workers"`
	FeedName string `json:"feed_name"`
}
