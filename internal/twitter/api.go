package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kordape/ottct-poller-service/internal/processor"
)

// testing path
// http://localhost:8080/v1/tweets/classify?userId=1277254376&maxResults=90&startTime=2022-01-12&endTime=2022-06-15

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

// Make sure Client implement TweetsFetcher interface
var _ processor.TweetsFetcher = &Client{}

type Client struct {
	httpClient  *http.Client
	bearerToken string
}

func New(client *http.Client, bearerToken string) *Client {
	return &Client{
		bearerToken: bearerToken,
		httpClient:  client,
	}
}

type getUserTweetsResponse struct {
	Data []tweet  `json:"data"`
	Meta metadata `json:"meta"`
}

type tweet struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Text      string `json:"text"`
}

// metadata left to enable pagination option in perspective
// can be removed if needed
type metadata struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

func (client *Client) FetchTweets(ctx context.Context, ftr processor.FetchTweetsRequest) (processor.FetchTweetsResponse, error) {
	baseUrl := fmt.Sprintf(getUsersTweetsUrl, ftr.EntityID)
	queryParams := []string{
		fmt.Sprintf("max_results=%d", ftr.MaxResults),
		"tweet.fields=id,text,created_at",
		fmt.Sprintf("start_time=%s", ftr.StartTime.Format(time.RFC3339)),
		fmt.Sprintf("end_time=%s", ftr.EndTime.Format(time.RFC3339)),
	}

	url := fmt.Sprintf("%s?%s", baseUrl, strings.Join(queryParams, "&"))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.bearerToken))
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with: %d", resp.StatusCode)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var tweeterResponse getUserTweetsResponse
	err = json.Unmarshal(response, &tweeterResponse)
	if err != nil {
		return nil, err
	}

	result := make([]processor.Tweet, len(tweeterResponse.Data))
	for i, tweet := range tweeterResponse.Data {
		result[i] = processor.Tweet{
			ID:        tweet.ID,
			Text:      tweet.Text,
			CreatedAt: tweet.CreatedAt,
		}
	}

	return result, nil
}
