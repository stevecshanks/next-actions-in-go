# Next Actions in Go

## Running in dev mode

```
make dev
```

## Running tests

```
make test
```

## Running in production

The master branch will continuously deploy from CircleCI to the specified server. Requires the server to have `docker-compose` installed, and a `docker-deploy` user with an SSH key:

```
adduser --disabled-password docker-deploy
usermod -aG docker docker-deploy

su docker-deploy
ssh-keygen -t rsa
cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys

cat ~/.ssh/id_rsa
```

Copy this private key and add it to CircleCI settings for the project.

You'll also need to set some environment variables in CircleCI:

- `DOCKER_HUB_TOKEN` - an access token for Docker Hub that can write to the relevant repositories
- `DOCKER_SERVER` - the IP address or hostname for the server you want to deploy to
- `SSH_KNOWN_HOSTS` - will be appended to the `~/.ssh/known_hosts` file in CI to allow SSH access to your Docker server without prompting to confirm the identity of the host. Can be generated using `ssh-keyscan $DOCKER_SERVER`
