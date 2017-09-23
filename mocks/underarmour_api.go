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
	mux.HandleFunc("/v7.1/user/self/", s.handleUser)
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
