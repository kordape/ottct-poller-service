package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	log.Info(fmt.Sprintf("Calling Twitter API with: %s", url))
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

	var twitterResponse getUserTweetsResponse
	err = json.Unmarshal(response, &twitterResponse)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Received response from Twitter API: %v", twitterResponse))
	result := make([]Tweet, len(twitterResponse.Data))
	for i, tweet := range twitterResponse.Data {
		result[i] = Tweet{
			ID:        tweet.ID,
			Text:      tweet.Text,
			CreatedAt: tweet.CreatedAt,
		}
	}

	return result, nil
}
