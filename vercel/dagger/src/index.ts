/**
 * A simple wrapper around the vercel cli
 * Vercel module provides a simple way to deploy, list and remove deployments from vercel
 */

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
    // create a cache volume
    const nodeCache = dag.cacheVolume("node")

    return dag
      .container()
      .from("node:lts-slim")
      .withSecretVariable("VERCEL_TOKEN", this.token)
      .withMountedDirectory('/app', this.currentWorkdir)
      .withMountedCache("/src/node_modules", nodeCache)
      .withWorkdir("/app")
      .withExec(["npm", "i", "-g", "vercel"])
  }
}


@object()
class Vercel {

  /**
   * Deploy the current directory to vercel
   * @param currentWorkdir path to the directory to deploy
   * @param token vercel token
   * @returns deployment URL
   */
  @func()
  async vercelDeploy(currentWorkdir: Directory, token: Secret): Promise<string> {
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

  /**
   * List all deployments in vercel
   * @param currentWorkdir path of the directory to list deployments
   * @param token vercel token
   * @returns list of deployments
   */
  @func()
  async vercelList(currentWorkdir: Directory, token: Secret): Promise<string> {
    const vercel = new VercelOptions(currentWorkdir, token)

    return await vercel
      .base()
      .withEnvVariable("CACHEBUSTER", Date.now().toString()) // invalidate cache to get a fresh list
      .withExec([
        "sh",
        "-c",
        "vercel --token $VERCEL_TOKEN list --yes"
      ])
      .stdout()
  }

  /**
   * Remove a deployment from vercel
   * @param currentWorkdir path of the directory to remove deployment
   * @param deploymentURL URL of the deployment to remove
   * @param token vercel token
   * @returns log of the removal
   */
  @func()
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
