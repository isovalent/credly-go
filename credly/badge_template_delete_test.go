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
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteBadgeTemplate(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to delete
	templateID := "template-123"

	// Set up the mock expectations with a 204 No Content response
	mockResponse := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}

	// Verify the request is correct
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "DELETE" &&
			req.URL.String() == "https://api.credly.com/v1/organizations/test-org-id/badge_templates/template-123"
	})).Return(mockResponse, nil)

	// Call the function under test
	err := client.DeleteBadgeTemplate(templateID)

	// Assertions
	assert.NoError(t, err)

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestDeleteBadgeTemplate_Error(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to delete
	templateID := "template-123"

	// Mock error response
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Invalid request"}`)),
	}

	// Set up mock expectations
	mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Call the function under test
	err := client.DeleteBadgeTemplate(templateID)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status code: 400")

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}
