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

type Client struct {
	httpClient  *http.Client
	bearerToken string
}

func New(bearerToken string) *Client {
	return &Client{
		bearerToken: bearerToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
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
	var queryParams []string
	queryParams = append(queryParams, fmt.Sprintf("max_results=%d", ftr.MaxResults))
	queryParams = append(queryParams, "tweet.fields=id,text,created_at")
	if ftr.StartTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("start_time=%s", ftr.StartTime))
	}
	if ftr.EndTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("end_time=%s", ftr.EndTime))
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
