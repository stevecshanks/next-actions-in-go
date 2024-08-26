# Next Actions in Go

## Running in dev mode

If you haven't already, copy `.env.example` to `.env` and set your environment variables.

```
make dev
```

## Running tests

```
make test
```

## Running in production

The `main` branch will continuously deploy from CircleCI to the specified server. Requires the server to have `docker-compose` installed, and a `docker-deploy` user with an SSH key:

```
adduser --disabled-password docker-deploy
usermod -aG docker docker-deploy

su docker-deploy
ssh-keygen -t ed25519
cp ~/.ssh/id_ed25519.pub ~/.ssh/authorized_keys

cat ~/.ssh/id_ed25519
```

Copy this private key and add it to CircleCI settings for the project as an Additional SSH Key, with the hostname set to the server you want to deploy to.

You'll also need to set some environment variables in CircleCI:

- `DOCKER_HUB_TOKEN` - an access token for Docker Hub that can write to the relevant repositories
- `DOCKER_SERVER` - the IP address or hostname for the server you want to deploy to
- `SSH_KNOWN_HOSTS` - will be appended to the `~/.ssh/known_hosts` file in CI to allow SSH access to your Docker server without prompting to confirm the identity of the host. Can be generated using `ssh-keyscan $DOCKER_SERVER`

Finally, on the server copy `.env.example` to `.env` and set your environment variables.
