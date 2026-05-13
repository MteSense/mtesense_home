# MteSense Home

个人导航主页，前端使用 Vue 3 + TypeScript，后端使用 Go + SQLite。Go 服务同时提供 API、上传文件和前端静态页面，目标部署环境为 Linux VPS。

## 功能

- 前台导航主页：分组导航、搜索框、搜索引擎跳转、语言切换、白天/黑夜主题。
- 搜索：输入时按标题、URL、描述本地过滤；按 Enter 或点击搜索按钮跳转到 Google、Bing 或 Baidu。
- 管理后台：管理员登录后维护分组、链接、外观、默认搜索引擎和启用的搜索引擎。
- 上传：支持背景图和图标类图片上传到 `public_uploads/`。

## 环境变量

复制 `.env.example` 为 `.env` 后按需设置：

```text
PORT=8080
DATABASE_PATH=data/app.db
UPLOAD_DIR=public_uploads
PUBLIC_SITE_URL=https://example.com
JWT_SECRET=replace-with-a-long-random-secret
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123456
```

首次启动时会自动创建管理员账号。如果用户已存在，后续修改 `ADMIN_PASSWORD` 不会覆盖数据库里的密码。

## 本地开发

后端：

```bash
go run ./cmd/server
```

前端：

```bash
cd web/app
npm install
npm run dev
```

开发时访问 `http://localhost:5173`，Vite 会把 `/api` 和 `/uploads` 代理到 `http://localhost:8080`。

## Linux VPS 构建

推荐在 Linux VPS 上直接构建：

```bash
sudo apt update
sudo apt install -y curl git

# 安装 Node.js、Go 后进入项目目录
cd /opt/mtesense-home
chmod +x scripts/build-linux.sh
./scripts/build-linux.sh
```

构建产物会生成到：

```text
release/mtesense-home/
  mtesense-home
  web/app/dist/
  data/
  public_uploads/
  .env.example
```

也可以在项目根目录手动构建：

```bash
cd web/app
npm ci
npm run build
cd ../..
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o mtesense-home ./cmd/server
./mtesense-home
```

注意：Linux 部署不需要 `mtesense-home.exe`。

## Docker 部署

Docker 部署支持三种模式：可以在本地或 VPS 上用源码直接构建运行，也可以直接拉取 GHCR 镜像运行，或者先把镜像发布到 Docker Hub 后再部署。

容器内默认路径：

```text
PORT=8080
DATABASE_PATH=/app/data/app.db
UPLOAD_DIR=/app/public_uploads
PUBLIC_SITE_URL=https://example.com
```

SQLite 数据和上传文件建议使用 Docker 卷持久化。首次启动时会根据 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 创建管理员账号；如果用户已经存在，后续修改环境变量不会覆盖数据库里的密码。生产环境请务必修改 `JWT_SECRET`、`ADMIN_USERNAME` 和 `ADMIN_PASSWORD`。

### 模式一：源码构建部署

适合 VPS 上有项目源码的情况：

```bash
docker compose up -d --build
```

也可以不用 Compose，手动构建并运行：

```bash
docker build -t mtesense-home:local .

docker run -d --name mtesense-home \
  -p 8080:8080 \
  -e JWT_SECRET=replace-with-a-long-random-secret \
  -e PUBLIC_SITE_URL=https://example.com \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=admin123456 \
  -v mtesense_data:/app/data \
  -v mtesense_uploads:/app/public_uploads \
  --restart unless-stopped \
  mtesense-home:local
```

### 模式二：GHCR 镜像部署

适合直接使用 GitHub Container Registry 上已经构建好的镜像部署：

```bash
docker pull ghcr.io/mtesense/mtesense-home:latest
```

在 VPS 上直接运行：

```bash
docker run -d --name mtesense-home \
  -p 8080:8080 \
  -e JWT_SECRET=replace-with-a-long-random-secret \
  -e PUBLIC_SITE_URL=https://example.com \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=admin123456 \
  -v mtesense_data:/app/data \
  -v mtesense_uploads:/app/public_uploads \
  --restart unless-stopped \
  ghcr.io/mtesense/mtesense-home:latest
```

如果 GHCR 镜像保持私有，需要先在 VPS 上登录：

```bash
docker login ghcr.io
```

### 模式三：Docker Hub 镜像部署

适合先构建并发布镜像，然后在任意 Linux VPS 上拉取即部署。先把 `your-dockerhub-name` 替换成自己的 Docker Hub 用户名或组织名：

```bash
docker build -t your-dockerhub-name/mtesense-home:latest .
docker login
docker push your-dockerhub-name/mtesense-home:latest
```

在 VPS 上直接运行：

```bash
docker run -d --name mtesense-home \
  -p 8080:8080 \
  -e JWT_SECRET=replace-with-a-long-random-secret \
  -e PUBLIC_SITE_URL=https://example.com \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=admin123456 \
  -v mtesense_data:/app/data \
  -v mtesense_uploads:/app/public_uploads \
  --restart unless-stopped \
  your-dockerhub-name/mtesense-home:latest
```

或者编辑 `docker-compose.hub.yml`，把 `image` 改成自己的 Docker Hub 镜像名，然后运行：

```bash
docker compose -f docker-compose.hub.yml up -d
```

常用维护命令：

```bash
docker compose logs -f
docker compose down
docker compose pull
docker compose up -d
```

## systemd 部署

示例以 `/opt/mtesense-home` 为部署目录：

```bash
sudo useradd --system --home /opt/mtesense-home --shell /usr/sbin/nologin mtesense
sudo mkdir -p /opt/mtesense-home
sudo cp -R release/mtesense-home/* /opt/mtesense-home/
sudo cp /opt/mtesense-home/.env.example /opt/mtesense-home/.env
sudo nano /opt/mtesense-home/.env
sudo chown -R mtesense:mtesense /opt/mtesense-home

sudo cp deploy/systemd/mtesense-home.service /etc/systemd/system/mtesense-home.service
sudo systemctl daemon-reload
sudo systemctl enable --now mtesense-home
sudo systemctl status mtesense-home
```

服务启动后访问：

```text
http://your-vps-ip:8080
```

如果使用 Nginx 反向代理，把外部 80/443 转发到 `127.0.0.1:8080` 即可。

## 数据备份

SQLite 数据默认位于 `data/app.db`。停止服务后备份这个文件即可。上传资源默认位于 `public_uploads/`，需要和数据库一起备份。
