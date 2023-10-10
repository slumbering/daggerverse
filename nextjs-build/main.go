package main

import (
	"fmt"
)

type NextjsBuild struct{}

const (
	// https://hub.docker.com/_/node
	nodeJSVersion = "18"

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
	yarnCachePath := fmt.Sprintf("/usr/local/share/.cache/yarn/%s", nodeJSVersion)
	yarnCacheVolume := fmt.Sprintf("dagger-io-yarn-%s", nodeJSVersion)

	build := node.Pipeline("build").
		WithDirectory("/app", source).
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
