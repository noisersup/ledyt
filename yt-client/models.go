package ytclient

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

type SearchResponse struct {
	Contents struct {
		TwoColumn struct {
			PrimaryContents struct {
				SectionList struct {
					Contents []struct {
						ItemSection struct {
							Contents []ResponseVideo `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type ResponseVideo struct {
	VideoRenderer struct {
		Title struct {
			Runs []struct {
				Text string `json:"text"`
			} `json:"runs"`
		} `json:"title"`
		OwnerText struct {
			Runs []struct {
				Text string `json:"text"`
			} `json:"runs"`
		} `json:"ownerText"`
	} `json:"videoRenderer"`
}

type SearchBody struct {
	Context Context `json:"context"`
	Query   string  `json:"query"`
}

type Context struct {
	Client SearchClient `json:"client"`
}

type SearchClient struct {
	Hl            string `json:"hl"`
	Gl            string `json:"gl"`
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
	TimeZone      string `json:"timeZone"`
	UtcOffset     int    `json:"utcOffsetMinutes"`
}
