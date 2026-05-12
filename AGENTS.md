# AGENTS.md

本文件记录当前项目的重要上下文，供后续 Codex/工程师接手时快速理解项目。请优先阅读并遵守。

## 重要操作约束

- 禁止批量删除文件或目录。
- 不要使用：
  - `del /s`
  - `rd /s`
  - `rmdir /s`
  - `Remove-Item -Recurse`
  - `rm -rf`
- 需要删除文件时，只能一次删除一个明确路径的文件，例如：
  - `Remove-Item "C:\path\to\file.txt"`
- 如果需要批量删除文件，应停止操作并让用户手动删除。

## 项目概览

- 项目名：`MteSense Home`
- 类型：个人导航主页。
- 目标部署环境：Linux VPS。
- 架构：单 Go 服务托管后端 API、上传文件和 Vue 前端静态文件。
- 后端：Go、`chi`、SQLite、JWT、bcrypt。
- 前端：Vue 3、TypeScript、Vite、Pinia、Vue Router、vue-i18n、lucide-vue-next。
- 数据库：SQLite，默认路径 `data/app.db`。
- 上传目录：默认 `public_uploads/`。
- 前端构建产物：`web/app/dist/`。
- Linux 发布产物目录：`release/mtesense-home/`。

## 当前目录结构重点

```text
cmd/server/main.go                 # Go 服务入口
internal/config/                   # 环境变量配置
internal/db/                       # SQLite 打开与迁移执行
internal/migrations/               # 嵌入式 SQL 迁移
internal/auth/                     # 管理员初始化、登录、JWT
internal/http/                     # chi 路由、响应封装、中间件、SPA 托管
internal/nav/                      # 导航分组/链接业务
internal/settings/                 # 外观和搜索配置
internal/storage/                  # 上传文件保存
web/app/src/                       # Vue 前端源码
scripts/build-linux.sh             # Linux VPS 构建脚本
deploy/systemd/mtesense-home.service # systemd 服务示例
```

本地存在 `node_modules/`、`web/app/dist/`、`.cache/`、`data/app.db` 等生成/运行时内容；这些不应作为主要源码依据。

## 后端行为

启动流程在 `cmd/server/main.go`：

1. 读取环境变量配置。
2. 打开 SQLite 数据库并启用外键。
3. 执行 `internal/migrations/*.sql`。
4. 根据 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 初始化管理员账号。
5. 确保上传目录存在。
6. 启动 HTTP 服务。

默认环境变量见 `.env.example`：

```text
PORT=8080
DATABASE_PATH=data/app.db
UPLOAD_DIR=public_uploads
PUBLIC_SITE_URL=https://example.com
JWT_SECRET=replace-with-a-long-random-secret
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123456
```

注意：管理员账号只在用户名不存在时创建；修改环境变量里的 `ADMIN_PASSWORD` 不会覆盖数据库中已有用户密码。

## API 约定

所有 JSON 响应统一包在：

```json
{ "data": ... }
```

错误响应：

```json
{ "error": "message" }
```

公开接口：

```text
GET  /api/v1/health
POST /api/v1/auth/login
GET  /api/v1/navigation
GET  /api/v1/settings
```

登录后接口：

```text
GET    /api/v1/me
GET    /api/v1/admin/navigation
POST   /api/v1/admin/groups
PUT    /api/v1/admin/groups/{id}
DELETE /api/v1/admin/groups/{id}
POST   /api/v1/admin/links
PUT    /api/v1/admin/links/{id}
DELETE /api/v1/admin/links/{id}
PUT    /api/v1/admin/settings
POST   /api/v1/admin/uploads
```

管理员接口使用 `Authorization: Bearer <token>`。

## 数据模型

迁移文件：`internal/migrations/001_init.sql`

主要表：

- `users`：管理员用户，密码为 bcrypt 哈希。
- `nav_groups`：导航分组，含 `sort_order` 和 `visible`。
- `nav_links`：导航链接，含标题、URL、图标、描述、排序、显示状态、新窗口打开。
- `settings`：以 `key + value_json` 保存外观和搜索配置。

默认公开数据：

- 默认站点标题：`MteSense`
- 默认主题：`dark`
- 默认搜索引擎：`google`
- 启用搜索引擎：`google`、`bing`、`baidu`
- 默认分组：`APP`
- 默认链接：Google、Baidu

注意：迁移中默认设置使用 `ON CONFLICT DO NOTHING`，所以修改迁移默认值不会改变已有 `data/app.db` 里的设置。

## 前端行为

主页：`web/app/src/pages/HomePage.vue`

- 首页显示站点标题和当前时间。
- 副标题字段仍存在于设置结构和后台表单中，但当前首页不渲染副标题。
- 首页右侧有管理入口、语言切换、主题切换。
- 搜索框支持：
  - 输入时按链接标题、URL、描述本地过滤。
  - 回车或点击搜索按钮，按当前搜索引擎跳转。
- 搜索引擎定义在 `web/app/src/api/searchEngines.ts`：
  - Google: `https://www.google.com/search?q=...`
  - Bing: `https://www.bing.com/search?q=...`
  - Baidu: `https://www.baidu.com/s?wd=...`
- 当前 UI 调整：
  - 标题为 `MteSense`。
  - 首页标题字号已调小。
  - 标题/时间组合居中。
  - 搜索输入文字居中。
  - 首页不显示 `Personal navigation` 副标题。

状态存储：

- JWT token：`localStorage["mtesense_token"]`
- 主题：`localStorage["mtesense_theme"]`
- 语言：`localStorage["mtesense_locale"]`
- 搜索引擎：`localStorage["mtesense_search_engine"]`

路由：

```text
/                    # 首页
/admin/login         # 管理员登录
/admin/links         # 分组和链接管理
/admin/appearance    # 外观与搜索配置
```

## 前端样式

主样式文件：`web/app/src/styles/base.css`

- 使用 CSS 变量维护深色/浅色主题。
- 首页背景为渐变和模糊叠层，风格接近 Sun-Panel。
- 卡片、搜索栏、后台面板使用半透明背景和模糊效果。
- 主要响应式断点：
  - `980px`：导航网格变为 3 列，后台布局转单列。
  - `640px`：首页头部纵向排列，导航网格变为 1 列。

## 本地开发

后端：

```powershell
$env:GOCACHE='E:\Codex\mtesense_home\.cache\go-build'
$env:GOMODCACHE='E:\Codex\mtesense_home\.cache\gomod'
go run ./cmd/server
```

前端开发服务器：

```powershell
cd web/app
npm install
npm run dev
```

Vite 配置：`web/app/vite.config.ts`

- 开发端口：`5173`
- `/api` 和 `/uploads` 代理到 `http://localhost:8080`

生产预览时，先 `npm run build`，再启动 Go 服务访问 `http://127.0.0.1:8080`。

## Linux VPS 构建和部署

推荐使用：

```bash
chmod +x scripts/build-linux.sh
./scripts/build-linux.sh
```

脚本会：

1. 在 `web/app` 下运行 `npm ci`。
2. 运行 `npm run build`。
3. 使用 `GOOS=linux`、`CGO_ENABLED=0` 构建 Linux 二进制 `mtesense-home`。
4. 输出到 `release/mtesense-home/`。

systemd 示例：

```text
deploy/systemd/mtesense-home.service
```

默认部署目录为 `/opt/mtesense-home`，服务用户为 `mtesense`，可写目录为：

- `/opt/mtesense-home/data`
- `/opt/mtesense-home/public_uploads`

Linux 部署不需要 `mtesense-home.exe`。

## 验证命令

后端：

```powershell
$env:GOCACHE='E:\Codex\mtesense_home\.cache\go-build'
$env:GOMODCACHE='E:\Codex\mtesense_home\.cache\gomod'
go test ./...
```

前端：

```powershell
cd web/app
$env:npm_config_cache='E:\Codex\mtesense_home\.cache\npm'
npm run build
```

当前项目还没有专门的 Go 单元测试文件；`go test ./...` 主要验证编译和包加载。

## 已知注意事项

- `README.md` 和部分中文字符串在当前 PowerShell 输出中可能出现乱码显示；编辑中文文件时请确保使用 UTF-8。
- `internal/settings/settings.go` 中硬编码 fallback 仍有 `MteSense Home`，而迁移和前端默认值是 `MteSense`；如果要彻底统一默认标题，应同步修改该 fallback。
- `internal/migrations/001_init.sql` 中的中文默认链接描述在当前输出中也可能显示乱码；若要调整默认种子数据，请用 UTF-8 检查。
- 删除分组会因为 SQLite 外键 `ON DELETE CASCADE` 删除该分组下链接；这属于用户触发的单个资源删除，不要用批量文件删除命令处理任何文件。
- 上传限制：单文件最大 5MB；允许扩展名 `.png`、`.jpg`、`.jpeg`、`.webp`、`.gif`、`.svg`。
- 公开导航接口只返回 `visible = 1` 的分组和链接；后台接口返回全部。
