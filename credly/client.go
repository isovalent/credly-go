// Package credly provides a client to interact with the Credly API,
// allowing operations such as issuing and retrieving badges for a specific organization.
package credly

import (
	"encoding/base64"
	"net/http"
)

// HTTPClientInterface defines the methods that http.Client and MockHTTPClient must implement.
// This interface allows for mocking and testing of HTTP requests.
type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a Credly API client used to interact with the organization's badges.
// It handles authentication and provides methods for making authorized HTTP requests.
type Client struct {
	// HTTPClient is the HTTP client used to make requests to the Credly API.
	HTTPClient HTTPClientInterface

	// authToken is the base64-encoded authentication token for API access.
	authToken string

	// OrganizationId is the unique identifier for the organization in Credly.
	OrganizationId string
}

// ErrBadgeAlreadyIssued indicates that a badge has already been issued to the user.
const ErrBadgeAlreadyIssued = "User already has this badge"

// NewClient creates a new instance of the Credly API client.
// It accepts an API token and the organization ID, returning a Client
// with an encoded authentication token and organization-specific settings.
//
// token: The API token provided by Credly for authentication.
// organizationId: The unique identifier for the organization in Credly.
// Returns: A new Client instance configured for Credly API interaction.
func NewClient(token, organizationId string) *Client {
	// Encode the token with base64 and append a separator "|"
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token + "|"))

	return &Client{
		HTTPClient:     &http.Client{},
		authToken:      encodedToken,
		OrganizationId: organizationId,
	}
}

// Do sends an HTTP request using the Client's HTTP client, adding the necessary
// authentication headers for the Credly API.
//
// req: The HTTP request to be sent.
// Returns: The HTTP response and any error encountered.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Add the required headers for Credly API authentication and content type.
	req.Header.Set("Authorization", "Basic "+c.authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute the HTTP request using the client's HTTP client.
	return c.HTTPClient.Do(req)
}
