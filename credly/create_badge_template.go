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
	"net/http"
)

// CreateBadgeTemplateParams represents the parameters for creating a new badge template.
type CreateBadgeTemplateParams struct {
	// Name is the name of the badge template (required)
	Name string `json:"name"`

	// Description is the description of the badge template (required)
	Description string `json:"description"`

	// ImageURL is the URL to the image to use for the badge template (required)
	// Will be sent as remote_upload_url inside an image object
	ImageURL string `json:"-"`

	// Skills is a list of skills associated with the badge (requires at least 3)
	Skills []string `json:"skills"`

	// GlobalActivityURL is the criteria URL for the badge (required)
	GlobalActivityURL string `json:"global_activity_url"`

	// Optional fields below
	Cost              string                  `json:"cost,omitempty"`              // "Free" or "Paid"
	Level             string                  `json:"level,omitempty"`             // "Foundational", "Intermediate", or "Advanced"
	TimeToEarn        string                  `json:"time_to_earn,omitempty"`      // "Hours", "Days", "Weeks", "Months", or "Years"
	TypeCategory      string                  `json:"type_category,omitempty"`     // "Experience", "Learning", "Validation", or "Certification"
	Translatable      *bool                   `json:"translatable,omitempty"`
	ReportingTags     []string                `json:"reporting_tags,omitempty"`
	Alignments        []BadgeTemplateAlignment `json:"alignments,omitempty"`
	Recommendations   []BadgeRecommendation   `json:"recommendations,omitempty"`
	LockBadgeState    *bool                   `json:"lock_badge_state,omitempty"`
	Activities        []BadgeTemplateActivity `json:"badge_template_activities,omitempty"`
}

// BadgeTemplateAlignment represents an educational standard alignment
type BadgeTemplateAlignment struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// BadgeRecommendation represents a recommendation for a badge template
type BadgeRecommendation struct {
	// Three possible types of recommendations:
	// 1. By ID (existing recommendation)
	ID string `json:"id,omitempty"`

	// 2. By Badge Template ID (recommending another badge)
	RecommendedBadgeTemplateID string `json:"recommended_badge_template_id,omitempty"`

	// 3. By Type (recommending a URL)
	Type        string `json:"type,omitempty"`        // badge, information, education, employment, participation, offer
	ActivityURL string `json:"activity_url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// BadgeTemplateActivity represents an activity related to a badge template
type BadgeTemplateActivity struct {
	// ActivityType is the type of activity
	// Valid values: "Award", "Schedule / Registration", "Badge", etc.
	ActivityType string `json:"activity_type"`

	// Title is the title of the activity (optional for Badge type)
	Title string `json:"title,omitempty"`

	// ActivityURL is the URL to the activity (optional for Badge type)
	ActivityURL string `json:"activity_url,omitempty"`

	// RequiredBadgeTemplateID is the ID of a required badge template (for Badge type only)
	RequiredBadgeTemplateID string `json:"required_badge_template_id,omitempty"`
}

// CreateBadgeTemplate creates a new badge template with the provided parameters.
//
// params: The parameters for creating the badge template.
// Returns: The created badge template, or an error if the operation fails.
func (c *Client) CreateBadgeTemplate(params *CreateBadgeTemplateParams) (b BadgeTemplate, err error) {
	// Validate required fields
	if params.Name == "" {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] name is required")
	}
	if params.Description == "" {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] description is required")
	}
	if params.ImageURL == "" {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] image URL is required")
	}
	if len(params.Skills) < 3 {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] at least three skills are required")
	}
	if params.GlobalActivityURL == "" {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] global activity URL (criteria) is required")
	}

	// Construct the API URL
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates", c.OrganizationId)

	// Create a copy of params for JSON marshalling
	requestData := make(map[string]interface{})
	
	// Marshal the struct into a map and set all fields except ImageURL
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to marshal parameters: %v", err)
	}
	
	if err := json.Unmarshal(paramsJSON, &requestData); err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to prepare request data: %v", err)
	}
	
	// Add image data in the format expected by the API
	requestData["image"] = map[string]string{
		"remote_upload_url": params.ImageURL,
	}
	
	// Marshal the request data into JSON
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to marshal request data: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to create request: %v", err)
	}

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data BadgeTemplate `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return b, fmt.Errorf("[credly.CreateBadgeTemplate] Failed to parse response: %v", err)
	}

	return response.Data, nil
}
