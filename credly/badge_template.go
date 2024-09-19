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
	"encoding/json"
	"fmt"
	"net/http"
)

// getBadgeTemplateResponse represents the response structure when fetching a specific badge template.
type getBadgeTemplateResponse struct {
	Data BadgeTemplate `json:"data"`
}

// getBadgeTemplatesResponse represents the response structure when fetching multiple badge templates.
type getBadgeTemplatesResponse struct {
	Data []BadgeTemplate `json:"data"`
}

// BadgeTemplate represents the details of a badge template in Credly.
type BadgeTemplate struct {
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name"`
	Skills     []string `json:"skills"`
	Url        string   `json:"url"`
	ImageUrl   string   `json:"image_url"`
	VanitySlug string   `json:"vanity_slug"`
}

// GetBadgeTemplate retrieves a specific badge template by its ID.
//
// templateId: The ID of the badge template to be retrieved.
// Returns: A BadgeTemplate representing the retrieved template, or an error if the operation fails.
func (c *Client) GetBadgeTemplate(templateId string) (b BadgeTemplate, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates/%s", c.OrganizationId, templateId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return b, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("[credly.GetBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	var badgeResp getBadgeTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&badgeResp); err != nil {
		return b, fmt.Errorf("[credly.GetBadgeTemplate] Failed to parse JSON data: %v", err)
	}

	return badgeResp.Data, nil
}

// GetBadgeTemplates retrieves all badge templates for the organization.
//
// Returns: A slice of BadgeTemplate representing all templates, or an error if the operation fails.
func (c *Client) GetBadgeTemplates() (b []BadgeTemplate, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates", c.OrganizationId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return b, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("[credly.GetBadgeTemplates] API request failed with status code: %d", resp.StatusCode)
	}

	var badgeResp getBadgeTemplatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&badgeResp); err != nil {
		return b, fmt.Errorf("[credly.GetBadgeTemplates] Failed to parse JSON data: %v", err)
	}

	return badgeResp.Data, nil
}
