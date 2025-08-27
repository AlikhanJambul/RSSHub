package domain

type Command struct {
	Name     string `json:"name"`
	NameArg  string `json:"name_arg"`
	URL      string `json:"url"`
	Num      int    `json:"num"`
	Interval string `json:"interval"`
	Workers  int32  `json:"workers"`
	FeedName string `json:"feed_name"`
}
