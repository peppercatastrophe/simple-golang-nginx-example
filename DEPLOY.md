# Deployment

CI builds both Go service images on GitHub Actions and pushes them to GHCR.
Merging to `main` triggers the pipeline; the VPS only **pulls** images and
recreates containers — it never compiles anything, which keeps the 1 CPU / 1 GB
box happy.

## Architecture

```
GitHub Actions (build)                VPS (1 CPU / 1 GB)
  build main-service  -> ghcr.io/.../main-service:latest   +---> docker compose -f docker-compose.prod.yml
  build items-service -> ghcr.io/.../items-service:latest  |     nginx:8080 -> main-service:8081
  ssh deploy                                                  -> items-service:8082
```

## One-time setup on the VPS

1. Install Docker + compose plugin (Debian/Ubuntu):

   ```bash
   curl -fsSL https://get.docker.com | sh
   sudo systemctl enable --now docker
   ```

2. Make sure the user that GitHub Actions connects as can use the docker
   socket (if it isn't `root`). Without this, deploys fail with
   `permission denied while trying to connect to the docker API`:

   ```bash
   sudo usermod -aG docker "$USER"
   ```

3. Add a 1 GB swap file (recommended safety margin on a 1 GB box):

   ```bash
   sudo fallocate -l 1G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

4. Clone the repo:

   ```bash
   sudo mkdir -p /opt/simple-golang-nginx-example
   sudo chown "$USER" /opt/simple-golang-nginx-example
   git clone git@github.com:peppercatastrophe/simple-golang-nginx-example.git \
     /opt/simple-golang-nginx-example
   ```

5. Create the env file (gitignored, never committed):

   ```bash
   cd /opt/simple-golang-nginx-example
   cat > .env <<'EOF'
   JWT_SECRET=change-me-to-a-long-random-string
   JWT_EXPIRES_MINUTES=60
   MAIN_SERVICE_URL=http://main-service:8081
   EOF
   ```

6. Test once:

   ```bash
   docker compose -f docker-compose.prod.yml pull
   docker compose -f docker-compose.prod.yml up -d
   ```

## GitHub secrets

Add these to **Settings → Secrets and variables → Actions**:

| Secret          | Value                                                        |
| --------------- | ------------------------------------------------------------ |
| `VPS_HOST`      | VPS IP or domain                                             |
| `VPS_USER`      | SSH user (e.g. `root` or your user)                          |
| `VPS_SSH_KEY`   | Private key that can log into the VPS                        |
| `GHCR_TOKEN`    | PAT with **Packages: Read** scope (used by the VPS to pull)  |

Create the PAT at https://github.com/settings/tokens (fine-grained, repo = this repo, permissions = Packages: Read). Give it `read:packages` on the account-level classic PAT if you prefer.

## Flow

1. Open a PR against `main`. Merging it pushes to `main`.
2. The `Build & Deploy` workflow builds and pushes both images, then SSHs into the VPS.
3. On the VPS: `git pull` → `docker compose pull` → `up -d` → prune old images.

Roll back by SSHing into the VPS and:

```bash
cd /opt/simple-golang-nginx-example
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```
