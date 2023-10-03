package main

import (
	"context"
	"fmt"
)

type NextjsBuild struct {}

const (
	// https://hub.docker.com/_/node
	nodeJSVersion = "18"

	// https://hub.docker.com/_/alpine
	alpineVersion = "3.18"
)

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

	return build.Directory("out/"), nil
}
