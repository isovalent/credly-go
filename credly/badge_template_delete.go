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
	"fmt"
	"net/http"
)

// DeleteBadgeTemplate deletes a badge template.
//
// templateId: The ID of the badge template to be deleted.
// Returns: An error if the operation fails, nil otherwise.
func (c *Client) DeleteBadgeTemplate(templateId string) error {
	url := fmt.Sprintf("https://api.credly.com/v1/organizations/%s/badge_templates/%s", c.OrganizationId, templateId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("[credly.DeleteBadgeTemplate] Failed to create request: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("[credly.DeleteBadgeTemplate] Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[credly.DeleteBadgeTemplate] API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
