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
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock of the http.Client used for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewClient(t *testing.T) {
	token := "test-token"
	orgId := "abcd-efgh-1234-5678"
	expectedToken := base64.StdEncoding.EncodeToString([]byte(token + "|"))

	client := NewClient(token, orgId)

	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, expectedToken, client.authToken)
}

func TestDo(t *testing.T) {
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockHTTPClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	req, err := http.NewRequest("GET", "https://api.credly.com/v1/some-endpoint", nil)
	assert.NoError(t, err)

	expectedResponse := &http.Response{
		StatusCode: 200,
	}

	mockHTTPClient.On("Do", mock.Anything).Return(expectedResponse, nil)

	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)

	// Check that the correct headers are set
	assert.Equal(t, "Basic "+client.authToken, req.Header.Get("Authorization"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))

	mockHTTPClient.AssertExpectations(t)
}
