package version

import "time"

// Version information with a link to the details of the Version
type Version struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`

	Published time.Time `json:"time"`
	Released  time.Time `json:"releaseTime"`
	Downloads struct {
		Client struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
			SHA1 string `json:"sha1"`
		} `json:"client"`
		Server struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
			SHA1 string `json:"sha1"`
		} `json:"server"`
	} `json:"downloads"`
	Java struct {
		Version int `json:"majorVersion"`
	} `json:"javaVersion"`
}
