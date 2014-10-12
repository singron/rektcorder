package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Twitch struct {
	ClientId string
}

var (
	broadcastRegex = regexp.MustCompile("^(?:(?:http://)?www\\.)?twitch\\.tv/[^/]*/b/(\\d+)$")
	videoIdRegex   = regexp.MustCompile("^((?:b|c)\\d+)$")
)

func (t *Twitch) Video(str string) (*Video, error) {
	var videoId string
	m := broadcastRegex.FindStringSubmatch(str)
	if m != nil {
		videoId = "a" + m[1]
	} else {
		m = videoIdRegex.FindStringSubmatch(str)
		if m == nil {
			return nil, fmt.Errorf("Could not extract videoId from '%s'\n", str)
		}
		videoId = m[1]
	}

	body, err := t.get("/videos/" + videoId)
	if err != nil {
		return nil, err
	}
	var v Video
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (t *Twitch) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken"+path, nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v2+json")
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s returned status %d\n", path, res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}

type Video struct {
	Title      string    `json: "title"`
	RecordedAt time.Time `json: "recorded_at"`
}
