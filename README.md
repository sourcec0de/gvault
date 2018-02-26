GVAULT
===
> Manage project `secrets` using Google Cloud KMS

This project is intended to allow Google Cloud Platform users to easily take advantage of Google Cloud KMS.
GVault enables you to store and manage secrets right in source control. By storing your secrets in source control
you can more easily package your application and it's required configuration for production. This also removes
the need for sharing secrets in traditional ways. You can simply grant read / write access using Google existing IAM roles.

# Demo
[![asciicast](https://asciinema.org/a/nLFlLdfMDEnv0AYqGkWcepK0w.png)](https://asciinema.org/a/nLFlLdfMDEnv0AYqGkWcepK0w)

# Benefits

### #1 - No Manual Key Management
Since GVault is built on top of Google Cloud KMS it's users simply need to be logged into the `gcloud` CLI
or have the `GOOGLE_APPLICATION_CREDENTIALS` environmental variable pointing to a service account with access to KMS resources.
Beyond that access to keys and keyrings is controlled via IAM. Key management is handled by Google Cloud KMS.
Keys are automatically rotated on a regular basis. You can easily give any need to know personel access to your keys
and therefore your vaults.

### Simple
Gvault has a very small footprint and CLI surface that any developer will be able to easily command.
Once the project is initialize adding a secret is as simple as `gvault secrets add MYSQL_PASSWORD=s71Dbl01-Z`

### No servers
Gvault stores your encrypted secrets in your projects repository. You can use your SCM tool of choice to track changes to secrets
without the worry of leaking them.


### Integrates with Google Container Builder
GVault support generating configurations for your `cloudbuild.yml`

### Integrates with Kubernetes
GVault can sync with kubernetes by creating versioned secrets that match your vaults contents.
Once the secret is in kubernetes you are free to mount it however you like.


# Getting Started
First install gvault

### Initialize
This will prompt you to set defaults for your vault.
- project
- location
- keyring
- key
```sh
cd ~/project_dir
gvault init
```

### Add a secret
```sh
gvault secrets add MYSQL_PASSWORD=s71Dbl01-Z
```

### Remove a secret
```sh
gvault secrets remove MYSQL_PASSWORD
```

### Retrieve a secret
```sh
gvault secrets get MYSQL_PASSWORD
```

### Import all key value pairs from a .env file
```sh
gvault secrets import /path/to/.env
```
