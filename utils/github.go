package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Release struct {
	Url     string `json:"html_url"`
	TagName string `json:"tag_name"`
}

const REPO = "https://api.github.com/repos/metal3d/idok/releases"

func CheckRelease() Release {
	resp, err := http.Get(REPO)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	r := []Release{}
	json.NewDecoder(resp.Body).Decode(&r)
	return r[0]
}
