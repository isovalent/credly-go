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
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArchiveBadgeTemplate(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to archive
	templateID := "template-123"

	// Mock response for archive
	responseTemplate := BadgeTemplate{
		Id:          "template-123",
		Name:        "Archived Badge",
		Description: "Badge that has been archived",
		State:       "archived",
		Skills:      []string{"Skill1", "Skill2", "Skill3"},
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"data": responseTemplate,
	})

	// Set up the mock expectations
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	// Verify the request is correct
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "PUT" &&
			req.URL.String() == "https://api.credly.com/v1/organizations/test-org-id/badge_templates/template-123/archive"
	})).Return(mockResponse, nil)

	// Call the function under test
	archivedTemplate, err := client.ArchiveBadgeTemplate(templateID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "template-123", archivedTemplate.Id)
	assert.Equal(t, "Archived Badge", archivedTemplate.Name)
	assert.Equal(t, "archived", archivedTemplate.State)

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestUnarchiveBadgeTemplate(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to unarchive
	templateID := "template-123"

	// Mock response for unarchive
	responseTemplate := BadgeTemplate{
		Id:          "template-123",
		Name:        "Unarchived Badge",
		Description: "Badge that has been unarchived",
		State:       "active",
		Skills:      []string{"Skill1", "Skill2", "Skill3"},
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"data": responseTemplate,
	})

	// Set up the mock expectations
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	// Verify the request is correct
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "PUT" &&
			req.URL.String() == "https://api.credly.com/v1/organizations/test-org-id/badge_templates/template-123/unarchive"
	})).Return(mockResponse, nil)

	// Call the function under test
	unarchivedTemplate, err := client.UnarchiveBadgeTemplate(templateID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "template-123", unarchivedTemplate.Id)
	assert.Equal(t, "Unarchived Badge", unarchivedTemplate.Name)
	assert.Equal(t, "active", unarchivedTemplate.State)

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestArchiveBadgeTemplate_Error(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to archive
	templateID := "template-123"

	// Mock error response
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Invalid request"}`)),
	}

	// Set up mock expectations
	mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Call the function under test
	_, err := client.ArchiveBadgeTemplate(templateID)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status code: 400")

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestUnarchiveBadgeTemplate_Error(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to unarchive
	templateID := "template-123"

	// Mock error response
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Invalid request"}`)),
	}

	// Set up mock expectations
	mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Call the function under test
	_, err := client.UnarchiveBadgeTemplate(templateID)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status code: 400")

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}
