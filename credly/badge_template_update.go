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
	"fmt"
	"io"
	"net/http"
)

// UpdateBadgeTemplateParams represents the parameters for updating a badge template.
type UpdateBadgeTemplateParams struct {
	// Name is the name of the badge template (required)
	Name string `json:"name"`

	// Description is the description of the badge template (required)
	Description string `json:"description"`

	// ImageURL is the URL to the image to use for the badge template (optional for updates)
	// Will be sent as remote_upload_url inside an image object if provided
	ImageURL string `json:"-"`

	// Skills is a list of skills associated with the badge (requires at least 3)
	Skills []string `json:"skills"`

	// GlobalActivityURL is the criteria URL for the badge (optional for updates)
	GlobalActivityURL string `json:"global_activity_url,omitempty"`

	// EarnThisBadgeURL is the public URL where users can register or enroll to earn the badge.
	EarnThisBadgeURL string `json:"earn_this_badge_url,omitempty"`

	// EnableEarnThisBadge controls whether the "Earn this Badge" button is displayed.
	// A pointer distinguishes an explicit false value from an omitted value.
	EnableEarnThisBadge *bool `json:"enable_earn_this_badge,omitempty"`

	// Optional fields below
	Cost             string                   `json:"cost,omitempty"`          // "Free" or "Paid"
	Level            string                   `json:"level,omitempty"`         // "Foundational", "Intermediate", or "Advanced"
	TimeToEarn       string                   `json:"time_to_earn,omitempty"`  // "Hours", "Days", "Weeks", "Months", or "Years"
	TypeCategory     string                   `json:"type_category,omitempty"` // "Experience", "Learning", "Validation", or "Certification"
	Translatable     *bool                    `json:"translatable,omitempty"`
	ReportingTags    []string                 `json:"reporting_tags,omitempty"`
	Alignments       []BadgeTemplateAlignment `json:"alignments,omitempty"`
	Recommendations  []BadgeRecommendation    `json:"recommendations,omitempty"`
	LockBadgeState   *bool                    `json:"lock_badge_state,omitempty"`
	Activities       []BadgeTemplateActivity  `json:"badge_template_activities,omitempty"`
	LinkedInShareMsg string                   `json:"linkedin_share_default_message,omitempty"`
}

// UpdateBadgeTemplate updates an existing badge template with the provided parameters.
//
// templateId: The ID of the badge template to update.
// params: The parameters for updating the badge template.
// Returns: The updated badge template, or an error if the operation fails.
func (c *Client) UpdateBadgeTemplate(templateId string, params *UpdateBadgeTemplateParams) (b BadgeTemplate, err error) {
	// Validate required fields
	if params.Name == "" {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] name is required")
	}
	if params.Description == "" {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] description is required")
	}
	if len(params.Skills) < 3 {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] at least three skills are required")
	}

	// Construct the API URL
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates/%s", c.OrganizationId, templateId)

	// Create a copy of params for JSON marshalling
	requestData := make(map[string]interface{})

	// Marshal the struct into a map and set all fields except ImageURL
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to marshal parameters: %v", err)
	}

	if err := json.Unmarshal(paramsJSON, &requestData); err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to prepare request data: %v", err)
	}

	// Add image data if provided
	if params.ImageURL != "" {
		requestData["image"] = map[string]string{
			"remote_upload_url": params.ImageURL,
		}
	}

	// Marshal the request data into JSON
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to marshal request data: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to create request: %v", err)
	}

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Parse body for more details if needed
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return b, fmt.Errorf("[credly.UpdateBadgeTemplate] API request failed with status code: %d, response: %s", resp.StatusCode, body)
		}
		// If no body, just return the status code
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data BadgeTemplate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return b, fmt.Errorf("[credly.UpdateBadgeTemplate] Failed to parse response: %v", err)
	}

	return response.Data, nil
}
