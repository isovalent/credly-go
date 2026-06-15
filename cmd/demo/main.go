package main

import (
	"fmt"
	"os"
	"time"

	"github.com/isovalent/credly-go/credly"
)

func main() {
	// Get credentials from environment variables
	token := os.Getenv("CREDLY_TOKEN")
	orgID := os.Getenv("CREDLY_ORG_ID")

	if token == "" || orgID == "" {
		fmt.Println("Error: CREDLY_TOKEN and CREDLY_ORG_ID environment variables must be set")
		os.Exit(1)
	}

	// Initialize the Credly client
	client := credly.NewClient(token, orgID)

	// Create a unique test badge name with timestamp to avoid conflicts
	timestamp := time.Now().Format("20060102-150405")
	badgeName := fmt.Sprintf("Test Badge Template %s", timestamp)

	// Create parameters for the new badge template
	templateParams := &credly.CreateBadgeTemplateParams{
		Name:              badgeName,
		Description:       "This is a test badge template created via the credly-go library.",
		ImageURL:          "https://labs-map.isovalent.com/labs/cilium-getting-started/badge.png", // Placeholder image
		Skills:            []string{"Testing", "Go Programming", "API Integration"},
		GlobalActivityURL: "https://example.com/criteria",
		// Activities
		Activities: []credly.BadgeTemplateActivity{
			{
				ActivityType:            "Award",
				Title:                   "Awarded for completing the test badge template demo",
				ActivityURL:             "https://example.com/award-criteria",
				RequiredBadgeTemplateID: "required-badge-template-id", // Example ID
			},
		},
		// Optional fields
		Level:         "Intermediate",
		TimeToEarn:    "Days",
		TypeCategory:  "Learning",
		ReportingTags: []string{"test", "api"},
	}

	fmt.Printf("Creating badge template '%s'...\n", badgeName)
	template, err := client.CreateBadgeTemplate(templateParams)
	if err != nil {
		fmt.Printf("Error creating badge template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created badge template with ID: %s\n", template.Id)
	fmt.Printf("Badge state: %s\n", template.State)

	// Wait for user to review the created badge template
	fmt.Println("Badge template created successfully. You can review it in your Credly account.")
	fmt.Println("Press Enter to continue with the demo...")
	_, _ = fmt.Scanln()

	// Update the badge template
	fmt.Println("Updating the badge template...")
	updateParams := &credly.UpdateBadgeTemplateParams{
		Name:              badgeName + " (Updated)",
		Description:       "This badge template has been updated via the credly-go library.",
		ImageURL:          "https://labs-map.isovalent.com/labs/cilium-cluster-mesh/badge.png", // Updated image
		Skills:            []string{"Testing", "Go Programming", "API Integration", "Updated Skill"},
		GlobalActivityURL: "https://example.com/updated-criteria",
		// Activities
		Activities: []credly.BadgeTemplateActivity{
			{
				ActivityType:            "Award",
				Title:                   "Awarded for completing the updated test badge template demo",
				ActivityURL:             "https://example.com/updated-award-criteria",
				RequiredBadgeTemplateID: "required-badge-template-id-updated", // Example ID
			},
		},
		// Optional fields
		Level:         "Advanced",
		TimeToEarn:    "Weeks",
		TypeCategory:  "Certification",
		ReportingTags: []string{"test", "api", "updated"},
	}
	updatedTemplate, err := client.UpdateBadgeTemplate(template.Id, updateParams)
	if err != nil {
		fmt.Printf("Error updating badge template: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Badge template updated successfully. New name: %s\n", updatedTemplate.Name)
	fmt.Printf("Badge state after update: %s\n", updatedTemplate.State)

	// Wait for user to review the updated badge template
	fmt.Println("Badge template updated successfully. You can review it in your Credly account.")
	fmt.Println("Press Enter to continue with the demo...")
	_, _ = fmt.Scanln()

	// Archive the badge template
	fmt.Println("Archiving the badge template...")
	archivedTemplate, err := client.ArchiveBadgeTemplate(template.Id)
	if err != nil {
		fmt.Printf("Error archiving badge template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Badge template archived. New state: %s\n", archivedTemplate.State)

	// Wait for user to review the archived badge template
	fmt.Println("Badge template archived successfully. You can review it in your Credly account.")
	fmt.Println("Press Enter to continue with the demo...")
	_, _ = fmt.Scanln()

	// Unarchive the badge template
	fmt.Println("Unarchiving the badge template...")
	unarchivedTemplate, err := client.UnarchiveBadgeTemplate(template.Id)
	if err != nil {
		fmt.Printf("Error unarchiving badge template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Badge template unarchived. New state: %s\n", unarchivedTemplate.State)

	// Wait for user to review the unarchived badge template
	fmt.Println("Badge template unarchived successfully. You can review it in your Credly account.")
	fmt.Println("Press Enter to continue with the demo...")
	_, _ = fmt.Scanln()

	// Prompt user to delete the badge template
	fmt.Println("Do you want to delete the badge template? (yes/no)")
	var response string
	_, _ = fmt.Scanln(&response)
	if response != "yes" {
		fmt.Println("Skipping deletion of badge template. Demo completed successfully!")
		os.Exit(0)
	}

	fmt.Println("Deleting the badge template...")
	err = client.DeleteBadgeTemplate(template.Id)
	if err != nil {
		fmt.Printf("Error deleting badge template: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Badge template successfully deleted.")

	fmt.Println("Demo completed successfully!")
}
