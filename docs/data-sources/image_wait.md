# Data Source: docker\_image\_wait

Waits until the specified docker image is available.

## Example Usage

```hcl
data "docker_image_wait" "example_image_wait" {
  registry = "yourregistry.azurecr.io"
  username = "username"
  password = "password"
  image    = "example-image:1.0.0"
  # 30 minutes
  timeout = 1800
}

# In your container definition:

resource "azurerm_container_group" "example_container_group" {
  # ...

  container {
    name   = "example"
	# full_image returns the full string you need.
	# It's also used to tell Terraform that this resource
	# depends on the docker_image_wait data source
    image  = data.docker_image_wait.example_image_wait.full_image
    cpu    = "0.3"
    memory = "0.3"

    # ...
  }
}
```

## Argument Reference

* `registry` - (Optional) Registry URL. Default: `registry.hub.docker.com`.
* `username` - (Optional) Username to log in at the registry.
* `password` - (Optional) Password to log in at the registry.
* `image` - (Required) Image to wait for (for example, `postgresql:13`).
* `fail_after_timeout` - (Optional) Whether to return an error if waiting times out after [`timeout`](#timeout) seconds.

## Attribute Reference

* `id` - (String) An ID that changes every time to make sure waiting isn't skipped.
* `exists` - (Bool) `true`, if the image exists after waiting for it.
* `full_image` - (String) The full image reference including registry, repository and tag.

## Timeouts

* `timeout` - (Optional) How long to wait for the image, in seconds, before timing out. Default: 600

## Import

This resource cannot be imported.


