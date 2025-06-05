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
	Id                      string                   `json:"id,omitempty"`
	Name                    string                   `json:"name"`
	Description             string                   `json:"description"`
	Skills                  []string                 `json:"skills"`
	Url                     string                   `json:"url"`
	ImageUrl                string                   `json:"image_url"`
	VanitySlug              string                   `json:"vanity_slug"`
	GlobalActivityURL       string                   `json:"global_activity_url"`
	State                   string                   `json:"state"` // active, archived, or draft
	AllowDuplicateBadges    bool                     `json:"allow_duplicate_badges"`
	BadgesCount             int                      `json:"badges_count"`
	Public                  bool                     `json:"public"`
	AllowDelete             bool                     `json:"allow_delete"`
	AllowArchive            bool                     `json:"allow_archive"`
	Cost                    string                   `json:"cost"`          // Free, Paid
	Level                   string                   `json:"level"`         // Foundational, Intermediate, Advanced
	TimeToEarn              string                   `json:"time_to_earn"`  // Hours, Days, Weeks, Months, Years
	TypeCategory            string                   `json:"type_category"` // Experience, Learning, Validation, Certification
	LockBadgeState          bool                     `json:"lock_badge_state"`
	RecipientType           string                   `json:"recipient_type"`
	ShowSkillTagLinks       bool                     `json:"show_skill_tag_links"`
	VisitlyPrintingDisabled bool                     `json:"printing_disabled"`
	Visibility              string                   `json:"visibility"` // public, private
	VariantsAllowed         bool                     `json:"variants_allowed"`
	VariantType             string                   `json:"variant_type"`
	ReportingTags           []string                 `json:"reporting_tags"`
	Alignments              []BadgeTemplateAlignment `json:"alignments"`
	Recommendations         []BadgeRecommendation    `json:"recommendations"`
	Activities              []BadgeTemplateActivity  `json:"badge_template_activities"`
	LinkedInShareMsg        string                   `json:"linkedin_share_default_message"`
	Translatable            bool                     `json:"translatable"`
	CreatedAt               string                   `json:"created_at"`
	UpdatedAt               string                   `json:"updated_at"`
	StateUpdatedAt          string                   `json:"state_updated_at"`
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
