package processor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestProcess(t *testing.T) {

	t.Run("failed fetching", func(t *testing.T) {
		fetcher := twitter.NewMockTweetsFetcher(t)
		classifier := predictor.NewMockFakeNewsClassifier(t)

		now := time.Now()
		expectedFetchRequest := twitter.FetchTweetsRequest{
			EntityID:   "entity",
			StartTime:  now,
			EndTime:    now,
			MaxResults: defaultFetchCount,
		}
		fetcher.On("FetchTweets", mock.Anything, mock.Anything, expectedFetchRequest).Return(
			twitter.FetchTweetsResponse{},
			errors.New("big error"),
		)

		process := GetProcessFn(logger.New("DEBUG"), fetcher, classifier)

		response := process(context.Background(), JobRequest{
			EntityID:  "entity",
			StartTime: now,
			EndTime:   now,
		})

		assert.Equal(t, "entity", response.EntityID)
		assert.Error(t, response.Error)

	})

	t.Run("failed classifying", func(t *testing.T) {
		fetcher := twitter.NewMockTweetsFetcher(t)
		classifier := predictor.NewMockFakeNewsClassifier(t)

		now := time.Now()
		expectedFetchRequest := twitter.FetchTweetsRequest{
			EntityID:   "entity",
			StartTime:  now,
			EndTime:    now,
			MaxResults: defaultFetchCount,
		}
		fetcher.On("FetchTweets", mock.Anything, mock.Anything, expectedFetchRequest).Return(
			twitter.FetchTweetsResponse([]twitter.Tweet{
				{
					ID:        "1",
					Text:      "Dummy 1",
					CreatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Dummy 2",
					CreatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Dummy 3",
					CreatedAt: now,
				},
			},
			),
			nil,
		)

		classifier.On("Classify", mock.Anything, predictor.ClassifyRequest([]string{
			"Dummy 1", "Dummy 2", "Dummy 3",
		})).Return(
			predictor.ClassifyResponse{},
			errors.New("big error"),
		)

		process := GetProcessFn(logger.New("DEBUG"), fetcher, classifier)

		response := process(context.Background(), JobRequest{
			EntityID:  "entity",
			StartTime: now,
			EndTime:   now,
		})

		assert.Equal(t, "entity", response.EntityID)
		assert.Error(t, response.Error)
	})

	t.Run("success", func(t *testing.T) {
		fetcher := twitter.NewMockTweetsFetcher(t)
		classifier := predictor.NewMockFakeNewsClassifier(t)

		now := time.Now()
		expectedFetchRequest := twitter.FetchTweetsRequest{
			EntityID:   "entity",
			StartTime:  now,
			EndTime:    now,
			MaxResults: defaultFetchCount,
		}
		fetcher.On("FetchTweets", mock.Anything, mock.Anything, expectedFetchRequest).Return(
			twitter.FetchTweetsResponse([]twitter.Tweet{
				{
					ID:        "1",
					Text:      "Dummy 1",
					CreatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Dummy 2",
					CreatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Dummy 3",
					CreatedAt: now,
				},
			},
			),
			nil,
		)

		classifier.On("Classify", mock.Anything, predictor.ClassifyRequest([]string{
			"Dummy 1", "Dummy 2", "Dummy 3",
		})).Return(
			predictor.ClassifyResponse{
				Classification: []predictor.Classification{
					predictor.Fake,
					predictor.Real,
					predictor.Fake,
				},
			},
			nil,
		)

		process := GetProcessFn(logger.New("DEBUG"), fetcher, classifier)

		response := process(context.Background(), JobRequest{
			EntityID:  "entity",
			StartTime: now,
			EndTime:   now,
		})

		assert.Equal(t, "entity", response.EntityID)
		assert.NoError(t, response.Error)
		assert.Equal(t, 2, len(response.FakeNewsTweets))
		assert.Equal(t, "Dummy 1", response.FakeNewsTweets[0].Content)
		assert.Equal(t, "Dummy 3", response.FakeNewsTweets[1].Content)
	})
}
