package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
)

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

type getUserTweetsResponse struct {
	Data []tweet  `json:"data"`
	Meta metadata `json:"meta"`
}

type tweet struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Text      string    `json:"text"`
}

// metadata left to enable pagination option in perspective
// can be removed if needed
type metadata struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

func (client *Client) FetchTweets(ctx context.Context, log logger.Interface, ftr FetchTweetsRequest) (FetchTweetsResponse, error) {
	resp, err := client.invokeFetchTweets(ctx, log, ftr.EntityID, ftr.MaxResults, ftr.StartTime, ftr.EndTime, "")
	if err != nil {
		return nil, fmt.Errorf("error invoking twitter api: %w", err)
	}

	result := []Tweet{}
	for _, tweet := range resp.Data {
		result = append(result, Tweet{
			ID:        tweet.ID,
			Text:      tweet.Text,
			CreatedAt: tweet.CreatedAt,
		})
	}

	nextPageToken := resp.Meta.NextToken
	for {
		if nextPageToken == "" {
			// reached end of results
			break
		}

		resp, err := client.invokeFetchTweets(ctx, log, ftr.EntityID, ftr.MaxResults, ftr.StartTime, ftr.EndTime, nextPageToken)
		if err != nil {
			return nil, fmt.Errorf("error invoking twitter api: %w", err)
		}

		nextPageToken = resp.Meta.NextToken

		for _, tweet := range resp.Data {
			result = append(result, Tweet{
				ID:        tweet.ID,
				Text:      tweet.Text,
				CreatedAt: tweet.CreatedAt,
			})
		}
	}

	log.Info(fmt.Sprintf("Received response from Twitter API with %d tweets", len(result)))

	return result, nil
}

func (c *Client) invokeFetchTweets(ctx context.Context, log logger.Interface, entityID string, maxResults int, start, end time.Time, paginationToken string) (getUserTweetsResponse, error) {
	baseUrl := fmt.Sprintf(getUsersTweetsUrl, entityID)
	queryParams := []string{
		fmt.Sprintf("max_results=%d", maxResults),
		"tweet.fields=id,text,created_at",
		fmt.Sprintf("start_time=%s", start.Format(time.RFC3339)),
		fmt.Sprintf("end_time=%s", end.Format(time.RFC3339)),
	}

	if paginationToken != "" {
		queryParams = append(queryParams, fmt.Sprintf("pagination_token=%s", paginationToken))
	}

	url := fmt.Sprintf("%s?%s", baseUrl, strings.Join(queryParams, "&"))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return getUserTweetsResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	log.Info(fmt.Sprintf("Calling Twitter API with: %s", url))
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return getUserTweetsResponse{}, fmt.Errorf("error doing request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getUserTweetsResponse{}, fmt.Errorf("request failed with: %d", resp.StatusCode)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return getUserTweetsResponse{}, fmt.Errorf("error reading response: %w", err)
	}

	var twitterResponse getUserTweetsResponse
	err = json.Unmarshal(response, &twitterResponse)
	if err != nil {
		return getUserTweetsResponse{}, err
	}

	return twitterResponse, nil
}
