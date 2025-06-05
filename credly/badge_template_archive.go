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

// ArchiveBadgeTemplate archives a badge template to prevent further badges from being issued.
//
// templateId: The ID of the badge template to archive.
// Returns: The archived BadgeTemplate, or an error if the operation fails.
func (c *Client) ArchiveBadgeTemplate(templateId string) (b BadgeTemplate, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates/%s/archive", c.OrganizationId, templateId)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return b, fmt.Errorf("[credly.ArchiveBadgeTemplate] Failed to create request: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, fmt.Errorf("[credly.ArchiveBadgeTemplate] Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("[credly.ArchiveBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data BadgeTemplate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return b, fmt.Errorf("[credly.ArchiveBadgeTemplate] Failed to parse response: %v", err)
	}

	return response.Data, nil
}

// UnarchiveBadgeTemplate unarchives a previously archived badge template to allow badges to be issued again.
//
// templateId: The ID of the badge template to unarchive.
// Returns: The unarchived BadgeTemplate, or an error if the operation fails.
func (c *Client) UnarchiveBadgeTemplate(templateId string) (b BadgeTemplate, err error) {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates/%s/unarchive", c.OrganizationId, templateId)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return b, fmt.Errorf("[credly.UnarchiveBadgeTemplate] Failed to create request: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return b, fmt.Errorf("[credly.UnarchiveBadgeTemplate] Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, fmt.Errorf("[credly.UnarchiveBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data BadgeTemplate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return b, fmt.Errorf("[credly.UnarchiveBadgeTemplate] Failed to parse response: %v", err)
	}

	return response.Data, nil
}
