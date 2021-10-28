package twitter

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
)

func TestToEvidenve(t *testing.T) {
	t.Run("lookup tweet", func(t *testing.T) {
		expectedText := "text"
		expectedUserID := "userID"
		evidence, err := toEvidence([]twitter.Tweet{
			{
				Text: "text",
				User: &twitter.User{
					ScreenName: "userID",
				},
			},
		}, 0)
		assert.Nil(t, err)
		assert.Equal(t, evidence.text, expectedText)
		assert.Equal(t, evidence.userID, expectedUserID)
	})
	t.Run("error if tweet not found", func(t *testing.T) {
		_, err := toEvidence([]twitter.Tweet{}, 0)
		assert.Error(t, err)
	})
}
