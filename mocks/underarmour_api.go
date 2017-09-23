package mocks

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type UnderArmourAPI struct {
	server *httptest.Server
}

func underArmourMux(s *UnderArmourAPI) *http.ServeMux {
	mux := http.NewServeMux()
	// everything must end in a trailing slash
	mux.HandleFunc("/v7.1/user/", s.handleUser)
	mux.HandleFunc("/v7.1/route/", s.handleRoute)
	return mux
}

func NewUnderArmourAPI() *UnderArmourAPI {
	s := &UnderArmourAPI{}
	s.server = httptest.NewServer(underArmourMux(s))
	return s
}

func (s *UnderArmourAPI) URL() string {
	return s.server.URL
}

func (s *UnderArmourAPI) Close() error {
	s.server.Close()
	return nil
}
func hasAuthToken(w http.ResponseWriter, req *http.Request) (ok bool) {
	if c, err := req.Cookie("auth-token"); err != nil || c.Value != "some_token.123" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func (s *UnderArmourAPI) handleUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if !hasAuthToken(w, req) {
		return
	}
	parts := strings.Split(req.URL.Path, "/")
	if parts[len(parts)-2] != "self" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`
			{
				"_diagnostics": {
					"validation_failures": [
						{
							"__all__": [
								"You must specify friends_with, requested_friendship_with, suggested_friends_for, mutual_friends_for, q, email, name or username"
							]
						}
					]
				},
				"_links": {
					"self": [
						{
							"href": "/v7.1/user/?limit=20&offset=0"
						}
					],
					"documentation": [
						{
							"href": "https://developer.underarmour.com/docs/v71_User"
						}
					]
				}
			}
	`))
		return
	}
	w.Write([]byte(`{
		"username": "Zach123",
		"first_name": "Zach",
		"last_name": "Person",
		"display_name": "Zach Person",
		"last_initial": "P.",
		"preferred_language": "en-US",
		"introduction": "sup dog",
		"gender": "M",
		"location": {
			"country": "US",
			"region": "NY",
			"locality": "New York City"
		},
		"time_zone": "America/New_York",
		"goal_statement": null,
		"hobbies": "running",
		"profile_statement": "",
		"id": 117774799,
		"date_joined": "2017-07-07T12:35:28+00:00"
	}`))
}

func (s *UnderArmourAPI) handleRoute(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if !hasAuthToken(w, req) {
		return
	}
	parts := strings.Split(req.URL.Path, "/")
	if parts[len(parts)-2] != "1784229029" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"_diagnostics":{"validation_failures":[["a filter is required: (user, users, close_to_location, (city, state, country), text_search)"]]},"_links":{"self":[{"href":"\/v7.1\/route\/?limit=20&offset=0"}],"documentation":[{"href":"https:\/\/developer.underarmour.com\/docs\/v71_Route"}]}}`))
		return
	}
	w.Write([]byte(`{
		"total_descent": -57.9541283985,
		"city": "Norwalk",
		"data_source": "run:RE",
		"description": "",
		"updated_datetime": "2017-09-18T19:06:52+00:00",
		"created_datetime": "2017-09-18T19:06:43+00:00",
		"country": "us",
		"start_point_type": "",
		"starting_location": {
			"type": "Point",
			"coordinates": [
				-73.418,
				41.10955
			]
		},
		"distance": 14459.28,
		"name": "RUNNING RUNNERS - 9",
		"climbs": null,
		"state": "CT",
		"max_elevation": 22.21,
		"postal_code": "",
		"min_elevation": 0.0,
		"images": [
		],
		"_links": {
			"activity_types": [
				{
					"href": "/v7.1/activity_type/16/",
					"id": "16"
				}
			],
			"privacy": [
				{
					"href": "/v7.1/privacy_option/1/",
					"id": "1"
				}
			],
			"self": [
				{
					"href": "/v7.1/route/1784229029/",
					"id": "1784229029"
				}
			],
			"alternate": [
				{
					"href": "/v7.1/route/1784229029/?format=kml&field_set=detailed",
					"id": "1784229029",
					"name": "kml"
				},
				{
					"href": "/v7.1/route/1784229029/?format=gpx&field_set=detailed",
					"id": "1784229029",
					"name": "gpx"
				}
			],
			"user": [
				{
					"href": "/v7.1/user/117774782/",
					"id": "117774782"
				}
			],
			"thumbnail": [
				{
					"href": "//drzetlglcbfx.cloudfront.net/routes/thumbnail/1784229029/1505761612?size=100x100"
				}
			],
			"documentation": [
				{
					"href": "https://developer.underarmour.com/docs/v71_Route"
				}
			]
		},
		"points": null,
		"total_ascent": 54.1198459896
	}`))
}
