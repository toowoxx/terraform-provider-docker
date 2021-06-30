package docker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	"gitlab.com/xdevs23/go-collections"
)

func dataSourceDockerImageWait() *schema.Resource {
	return &schema.Resource{
		Description: "Waits for a docker image to be available",
		ReadContext: dataSourceDockerImageWaitRead,
		Schema: map[string]*schema.Schema{
			"registry": {
				Description: "Registry URL",
				Type:        schema.TypeString,
				Default:     "registry.hub.docker.com",
				Optional:    true,
				ForceNew:    true,
			},
			"username": {
				Description: "Username to log in",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"password": {
				Description: "Password to log in",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"image": {
				Description: "Docker image",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"timeout": {
				Description: "How long to wait, in seconds",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
				ForceNew:    true,
			},
			"id": {
				Description: "Returns an ID that changes every time",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"exists": {
				Description: "Returns true if the image exists after waiting for it",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceDockerImageWaitRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	url := "https://" + d.Get("registry").(string)
	username := ""
	password := ""
	if strI, ok := d.GetOk("username"); ok {
		username = strI.(string)
	}
	if strI, ok := d.GetOk("password"); ok {
		password = strI.(string)
	}
	image := d.Get("image").(string)

	var repository string
	tag := "latest"
	splitImage := strings.Split(image, ":")
	if len(splitImage) == 0 {
		return diag.Errorf("invalid image \"%s\", please specify an image", image)
	}
	if len(splitImage) >= 1 {
		repository = splitImage[0]
	}
	if len(splitImage) == 2 {
		tag = splitImage[1]
	}
	if len(splitImage) > 2 {
		return diag.Errorf("invalid image \"%s\": found more than one colon", image)
	}

	if err := waitForImage(url, username, password, repository, tag, d.Get("timeout").(int)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", time.Now().Unix()))
	if err := d.Set("exists", true); err != nil {
		return diag.FromErr(errors.Wrap(err, "bug: could not set 'exists' output"))
	}

	return nil
}

func waitForImage(url string, username string, password string, repository string, tag string, timeout int) error {
	hub, err := registry.New(url, username, password)
	if err != nil {
		return errors.Wrap(err, "could not connect/log in to registry")
	}

	ranIntoTimeout := false

	time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		ranIntoTimeout = true
	})

	for {
		tags, err := hub.Tags(repository)
		if err != nil || !collections.Include(tags, tag) {
			if err != nil {
				log.Println(err)
			}
			if ranIntoTimeout {
				break
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	return nil
}
