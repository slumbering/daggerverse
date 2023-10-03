package main

import (
	"context"
	"fmt"
)

type NextjsBuild struct{}

const (
	// https://hub.docker.com/_/node
	nodeJSVersion = "18"

	// https://hub.docker.com/_/alpine
	alpineVersion = "3.18"
)

// Use a directory to build a NextJS app.
func (d *Directory) BuildNextJS(ctx context.Context) (*Directory, error) {
	node := dag.Container(ContainerOpts{Platform: Platform("linux/amd64")}).
		From(fmt.Sprintf("node:%s", nodeJSVersion))
	yarnCachePath := fmt.Sprintf("/usr/local/share/.cache/yarn/%s", nodeJSVersion)
	yarnCacheVolume := fmt.Sprintf("dagger-io-yarn-%s", nodeJSVersion)

	build, err := node.Pipeline("build").
		WithDirectory("/app", d).
		WithMountedCache(yarnCachePath, dag.CacheVolume(yarnCacheVolume)).
		WithEnvVariable("YARN_CACHE_FOLDER", yarnCachePath).
		WithWorkdir("/app").
		WithExec([]string{"yarn", "install"}).
		WithExec([]string{"yarn", "export"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	// Export the build directory to be used with other modules, such as deploying with `flyctl deploy`.
	return build.Directory("out/"), nil
}
