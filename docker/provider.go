package docker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"docker_image_wait": dataSourceDockerImageWait(),
		},
		Schema: map[string]*schema.Schema{},
	}
}
