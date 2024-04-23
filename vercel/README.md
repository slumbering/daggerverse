[![tested-with-dagger-version](https://img.shields.io/badge/Tested%20with%20dagger-0.11.0-success?style=for-the-badge)](https://github.com/dagger/dagger/releases/tag/v0.11.0)

# Deploy to Vercel

This module aims to deploy your projects to Vercel

## Usage

### Deploy to Vercel

```shell
dagger call vercel-deploy --current-workdir my/project/workdir --token env:VERCEL_TOKEN
```

### List available sites

```shell
dagger call vercel-list --current-workdir my/project/workdir --token env:VERCEL_TOKEN
```

### Remove a deployment

```shell
dagger call vercel-remove --current-workdir my/project/workdir --token env:VERCEL_TOKEN --deployment-url https://app-my-project-id.vercel.app

```

### Todo
| Command                | Done |
|------------------------|------|
| Deploy a project to Vercel  | ✅    |
| List recent deployments for the current Vercel Project | ✅    |
| Build a Vercel Project locally or in a CI environment       | ❌    |
| Remove a deployment       | ✅    |
