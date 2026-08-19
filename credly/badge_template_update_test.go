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

func TestUpdateBadgeTemplate(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to update
	templateID := "template-123"

	// Parameters for updating the badge template
	enableEarnThisBadge := false
	updateParams := &UpdateBadgeTemplateParams{
		Name:                "Updated Badge",
		Description:         "Updated badge template description",
		ImageURL:            "https://example.com/updated-image.png",
		Skills:              []string{"UpdatedSkill1", "UpdatedSkill2", "UpdatedSkill3"},
		GlobalActivityURL:   "https://example.com/updated-criteria",
		EarnThisBadgeURL:    "https://example.com/updated-enroll",
		EnableEarnThisBadge: &enableEarnThisBadge,
		ReportingTags:       []string{"updated", "tag", "test"},
	}

	// Mock response for update
	responseTemplate := BadgeTemplate{
		Id:                "template-123",
		Name:              "Updated Badge",
		Description:       "Updated badge template description",
		Skills:            []string{"UpdatedSkill1", "UpdatedSkill2", "UpdatedSkill3"},
		GlobalActivityURL: "https://example.com/updated-criteria",
		State:             "active",
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"data": responseTemplate,
	})

	// Set up the mock expectations
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	// Verify that the right parameters are sent to the API
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		// Verify the HTTP method and URL
		if req.Method != "PUT" {
			return false
		}
		if req.URL.String() != "https://api.credly.com/v1/organizations/test-org-id/badge_templates/template-123" {
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
		if requestBody["name"] != "Updated Badge" {
			return false
		}
		if requestBody["description"] != "Updated badge template description" {
			return false
		}
		if requestBody["global_activity_url"] != "https://example.com/updated-criteria" {
			return false
		}
		if requestBody["earn_this_badge_url"] != "https://example.com/updated-enroll" {
			return false
		}
		if requestBody["enable_earn_this_badge"] != false {
			return false
		}

		// Check image
		image, ok := requestBody["image"].(map[string]interface{})
		if !ok || image["remote_upload_url"] != "https://example.com/updated-image.png" {
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
	updatedTemplate, err := client.UpdateBadgeTemplate(templateID, updateParams)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "template-123", updatedTemplate.Id)
	assert.Equal(t, "Updated Badge", updatedTemplate.Name)
	assert.Equal(t, "Updated badge template description", updatedTemplate.Description)
	assert.Equal(t, []string{"UpdatedSkill1", "UpdatedSkill2", "UpdatedSkill3"}, updatedTemplate.Skills)

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestUpdateBadgeTemplate_Error(t *testing.T) {
	// Mock HTTP client
	mockHTTPClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient:     mockHTTPClient,
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to update
	templateID := "template-123"

	// Parameters for updating the badge template
	updateParams := &UpdateBadgeTemplateParams{
		Name:              "Updated Badge",
		Description:       "Updated badge template description",
		Skills:            []string{"UpdatedSkill1", "UpdatedSkill2", "UpdatedSkill3"},
		GlobalActivityURL: "https://example.com/updated-criteria",
	}

	// Mock error response
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error":"Invalid request"}`)),
	}

	// Set up mock expectations
	mockHTTPClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Call the function under test
	_, err := client.UpdateBadgeTemplate(templateID, updateParams)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status code: 400")

	// Verify all expectations were met
	mockHTTPClient.AssertExpectations(t)
}

func TestUpdateBadgeTemplate_ValidationErrors(t *testing.T) {
	client := &Client{
		OrganizationId: "test-org-id",
		authToken:      "test-auth-token",
	}

	// Template ID to update
	templateID := "template-123"

	testCases := []struct {
		name        string
		params      *UpdateBadgeTemplateParams
		expectedErr string
	}{
		{
			name:        "Missing name",
			params:      &UpdateBadgeTemplateParams{Description: "Description", Skills: []string{"Skill1", "Skill2", "Skill3"}},
			expectedErr: "name is required",
		},
		{
			name:        "Missing description",
			params:      &UpdateBadgeTemplateParams{Name: "Name", Skills: []string{"Skill1", "Skill2", "Skill3"}},
			expectedErr: "description is required",
		},
		{
			name:        "Missing skills",
			params:      &UpdateBadgeTemplateParams{Name: "Name", Description: "Description"},
			expectedErr: "at least three skills are required",
		},
		{
			name:        "Not enough skills",
			params:      &UpdateBadgeTemplateParams{Name: "Name", Description: "Description", Skills: []string{"Skill1", "Skill2"}},
			expectedErr: "at least three skills are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.UpdateBadgeTemplate(templateID, tc.params)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}
