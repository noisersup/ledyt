package yt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/noisersup/ledyt/backend/common"
)

type YoutubeClient struct {
	C http.Client
}

func (c YoutubeClient) Search(query string) ([]common.Video, error) {
	uri, err := url.Parse("https://www.youtube.com/youtubei/v1/search")
	if err != nil {
		return nil, err
	}

	q := uri.Query()
	q.Add("key", "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8")
	uri.RawQuery = q.Encode()

	reqbody, err := searchBody(query)
	if err != nil {
		return nil, err
	}

	r, err := c.C.Post(uri.String(), "application/json", bytes.NewBuffer(reqbody))
	if err != nil {
		err = fmt.Errorf("[%d] %s", r.StatusCode, err.Error())
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse

	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	return getVideosFromResponse(searchResp.Contents.TwoColumn.PrimaryContents.SectionList.Contents[0].ItemSection.Contents), nil
}

func searchBody(query string) ([]byte, error) {
	search := SearchBody{
		Context{
			SearchClient{
				Hl:            "pl",
				Gl:            "PL",
				ClientName:    "WEB",
				ClientVersion: "2.20220201.01.00",
				TimeZone:      "Europe/Warsaw",
				UtcOffset:     60,
			},
		},
		query,
	}

	return json.Marshal(search)
}

func getVideosFromResponse(respVideos []ResponseVideo) []common.Video {
	outVideos := []common.Video{}
	for _, rv := range respVideos {
		channelRuns := rv.VideoRenderer.OwnerText.Runs
		videoRuns := rv.VideoRenderer.Title.Runs

		if len(videoRuns) == 0 || len(channelRuns) == 0 {
			continue //Ignore horizontal card lists, shelf renderers, etc...
		}

		ch := common.Channel{
			Name: channelRuns[0].Text,
		}
		v := common.Video{
			Title:   videoRuns[0].Text,
			URL:     fmt.Sprintf("https://youtube.com/watch?v=%s", rv.VideoRenderer.VideoId),
			Channel: &ch,
		}
		outVideos = append(outVideos, v)
	}
	return outVideos
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
		VideoId string `json:"videoId"`
		Title   struct {
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
