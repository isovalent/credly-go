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

func TestCreateBadgeTemplate(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Test badge template creation data
	enableEarnThisBadge := true
	templateToCreate := &CreateBadgeTemplateParams{
		Name:                "Test Badge",
		Description:         "Test Badge Description",
		ImageURL:            "https://example.com/image.png",
		Skills:              []string{"Skill1", "Skill2", "Skill3"},
		GlobalActivityURL:   "https://example.com/criteria",
		EarnThisBadgeURL:    "https://example.com/enroll",
		EnableEarnThisBadge: &enableEarnThisBadge,
	}

	// Mock response
	responseTemplate := BadgeTemplate{
		Id:                "new-badge-template-id",
		Name:              "Test Badge",
		Description:       "Test Badge Description",
		ImageUrl:          "https://example.com/image.png",
		Skills:            []string{"Skill1", "Skill2", "Skill3"},
		GlobalActivityURL: "https://example.com/criteria",
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"data": responseTemplate,
	})

	// Set up the mock expectations
	mockResponse := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	// Verify that the right parameters are sent to the API
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		// Verify the HTTP method and URL
		if req.Method != "POST" {
			return false
		}
		if req.URL.String() != "https://api.credly.com/v1/organizations/test-org-id/badge_templates" {
			return false
		}

		// Parse and verify the request body
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body)) // Reset body for later use

		var requestBody map[string]interface{}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			return false
		}

		// Check essential fields
		if requestBody["name"] != "Test Badge" {
			return false
		}
		if requestBody["description"] != "Test Badge Description" {
			return false
		}
		if requestBody["global_activity_url"] != "https://example.com/criteria" {
			return false
		}
		if requestBody["earn_this_badge_url"] != "https://example.com/enroll" {
			return false
		}
		if requestBody["enable_earn_this_badge"] != true {
			return false
		}

		// Check image
		image, ok := requestBody["image"].(map[string]interface{})
		if !ok || image["remote_upload_url"] != "https://example.com/image.png" {
			return false
		}

		// Check skills array
		skills, ok := requestBody["skills"].([]interface{})
		if !ok || len(skills) != 3 {
			return false
		}

		return true
	})).Return(mockResponse, nil)

	// Call the function under test
	createdTemplate, err := client.CreateBadgeTemplate(templateToCreate)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "new-badge-template-id", createdTemplate.Id)
	assert.Equal(t, "Test Badge", createdTemplate.Name)
	assert.Equal(t, []string{"Skill1", "Skill2", "Skill3"}, createdTemplate.Skills)

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestCreateBadgeTemplate_Error(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Test badge template creation data
	templateToCreate := &CreateBadgeTemplateParams{
		Name:              "Test Badge",
		Description:       "Test Badge Description",
		ImageURL:          "https://example.com/image.png",
		Skills:            []string{"Skill1", "Skill2", "Skill3"},
		GlobalActivityURL: "https://example.com/criteria",
	}

	// Mock error response
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Invalid request"}`)),
	}

	// Set up mock expectations
	mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Call the function under test
	_, err := client.CreateBadgeTemplate(templateToCreate)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status code: 400")

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestCreateBadgeTemplate_ValidationErrors(t *testing.T) {
	client := &Client{
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	testCases := []struct {
		name        string
		template    *CreateBadgeTemplateParams
		expectedErr string
	}{
		{
			name:        "Missing name",
			template:    &CreateBadgeTemplateParams{Description: "Description", Skills: []string{"Skill1", "Skill2", "Skill3"}, ImageURL: "https://example.com/image.png", GlobalActivityURL: "https://example.com/criteria"},
			expectedErr: "name is required",
		},
		{
			name:        "Missing description",
			template:    &CreateBadgeTemplateParams{Name: "Name", Skills: []string{"Skill1", "Skill2", "Skill3"}, ImageURL: "https://example.com/image.png", GlobalActivityURL: "https://example.com/criteria"},
			expectedErr: "description is required",
		},
		{
			name:        "Missing image URL",
			template:    &CreateBadgeTemplateParams{Name: "Name", Description: "Description", Skills: []string{"Skill1", "Skill2", "Skill3"}, GlobalActivityURL: "https://example.com/criteria"},
			expectedErr: "image URL is required",
		},
		{
			name:        "Missing skills",
			template:    &CreateBadgeTemplateParams{Name: "Name", Description: "Description", ImageURL: "https://example.com/image.png", GlobalActivityURL: "https://example.com/criteria"},
			expectedErr: "at least three skills are required",
		},
		{
			name:        "Not enough skills",
			template:    &CreateBadgeTemplateParams{Name: "Name", Description: "Description", ImageURL: "https://example.com/image.png", Skills: []string{"Skill1", "Skill2"}, GlobalActivityURL: "https://example.com/criteria"},
			expectedErr: "at least three skills are required",
		},
		{
			name:        "Missing global activity URL (criteria)",
			template:    &CreateBadgeTemplateParams{Name: "Name", Description: "Description", ImageURL: "https://example.com/image.png", Skills: []string{"Skill1", "Skill2", "Skill3"}},
			expectedErr: "global activity URL (criteria) is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.CreateBadgeTemplate(tc.template)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}
