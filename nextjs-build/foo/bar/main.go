// nextjs-build is a Dagger module that provides a method for building a Next.js application within a containerized environment.
//
// The module is implemented as a struct NextjsBuild with a method Build that takes a source directory (source) as input, representing the root of the Next.js project.
// The resulting container is configured for serving the static site, with Nginx acting as the web server.

package main

import (
	"context"
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

// Build initiates the construction process of a Next.js application within a containerized environment.
// This method takes a source directory (source) as input, representing the root of the Next.js project.
// The resulting container is configured for serving the static site, with Nginx acting as the web server.
func (m *NextjsBuild) Build(ctx context.Context, source *Directory) *Container {
	node := dag.Container().
		From(fmt.Sprintf("node:%s", nodeJSVersion))
	yarnCachePath := "/Users/slumbering/Library/Caches/Yarn/v6"
	yarnCacheVolume := fmt.Sprintf("dagger-io-yarn-%s", nodeJSVersion)

	build := node.Pipeline("build").
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

	return dag.Container(ContainerOpts{Platform: Platform("linux/amd64")}).
		From(fmt.Sprintf("nginx:%s", nginxVersion)).
		WithDirectory("/usr/share/nginx/html", build.Directory("/app/out")).
		WithNewFile("/etc/nginx/conf.d/default.conf", ContainerWithNewFileOpts{
			Contents: `
				server {
					listen       80;
					listen  [::]:80;
					server_name  localhost;

					location / {
						root   /usr/share/nginx/html;
						try_files $uri $uri/index.html $uri.html =404;
					}

					error_page  404 	/404.html;

					error_page  500 502 503 504 	/50x.html;
				}
			`})
}
