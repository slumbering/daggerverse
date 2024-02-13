import { dag, Container, object, func, Secret, Directory, field } from "@dagger.io/dagger"


@object()
class VercelOptions {
  @field()
  currentWorkdir: Directory

  @field()
  token: Secret

  @field()
  deploymentURL?: string

  constructor(currentWorkdir: Directory, token: Secret, deploymentURL?: string) {
    this.currentWorkdir = currentWorkdir
    this.token = token
    this.deploymentURL = deploymentURL
  }

  @func()
  // Set up a container with the vercel cli installed
  base(): Container {
    return dag
      .container()
      .from("node:lts-slim")
      .withExec(["yarn", "global", "add", "vercel"])
      .withSecretVariable("VERCEL_TOKEN", this.token)
      .withMountedDirectory('/app', this.currentWorkdir)
      .withWorkdir("/app")
  }
}


@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Vercel {

  @func()
  // Deploy the current directory to vercel
  async vercelProd(currentWorkdir: Directory, token: Secret): Promise<string> {
    const vercel = new VercelOptions(currentWorkdir, token)
    return await vercel
      .base()
      .withExec([
        "sh",
        "-c",
        "vercel --token $VERCEL_TOKEN --yes"
      ])
      .stdout()
  }

  @func()
  // List available sites for the current directory
  async vercelList(currentWorkdir: Directory, token: Secret): Promise<string> {
    const vercel = new VercelOptions(currentWorkdir, token)

    return await vercel
      .base()
      .withExec([
        "sh",
        "-c",
        "vercel --token $VERCEL_TOKEN list --yes"
      ])
      .stdout()
  }

  @func()
  // Remove a given deployment from vercel
  async vercelRemove(currentWorkdir: Directory, deploymentURL: string, token: Secret): Promise<string> {
    const vercel = new VercelOptions(currentWorkdir, token)

    return await vercel
      .base()
      .withExec([
        "sh",
        "-c",
        `vercel --token $VERCEL_TOKEN remove ${deploymentURL} --yes`
      ])
      .stdout()
  }
}
