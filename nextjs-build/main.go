package main

import (
	"fmt"
)

type NextjsBuild struct{}

const (
	// https://hub.docker.com/_/node
	nodeJSVersion = "slim"

	// https://hub.docker.com/_/nginx
	nginxVersion = "1.25"

	// https://hub.docker.com/_/alpine
	alpineVersion = "3.18"
)

func (d *Directory) NextJSBuild() *Container {
	return (&NextjsBuild{}).Build(d)
}

// Use a directory to build a NextJS app.
func (m *NextjsBuild) Build(source *Directory) *Container {
	node := dag.Container(ContainerOpts{Platform: Platform("linux/amd64")}).
		From(fmt.Sprintf("node:%s", nodeJSVersion))
	yarnCachePath := "/Users/slumbering/Library/Caches/Yarn/v6"
	yarnCacheVolume := fmt.Sprintf("dagger-io-yarn-%s", nodeJSVersion)

	return node.Pipeline("build").
		WithDirectory("/app", source, ContainerWithDirectoryOpts{
			Exclude: []string{
				"node_modules",
				".next",
				"out",
			},
		}).
		WithMountedCache(yarnCachePath, dag.CacheVolume(yarnCacheVolume)).
		WithEnvVariable("YARN_CACHE_FOLDER", yarnCachePath).
		WithWorkdir("/app").
		WithExec([]string{"yarn", "install"}).
		WithExec([]string{"yarn", "export"})
}
