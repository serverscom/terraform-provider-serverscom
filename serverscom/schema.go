package serverscom

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setResourceDataFromMap(d *schema.ResourceData, m map[string]any) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("Unable to set `%s` attribute: %s", key, err)
		}
	}

	return nil
}
