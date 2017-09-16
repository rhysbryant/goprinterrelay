package main

import (
	"net/http"
)

type ApplicationInfo struct {
	Version       string `json:"version"`
	FeatureConfig struct {
		Camera struct {
			AutoStart bool `json:"autoStart"`
		} `json:"camera"`
	} `json:"featureConfig"`
}

func ApplicationInfoGetRequest(w http.ResponseWriter, r *http.Request) {
	WriteJsonResponse(w, applicationInfo)
}
