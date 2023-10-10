package main

import (
	"context"
	"fmt"
	"time"
)

type Flyio struct {}

const (
	// https://hub.docker.com/r/flyio/flyctl/tags
	flyctlVersion = "0.1.103"
	appMemoryMB = "256"
	// https://fly.io/docs/reference/configuration/#picking-a-deployment-strategy
	deployStrategy = "bluegreen"
	// wait this many seconds for the app to finish deploying
	waitSeconds = "60"
	// https://fly.io/dashboard/dagger
	flyOrg = "dagger"

	appImageRegistry = "registry.fly.io"
)

func flyTokenSecret(token string) *Secret {
	if token == "" {
		panic("FLY_API_TOKEN env var must be set")
	}
	return dag.SetSecret("FLY_API_TOKEN", token)
}

func flyctl(appName string, token string) *Container {
	c := dag.Pipeline("flyctl")
	flyctl := c.Container(ContainerOpts{Platform: Platform("linux/amd64")}).Pipeline("auth").
		From(fmt.Sprintf("flyio/flyctl:v%s", flyctlVersion)).
		WithSecretVariable("FLY_API_TOKEN", flyTokenSecret(token)).
		WithEnvVariable("RUN_AT", time.Now().String()).
		WithNewFile("fly.toml", ContainerWithNewFileOpts{
			Contents: fmt.Sprintf(`
# https://fly.io/docs/reference/configuration/
app = "%s"

[http_service]
	internal_port = 80
	force_https = true

[[http_service.checks]]
	interval = "5s"
	timeout = "4s"
	method = "GET"
	path = "/"
			`, appName)})

	return flyctl
}

func (d *Directory) FlyioDeploy(ctx context.Context, appName, token string) (*Container, error) {

	app := d.NextJsbuild()

	_, err := flyctl(appName, token).
		WithExec([]string{"status"}).
		Sync(ctx)
	if err != nil {
		_, err = flyctl(appName, token).
			WithExec([]string{"apps", "create", appName, "--org", flyOrg}).
			Sync(ctx)
		if err != nil {
			panic(err)
		}
	}

	appImageRef, err := app.
		WithRegistryAuth(appImageRegistry, "x", flyTokenSecret(token)).
		Publish(ctx, fmt.Sprintf("%s/%s", appImageRegistry, appName))
	if err != nil {
		panic(err)
	}

	return flyctl(appName, token).Pipeline("deploy").
		WithDirectory("/app", d).
		WithExec([]string{
			"deploy", "--now",
			"--app", appName,
			"--image", appImageRef,
			"--vm-memory", appMemoryMB,
			"--ha=false", // these are preview environments, we do not want to pay for 2 instances (a.k.a. HA)
			"--strategy", deployStrategy,
			"--wait-timeout", waitSeconds,
		}).
		Sync(ctx)
}
