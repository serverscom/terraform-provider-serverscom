package serverscom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"strings"
)

func expandIntList(elements []interface{}) []int {
	expandedIntList := make([]int, len(elements))
	for i, v := range elements {
		expandedIntList[i] = v.(int)
	}

	return expandedIntList
}

func expandedStringList(elements []interface{}) []string {
	expandedStringList := make([]string, len(elements))
	for i, v := range elements {
		expandedStringList[i] = v.(string)
	}

	return expandedStringList
}

func normalizeString(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

func compareStrings(k, old, new string, d *schema.ResourceData) bool {
	return normalizeString(old) == normalizeString(new)
}
