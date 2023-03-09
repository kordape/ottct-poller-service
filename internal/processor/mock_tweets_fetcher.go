// Code generated by mockery. DO NOT EDIT.

package processor

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockTweetsFetcher is an autogenerated mock type for the TweetsFetcher type
type MockTweetsFetcher struct {
	mock.Mock
}

// FetchTweets provides a mock function with given fields: _a0, _a1
func (_m *MockTweetsFetcher) FetchTweets(_a0 context.Context, _a1 FetchTweetsRequest) (FetchTweetsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 FetchTweetsResponse
	if rf, ok := ret.Get(0).(func(context.Context, FetchTweetsRequest) FetchTweetsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(FetchTweetsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, FetchTweetsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewMockTweetsFetcherT interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTweetsFetcher creates a new instance of MockTweetsFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTweetsFetcher(t NewMockTweetsFetcherT) *MockTweetsFetcher {
	mock := &MockTweetsFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}