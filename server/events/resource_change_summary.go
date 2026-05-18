package events

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TerraformPlanJSON represents the structure of terraform show -json output
type TerraformPlanJSON struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

// ResourceChange represents a single resource change in the plan
type ResourceChange struct {
	Address string       `json:"address"`
	Change  ChangeDetail `json:"change"`
}

// ChangeDetail contains the actions for a resource change
type ChangeDetail struct {
	Actions []string `json:"actions"`
}

// GenerateResourceChangeSummary converts terraform plan/apply JSON into a human-readable summary
// Works for both plan output and apply output since they share the same resource_changes structure
// Filters out resources with no-op actions to show only meaningful changes
func GenerateResourceChangeSummary(planJSON string) string {
	if planJSON == "" {
		return ""
	}

	var plan TerraformPlanJSON
	if err := json.Unmarshal([]byte(planJSON), &plan); err != nil {
		return fmt.Sprintf("Error parsing plan JSON: %v", err)
	}

	resourceChanges := plan.ResourceChanges
	if len(resourceChanges) == 0 {
		return "No resource changes in this plan."
	}

	// Filter out no-op resources
	var meaningfulChanges []ResourceChange
	for _, resource := range resourceChanges {
		if !isNoOp(resource.Change.Actions) {
			meaningfulChanges = append(meaningfulChanges, resource)
		}
	}

	// If all changes were no-ops, return early
	if len(meaningfulChanges) == 0 {
		return "No resource changes in this plan."
	}

	lines := []string{
		fmt.Sprintf("Plan Summary (%d resource(s)):", len(meaningfulChanges)),
		strings.Repeat("-", 60),
	}

	for _, resource := range meaningfulChanges {
		actionDesc := getActionDescription(resource.Change.Actions)
		lines = append(lines, fmt.Sprintf("%s - %s", resource.Address, actionDesc))
	}

	return strings.Join(lines, "\n")
}

// isNoOp checks if actions represent a no-op (no actual changes)
func isNoOp(actions []string) bool {
	if len(actions) == 0 {
		return true
	}
	if len(actions) == 1 && actions[0] == "no-op" {
		return true
	}
	return false
}

// getActionDescription converts action list to human-readable description
func getActionDescription(actions []string) string {
	if len(actions) == 0 {
		return "no-op"
	}

	if len(actions) == 1 {
		actionMap := map[string]string{
			"create": "will be created",
			"read":   "will be read",
			"update": "will be updated",
			"delete": "will be destroyed",
			"no-op":  "no changes",
		}
		if desc, ok := actionMap[actions[0]]; ok {
			return desc
		}
		return actions[0]
	}

	// Multiple actions usually indicate replacement
	hasDelete := false
	hasCreate := false
	for _, action := range actions {
		if action == "delete" {
			hasDelete = true
		}
		if action == "create" {
			hasCreate = true
		}
	}

	if hasDelete && hasCreate {
		if actions[0] == "create" {
			return "will be replaced (create before destroy)"
		}
		return "will be replaced (destroy and create)"
	}

	return fmt.Sprintf("actions: %s", strings.Join(actions, ", "))
}
