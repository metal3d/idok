package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Release struct {
	Url     string `json:"html_url"`
	TagName string `json:"tag_name"`
	Pre     bool   `json:"prerelease"`
}

const REPO = "https://api.github.com/repos/metal3d/idok/releases"

func CheckRelease() (*Release, error) {
	resp, err := http.Get(REPO)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	r := []Release{}
	json.NewDecoder(resp.Body).Decode(&r)
	for _, release := range r {
		if !release.Pre {
			return &release, nil
		}
	}
	return nil, errors.New("Error fetching releases")
}
