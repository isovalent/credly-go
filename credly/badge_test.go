// Copyright 2024 Cisco Systems, Inc. and its affiliates

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package credly

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIssueBadge(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	templateId := "template-123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	expectedBadge := BadgeInfo{
		Id:       "badge-123",
		ImageUrl: "http://image.url",
		Url:      "http://badge.url",
		IssuedAt: time.Now(),
		State:    "issued",
	}

	responseBody, _ := json.Marshal(issueBadgeResponse{
		Data: expectedBadge,
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	badge, err := client.IssueBadge(templateId, email, firstName, lastName)

	assert.NoError(t, err)
	assert.Equal(t, expectedBadge.Id, badge.Id)
	assert.Equal(t, expectedBadge.ImageUrl, badge.ImageUrl)
	assert.Equal(t, expectedBadge.Url, badge.Url)
	assert.Equal(t, expectedBadge.State, badge.State)
	// Optionally compare IssuedAt with some tolerance
	// assert.WithinDuration(t, expectedBadge.IssuedAt, badge.IssuedAt, time.Second)
	mockClient.AssertExpectations(t)
}

func TestIssueBadge_BadgeAlreadyIssued(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	templateId := "template-123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Simulate a 422 response indicating the badge is already issued
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil)

	badge, err := client.IssueBadge(templateId, email, firstName, lastName)

	assert.Error(t, err)
	assert.Equal(t, ErrBadgeAlreadyIssued, err.Error())
	assert.Empty(t, badge)
	mockClient.AssertExpectations(t)
}

func TestIssueBadge_Failure(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	templateId := "template-123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Simulate a failure response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil)

	badge, err := client.IssueBadge(templateId, email, firstName, lastName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed")
	assert.Empty(t, badge)
	mockClient.AssertExpectations(t)
}

func TestGetBadges_NoCollections(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	email := "test@example.com"

	expectedBadges := []BadgeInfo{
		{Id: "badge-123", ImageUrl: "http://image.url", Url: "http://badge.url"},
		{Id: "badge-456", ImageUrl: "http://image2.url", Url: "http://badge2.url"},
	}

	responseBody, _ := json.Marshal(getBadgesResponse{
		Data: expectedBadges,
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	badges, err := client.GetBadges(email, []string{})

	assert.NoError(t, err)
	assert.Equal(t, expectedBadges, badges)
}

func TestGetBadges_WithCollections(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{HTTPClient: mockClient}

	email := "test@example.com"
	collections := []string{"collection1", "collection2"}

	expectedBadges := []BadgeInfo{
		{Id: "badge-123", ImageUrl: "http://image.url", Url: "http://badge.url"},
		{Id: "badge-456", ImageUrl: "http://image2.url", Url: "http://badge2.url"},
	}

	responseBody, _ := json.Marshal(getBadgesResponse{
		Data: expectedBadges,
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	badges, err := client.GetBadges(email, collections)

	assert.NoError(t, err)
	assert.Equal(t, expectedBadges, badges)
}

func TestGetBadge(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	email := "test@example.com"
	badgeId := "badge-123"

	expectedBadge := BadgeInfo{
		Id: "badge-123",
	}

	responseBody, _ := json.Marshal(getBadgesResponse{
		Data: []BadgeInfo{expectedBadge},
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	badge, err := client.GetBadge(email, badgeId)

	assert.NoError(t, err)
	assert.Equal(t, expectedBadge, badge)
	mockClient.AssertExpectations(t)
}


func TestGetBadges_Failure(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	email := "test@example.com"

	// Simulate a failure response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil)

	badges, err := client.GetBadges(email, []string{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed")
	assert.Empty(t, badges)
	mockClient.AssertExpectations(t)
}
