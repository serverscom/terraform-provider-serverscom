package serverscom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"strings"
)

func expandIntList(elements []any) []int {
	expandedIntList := make([]int, len(elements))
	for i, v := range elements {
		expandedIntList[i] = v.(int)
	}

	return expandedIntList
}

func expandedStringList(elements []any) []string {
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

// optionalConfigInt returns a pointer to the config value for key, or nil when
// the attribute is absent from configuration. Unlike schema.ResourceData.GetOk
// it treats an explicit 0 as a set value, which matters when the API
// distinguishes "unset" from "zero" via a pointer field (e.g. a DNS MX record
// with priority 0).
func optionalConfigInt(d *schema.ResourceData, key string) *int {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() {
		return nil
	}
	v := rawConfig.GetAttr(key)
	if v.IsNull() {
		return nil
	}
	i := d.Get(key).(int)
	return &i
}
