package serverscom

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func createClient() (*scgo.Client, error) {
	if os.Getenv("SERVERSCOM_TOKEN") == "" {
		return nil, fmt.Errorf("you must set SERVERSCOM_TOKEN")
	}

	if os.Getenv("SERVERSCOM_API_URL") == "" {
		return nil, fmt.Errorf("you must set SERVERSCOM_API_URL")
	}

	return scgo.NewClientWithEndpoint(os.Getenv("SERVERSCOM_TOKEN"), os.Getenv("SERVERSCOM_API_URL")), nil
}
