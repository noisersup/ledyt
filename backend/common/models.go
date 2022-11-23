package common

type Video struct {
	Title   string
	Channel *Channel
	URL     string
	//TODO: length
}

type Channel struct {
	Name string
	URL  string
}
