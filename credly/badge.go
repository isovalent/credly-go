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
	"net/url"
	"strings"
	"time"
)

// issueBadgeResponse represents the response structure when a badge is issued.
// see https://www.credly.com/docs/issued_badges
type issueBadgeResponse struct {
	Data BadgeInfo `json:"data"`
}

// getBadgesResponse represents the response structure when fetching multiple badges.
type getBadgesResponse struct {
	Data []BadgeInfo `json:"data"`
}


// BadgeInfo represents the details of an issued badge.
type BadgeInfo struct {
	Id       string    `json:"id"`
	ImageUrl string    `json:"image_url"`
	Url      string    `json:"badge_url"`
	IssuedAt time.Time `json:"issued_at"`
	State    string    `json:"state"`

	Image struct {
		Url string `json:"url"`
	} `json:"image"`

	Template BadgeTemplate `json:"badge_template"`

	User struct {
		Id        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Url       string `json:"url"`
	} `json:"user"`
}


// IssueBadge issues a new badge to a user based on their email and personal details.
//
// templateId: The ID of the badge template to be issued.
// email: The recipient's email address.
// firstName: The recipient's first name.
// lastName: The recipient's last name.
// Returns: BadgeInfo representing the issued badge, or an error if the operation fails.
func (c *Client) IssueBadge(templateId, email, firstName, lastName string) (i BadgeInfo, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badges", c.OrganizationId)

	now := time.Now()
	issuedAt := now.Format("2006-01-02 15:04:05 -0700")

	params := map[string]interface{}{
		"badge_template_id":    templateId,
		"recipient_email":      email,
		"issued_to_first_name": firstName,
		"issued_to_last_name":  lastName,
		"issued_at":            issuedAt,
	}
	reqBody, err := json.Marshal(params)
	if err != nil {
		return i, fmt.Errorf("[credly.IssueBadge] Failed to marshal parameters: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return i, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return i, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		// Contact already has badge
		return i, fmt.Errorf(ErrBadgeAlreadyIssued)
	}

	if resp.StatusCode != http.StatusCreated {
		return i, fmt.Errorf("[credly.IssueBadge] API request failed with status code: %d", resp.StatusCode)
	}

	var badgeResp issueBadgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&badgeResp); err != nil {
		return i, fmt.Errorf("[credly.IssueBadge] Failed to parse JSON data: %v", err)
	}

	return badgeResp.Data, nil
}

// GetBadges retrieves all badges for a given email, optionally filtered by collections.
//
// email: The recipient's email address.
// collections: A list of collection tags to filter badges.
// Returns: A slice of BadgeInfo representing the retrieved badges, or an error if the operation fails.
func (c *Client) GetBadges(email string, collections []string) (b []BadgeInfo, err error) {
	qUrl := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badges", c.OrganizationId)
	qUrl = fmt.Sprintf("%s?filter=recipient_email_all::%s", qUrl, url.QueryEscape(email))

	if len(collections) > 0 {
		colFilter := fmt.Sprintf("|badge_templates[reporting_tags]::%s", strings.Join(collections, ","))
		qUrl = fmt.Sprintf("%s%s", qUrl, url.QueryEscape(colFilter))
	}

	req, err := http.NewRequest("GET", qUrl, nil)
	if err != nil {
		return b, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("[credly.GetBadges] API request failed with status code: %d", resp.StatusCode)
	}

	var badgesResp getBadgesResponse
	if err := json.NewDecoder(resp.Body).Decode(&badgesResp); err != nil {
		return b, fmt.Errorf("[credly.GetBadges] Failed to parse JSON data: %v", err)
	}

	return badgesResp.Data, nil
}

// GetBadge retrieves a specific badge for a given email and badge ID.
//
// email: The recipient's email address.
// badgeId: The ID of the badge to be retrieved.
// Returns: A BadgeInfo representing the retrieved badge, or an error if the operation fails.
func (c *Client) GetBadge(email, badgeId string) (b BadgeInfo, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badges", c.OrganizationId)
	url = fmt.Sprintf("%s?filter=recipient_email_all::%s|badge_template_id::%s", url, email, badgeId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return b, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	var badgesResp getBadgesResponse
	if err := json.NewDecoder(resp.Body).Decode(&badgesResp); err != nil {
		return b, fmt.Errorf("Failed to parse JSON data: %v", err)
	}

	if len(badgesResp.Data) == 0 {
		return b, nil
	}

	return badgesResp.Data[0], nil
}
