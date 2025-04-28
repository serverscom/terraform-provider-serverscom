package serverscom

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HashString(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func HashStringStateFunc() schema.SchemaStateFunc {
	return func(v interface{}) string {
		switch v.(type) {
		case string:
			return HashString(v.(string))
		default:
			return ""
		}
	}
}

// hashcode.String in the terraform-plugin-sdk was made internal to the SDK in v2.
// Embed the implementation here to allow same hash function to continue to be used
// by the code in this provider that used it for hash computation.
func SDKHashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// hashFilter returns md5 hash of filter
func hashFilter(filter map[string]any) (string, error) {
	filterBytes, err := json.Marshal(filter)
	if err != nil {
		return "", fmt.Errorf("error marshaling filter: %s", err)
	}
	hash := md5.Sum(filterBytes)
	return hex.EncodeToString(hash[:]), nil
}
