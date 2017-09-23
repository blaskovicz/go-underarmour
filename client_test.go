package underarmour

import (
	"os"
	"testing"
	"time"

	"github.com/blaskovicz/go-underarmour/mocks"
	"github.com/blaskovicz/go-underarmour/models"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	s := mocks.NewUnderArmourAPI()
	defer s.Close()
	os.Setenv("UNDERARMOUR_ROOT_URI", s.URL())
	os.Setenv("UNDERARMOUR_COOKIE_AUTH_TOKEN", "some_token.123")
	var err error
	var client *Client
	var user *models.User
	t.Run("Init", func(t *testing.T) {
		client = New()
		require.NotNil(t, client, "client was nil")
	})
	t.Run("ReadUser", func(t *testing.T) {
		require.NotNil(t, client)
		user, err = client.ReadUser("self")
		require.NoError(t, err, "read user failed")
		require.NotNil(t, user, "user was nil")
		require.Equal(t, "Zach", user.FirstName)
		require.Equal(t, "Person", user.LastName)
		require.Equal(t, "Zach123", user.Username)
		require.Equal(t, "Zach Person", user.DisplayName)
		require.Equal(t, "P.", user.LastInitial)
		require.Equal(t, "M", user.Gender)
		require.Equal(t, "en-US", user.PreferredLanguage)
		require.Equal(t, "New York City", user.Location.Locality)
		require.Equal(t, "NY", user.Location.Region)
		require.Equal(t, "US", user.Location.Country)
		require.Equal(t, "running", user.Hobbies)
		require.Equal(t, "sup dog", user.Introduction)
		require.Equal(t, "America/New_York", user.TimeZone)
		require.Equal(t, "", user.GoalStatement)
		require.Equal(t, "", user.ProfileStatement)
		require.Equal(t, 2017, user.DateJoined.Year())
		require.Equal(t, time.Month(7), user.DateJoined.Month())
		require.Equal(t, 7, user.DateJoined.Day())
		require.Equal(t, 117774799, user.ID)
	})
}
