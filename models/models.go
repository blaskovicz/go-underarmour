package models

import "time"

type ErrorResponse struct {
	// eg: diagnostics.validation_failures[0].__all__[0]
	Diagnostics map[string][]map[string][]string `json:"_diagnostics,omitempty"`
	// eg: links.self[0].href
	Links map[string][]map[string]string `json:"_links,omitempty"`
}
type User struct {
	ID                int       `json:"id"`
	Gender            string    `json:"gender"`
	PreferredLanguage string    `json:"preferred_language"`
	Introduction      string    `json:"introduction"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	DisplayName       string    `json:"display_name"`
	LastInitial       string    `json:"last_initial"`
	DateJoined        time.Time `json:"date_joined"`
	ProfileStatement  string    `json:"profile_statement"`
	Hobbies           string    `json:"hobbies"`
	TimeZone          string    `json:"time_zone"`
	GoalStatement     string    `json:"goal_statement"`
	Location          struct {
		Country  string `json:"country"`
		Region   string `json:"region"`
		Locality string `json:"locality"`
	} `json:"location"`
}
