package twitter

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/twitter/mocks"
	"github.com/stretchr/testify/assert"
)

func newHTTPCli(f roundTripperFunc) *http.Client {
	return &http.Client{
		Transport: f,
		Timeout:   time.Millisecond * 100,
	}
}

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestFetchTweets(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		client := newHTTPCli(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mocks.SuccessResponse)),
			}, nil
		})

		api := New(client, "futile")

		resp, err := api.FetchTweets(context.Background(), logger.New("DEBUG"), FetchTweetsRequest{})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
	})

	t.Run("fail", func(t *testing.T) {
		client := newHTTPCli(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(bytes.NewBufferString(mocks.FailResponse)),
			}, nil
		})

		api := New(client, "futile")

		resp, err := api.FetchTweets(context.Background(), logger.New("DEBUG"), FetchTweetsRequest{})

		assert.Error(t, err)
		assert.Empty(t, resp)
	})
}
