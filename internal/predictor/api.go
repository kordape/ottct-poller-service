package predictor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kordape/ottct-poller-service/internal/processor"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(client *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: client,
		baseURL:    baseURL,
	}
}

type request struct {
	Tweet string `json:"tweet"`
}
type response struct {
	Prediction []int `json:"prediction"`
}

func (c *Client) Classify(ctx context.Context, requests []processor.ClassifyTweetsRequest) (processor.ClassifyTweetsResponse, error) {
	predictRequest := make([]request, len(requests))
	for i, r := range requests {
		predictRequest[i] = request{
			Tweet: r.Tweet,
		}
	}

	buf, err := json.Marshal(predictRequest)
	if err != nil {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("error marshalling request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("error doing http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("request failed with: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("error reading response: %w", err)
	}

	var predictions response
	err = json.Unmarshal(body, &predictions)
	if err != nil {
		return processor.ClassifyTweetsResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result := processor.ClassifyTweetsResponse{}
	classifications := make([]processor.Classification, len(predictions.Prediction))
	for i, p := range predictions.Prediction {
		classifications[i] = processor.Classification(p)
	}
	result.Classification = classifications

	return result, nil
}
