package underarmour

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/blaskovicz/go-underarmour/models"
	gpx "github.com/ptrv/go-gpx"
)

const (
	DefaultRootURI = "https://api.ua.com"
	envVarPrefix   = "UNDERARMOUR"
)

type Client struct {
	rootURI         string
	cookieAuthToken string
	accessToken     string // TODO
}

func New() *Client {
	rootURI := os.Getenv(envVarPrefix + "_ROOT_URI")
	if rootURI == "" {
		rootURI = DefaultRootURI
	}
	cookieAuthToken := os.Getenv(envVarPrefix + "_COOKIE_AUTH_TOKEN")
	return &Client{rootURI: rootURI, cookieAuthToken: cookieAuthToken}
}

func (c *Client) uri(path string, pathArgs ...interface{}) string {
	return fmt.Sprintf("%s%s", c.rootURI, fmt.Sprintf(path, pathArgs...))
}
func (c *Client) doWithResponse(req *http.Request) (*http.Response, error) {
	if c.cookieAuthToken == "" {
		return nil, fmt.Errorf("missing auth-token for request")
	}
	req.AddCookie(&http.Cookie{Name: "auth-token", Value: c.cookieAuthToken})
	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	return resp, nil
}
func (c *Client) do(req *http.Request, decodeTarget interface{}) error {
	if decodeTarget != nil {
		if decodeKind := reflect.TypeOf(decodeTarget).Kind(); decodeKind != reflect.Ptr {
			return fmt.Errorf("invalid decode target type %s (need %s)", decodeKind.String(), reflect.Ptr.String())
		}
	}
	resp, err := c.doWithResponse(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if decodeTarget != nil {
			err = json.NewDecoder(resp.Body).Decode(decodeTarget)
			if err != nil {
				return fmt.Errorf("failed to decode payload: %s", err)
			}
		}
		return nil
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read %s response payload: %s", resp.Status, err)
	}
	if rawBody != nil && len(rawBody) != 0 {
		var e models.ErrorResponse
		if err = json.Unmarshal(rawBody, &e); err != nil {
			return fmt.Errorf("failed to decode %s error payload (%s): %s", resp.Status, string(rawBody), err)
		}
		return fmt.Errorf("request failed with %s: %v", resp.Status, e)
	} else {
		return fmt.Errorf("request failed with %s: %s", resp.Status, string(rawBody))
	}
}

// https://developer.underarmour.com/docs/v71_User/
func (c *Client) ReadUser(userPk string) (*models.User, error) {
	req, err := http.NewRequest("GET", c.uri("/v7.1/user/%s/", userPk), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var u models.User
	if err = c.do(req, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) ReadRoute(routeID int) (*models.Route, error) {
	req, err := http.NewRequest("GET", c.uri("/v7.1/route/%d/", routeID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	var r models.Route
	if err = c.do(req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) ReadRouteGPX(routeID int) (*gpx.Gpx, error) {
	req, err := http.NewRequest("GET", c.uri("/v7.1/route/%d/?format=gpx&field_set=detailed", routeID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	resp, err := c.doWithResponse(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return gpx.Parse(resp.Body)
}
