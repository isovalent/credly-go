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
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBadgeTemplate(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	templateId := "template-123"

	expectedTemplate := BadgeTemplate{
		Id:       "template-123",
		Name:     "Test Badge",
		ImageUrl: "http://image.url",
	}

	responseBody, _ := json.Marshal(getBadgeTemplateResponse{
		Data: expectedTemplate,
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	template, err := client.GetBadgeTemplate(templateId)

	assert.NoError(t, err)
	assert.Equal(t, expectedTemplate, template)
	mockClient.AssertExpectations(t)
}

func TestGetBadgeTemplates(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	expectedTemplates := []BadgeTemplate{
		{Id: "template-123", Name: "Badge 1", ImageUrl: "http://image1.url"},
		{Id: "template-456", Name: "Badge 2", ImageUrl: "http://image2.url"},
	}

	responseBody, _ := json.Marshal(getBadgeTemplatesResponse{
		Data: expectedTemplates,
	})

	// Simulate a successful response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	templates, err := client.GetBadgeTemplates()

	assert.NoError(t, err)
	assert.Equal(t, expectedTemplates, templates)
	mockClient.AssertExpectations(t)
}

func TestGetBadgeTemplate_Failure(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{
		HTTPClient: mockClient,
		authToken:  base64.StdEncoding.EncodeToString([]byte("test-token" + "|")),
	}

	templateId := "template-123"

	// Simulate a failure response
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil)

	template, err := client.GetBadgeTemplate(templateId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed")
	assert.Empty(t, template)
	mockClient.AssertExpectations(t)
}
