package mixpanelproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/inverse-inc/go-utils/sharedutils"
	"github.com/tidwall/gjson"
)

// Sample request
// Request URL: https://api.mixpanel.com/track/?verbose=1&ip=0&_=1665604206698
// Request Method: POST
// Status Code: 200
// Remote Address: 35.190.25.25:443
// Referrer Policy: strict-origin-when-cross-origin
//
// [{
// 	"event": "route",
// 	"properties": {
// 		"$os": "Mac OS X",
// 		"$browser": "Chrome",
// 		"$referrer": "https://localhost:1443/admin",
// 		"$referring_domain": "localhost:1443",
// 		"$current_url": "https://localhost:1443/admin#/configuration/policies_access_control",
// 		"$browser_version": 106,
// 		"$screen_height": 1440,
// 		"$screen_width": 2560,
// 		"mp_lib": "web",
// 		"$lib_version": "2.45.0",
// 		"$insert_id": "4lw3v3g57zsq8g2h",
// 		"time": 1665604202.076,
// 		"distinct_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
// 		"$device_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
// 		"$initial_referrer": "$direct",
// 		"$initial_referring_domain": "$direct",
// 		"git_commit_id": [
// 			null
// 		],
// 		"version": [
// 			"12.1.1"
// 		],
// 		"_version": [
// 			null
// 		],
// 		"_git_commit_id": [
// 			null
// 		],
// 		"fromName": "statusDashboard",
// 		"fromUrl": "/status/dashboard",
// 		"toUrl": "/configuration/policies_access_control",
// 		"token": "changethishere"
// 	}
// }]

func TestValidMassageRequest(t *testing.T) {
	k := "Zammit4CEO"
	p := MixpanelProxy{MixpanelKey: k}

	eventsCount := 2
	b := []byte(`[{
	"event": "route",
	"properties": {
		"$os": "Mac OS X",
		"$browser": "Chrome",
		"$referrer": "https://localhost:1443/admin",
		"$referring_domain": "localhost:1443",
		"$current_url": "https://localhost:1443/admin#/configuration/policies_access_control",
		"$browser_version": 106,
		"$screen_height": 1440,
		"$screen_width": 2560,
		"mp_lib": "web",
		"$lib_version": "2.45.0",
		"$insert_id": "4lw3v3g57zsq8g2h",
		"time": 1665604202.076,
		"distinct_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
		"$device_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
		"$initial_referrer": "$direct",
		"$initial_referring_domain": "$direct",
		"git_commit_id": [
			null
		],
		"version": [
			"12.1.1"
		],
		"_version": [
			null
		],
		"_git_commit_id": [
			null
		],
		"fromName": "statusDashboard",
		"fromUrl": "/status/dashboard",
		"toUrl": "/configuration/policies_access_control",
		"token": "changethishere"
	}
},
{
	"event": "route",
	"properties": {
		"$os": "Mac OS X",
		"$browser": "Chrome",
		"$referrer": "https://localhost:1443/admin",
		"$referring_domain": "localhost:1443",
		"$current_url": "https://localhost:1443/admin#/configuration/policies_access_control",
		"$browser_version": 106,
		"$screen_height": 1440,
		"$screen_width": 2560,
		"mp_lib": "web",
		"$lib_version": "2.45.0",
		"$insert_id": "4lw3v3g57zsq8g2h",
		"time": 1665604202.076,
		"distinct_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
		"$device_id": "183c70fdf1f799-029d8ef95fcbe7-1a525635-1fa400-183c70fdf20c5c",
		"$initial_referrer": "$direct",
		"$initial_referring_domain": "$direct",
		"git_commit_id": [
			null
		],
		"version": [
			"12.1.1"
		],
		"_version": [
			null
		],
		"_git_commit_id": [
			null
		],
		"fromName": "statusDashboard",
		"fromUrl": "/status/dashboard",
		"toUrl": "/configuration/policies_access_control",
		"token": "changethishere"
	}
}]`)

	req, err := http.NewRequest("POST", "https://analytics.packetfence.org/track/?verbose=1&ip=0&_=1665604206698", bytes.NewBuffer(b))
	sharedutils.CheckError(err)

	err = p.MassageRequestBody(req)
	sharedutils.CheckError(err)

	newBody, err := ioutil.ReadAll(req.Body)
	sharedutils.CheckError(err)

	tokens := gjson.GetBytes(newBody, "#.properties.token").Array()
	if len(tokens) != eventsCount {
		t.Error("Invalid amount of tokens retrieved from the payload")
	}

	for i, res := range tokens {
		if res.String() != k {
			t.Errorf("Invalid token at index %d: '%s'", i, res.String())
		}
	}

	for _, toClear := range keysToClear {
		data := gjson.GetBytes(newBody, "#.properties."+toClear).Array()
		if len(data) != eventsCount {
			t.Error("Invalid amount of tokens retrieved from the payload")
		}

		for i, res := range data {
			if res.String() != clearedValue {
				t.Errorf("Invalid cleared field %s at index %d: '%s'", toClear, i, res.String())
			}
		}
	}

}

func TestInvalidMassageRequest(t *testing.T) {
	k := "Zammit4CEO"
	p := MixpanelProxy{MixpanelKey: k}

	b := []byte(`[{}]`)
	req, err := http.NewRequest("POST", "https://analytics.packetfence.org/track/?verbose=1&ip=0&_=1665604206698", bytes.NewBuffer(b))
	sharedutils.CheckError(err)

	err = p.MassageRequestBody(req)
	sharedutils.CheckError(err)

	newBody, err := ioutil.ReadAll(req.Body)
	sharedutils.CheckError(err)

	if string(newBody) != string(b) {
		t.Error("Badly handled a JSON containing an unknown event format")
	}

	b = []byte(`{}`)
	req, err = http.NewRequest("POST", "https://analytics.packetfence.org/track/?verbose=1&ip=0&_=1665604206698", bytes.NewBuffer(b))
	sharedutils.CheckError(err)

	err = p.MassageRequestBody(req)
	sharedutils.CheckError(err)

	newBody, err = ioutil.ReadAll(req.Body)
	sharedutils.CheckError(err)

	if string(newBody) != string(b) {
		t.Error("Badly handled a JSON containing an unknown event format")
	}

	b = []byte(`notjson`)
	req, err = http.NewRequest("POST", "https://analytics.packetfence.org/track/?verbose=1&ip=0&_=1665604206698", bytes.NewBuffer(b))
	sharedutils.CheckError(err)

	err = p.MassageRequestBody(req)
	sharedutils.CheckError(err)

	newBody, err = ioutil.ReadAll(req.Body)
	sharedutils.CheckError(err)

	if string(newBody) != string(b) {
		t.Error("Badly handled a payload that wasn't JSON")
	}
}
