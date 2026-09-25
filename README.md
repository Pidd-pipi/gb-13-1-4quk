# 校园二手书交易平台（campusbooks）

面向高校学生的二手书交易平台：学生用学号 + 学校邮箱注册（邮箱验证码），发布闲置教材与书籍出售/交换，按学科、课程名搜索所需教材，通过站内消息沟通交易，交易完成后互评，平台自动记录收藏与浏览历史并推荐同院系同学发布的书籍。

## 快速启动（Docker Compose，首选）

```bash
docker compose up -d --build
```

启动后访问：

| 入口 | 地址 |
| --- | --- |
| 前端 | http://localhost:8011 |
| 后端 API | http://localhost:3011/api/v1 |
| 健康检查 | http://localhost:3011/healthz |
| MinIO 控制台 | http://localhost:9008 （minioadmin / minioadmin） |

演示账号（由种子数据自动创建）：

| 角色 | 账号 | 密码 |
| --- | --- | --- |
| 管理员 | admin@campusbooks.local | admin123 |
| 学生 | zhang@campusbooks.local / li@campusbooks.local / wang@campusbooks.local / chen@campusbooks.local | student123 |

> 注意：docker-compose.yml 顶层声明 `name: campusbooks`，所有容器名带 `${COMPOSE_PROJECT_NAME:-campusbooks}` 前缀，因此项目可放在任意目录名（包括中文目录）下正常启动。

## 主要功能

1. **用户注册与认证**：学号 + 学校邮箱注册，邮箱验证码校验（开发环境 `send-code` 直接返回 `dev_code` 便于联调），JWT 登录，完善姓名/院系/联系方式，支持头像上传（MinIO）。
2. **发布闲置书籍**：书名/作者/ISBN/原价/售价/新旧程度/学科分类/课程名/交易方式/校区/描述，最多 5 张实物图片。
3. **搜索与浏览**：按书名/作者/ISBN/课程名关键词搜索，按学科分类、新旧程度、价格区间筛选，按价格/发布时间/浏览量排序。
4. **求购信息**：买家发布求购（书名/期望价格/新旧要求/学科分类），卖家可一键联系求购者（自动创建会话）。
5. **交易沟通**：内置站内消息（文字 + 图片），会话未读数，支持标记书籍「已预约 / 已售出」。
6. **交易评价**：交易完成后买卖双方互评（好评/中评/差评 + 文字），用户主页展示好评率与历史评价，差评率 ≥40% 且评价数 ≥3 时标记风险提示。
7. **收藏与浏览历史**：收藏书籍、自动记录浏览历史、首页推荐同院系同学在售书籍。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | Vue 3 + TypeScript，使用 Vant 4 移动端组件库（适配手机端），Vite 构建工具 |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis（验证码存储，Redis 不可用时自动降级为内存存储） |
| 对象存储 | MinIO |
| 认证 | JWT + RBAC |
| 日志 | `log/slog` 结构化日志 |
| 参数校验 | `github.com/go-playground/validator/v10` |
| 接口文档 | 见下方完整 API 清单（README） |

## 项目目录结构

```
.
├── backend/
│   ├── cmd/server/main.go            # 入口：装配配置/依赖/启动
│   ├── internal/
│   │   ├── config/                   # 环境变量配置
│   │   ├── database/                 # 连接/自动迁移/种子数据
│   │   ├── model/                    # 每个实体一个文件
│   │   ├── dto/                      # 每个实体一个 DTO 文件
│   │   ├── repository/               # 每个实体一个 repository 文件
│   │   ├── service/                  # 每个实体一个 service 文件
│   │   ├── handler/                  # 每个实体一个 handler 文件
│   │   ├── router/                   # 每个实体一个路由注册文件
│   │   ├── middleware/               # auth / rbac / request_id / audit / recovery / rate_limiter / logger / error_handler
│   │   ├── constants/                # enums / error_codes / log_templates / messages
│   │   └── util/                     # jwt / logger / app_error / formatters / password / pagination / redis / minio
│   ├── migrations/                   # 数据库表结构参考 SQL
│   ├── Dockerfile
│   ├── go.mod / go.sum
├── frontend/
│   ├── src/
│   │   ├── api/                      # 每个实体一个 API 文件
│   │   ├── components/               # 共享组件：StatusBadge / ConditionTag / EmptyState / BookCard / ImageUploader / SubjectPicker / AppTabbar
│   │   ├── pages/                    # 每个模块一个页面
│   │   ├── stores/                   # auth / book / wish / conversation / evaluation
│   │   ├── hooks/                    # useAuth / usePagination
│   │   ├── utils/                    # request / format / upload
│   │   ├── constants/                # enums.ts（与后端枚举对应）
│   │   ├── router/                   # 路由守卫
│   │   └── types/                    # 类型定义
│   ├── Dockerfile
│   └── nginx.conf                    # 前端路由 + /api 反向代理
├── database/init.sql                 # MySQL 初始化（建库建账号，幂等）
├── docker-compose.yml
├── .env / .env.example
└── README.md
```

**严禁合并职责到单一文件**：每个核心实体拆分为 model / dto / repository / service / handler / router / constants 独立文件，前端按模块拆分 api / stores / pages / components，不把多个实体或所有页面写进同一个文件。

## 环境变量说明

详见 `.env.example`（含注释）。关键键：

- `COMPOSE_PROJECT_NAME=campusbooks`
- 数据库：`DB_NAME` / `DB_USER` / `DB_PASSWORD` / `DB_ROOT_PASSWORD`
- 认证：`JWT_SECRET`（生产环境务必更换）
- 端口：`FRONTEND_PORT=8011` / `BACKEND_PORT=3011` / `DB_PORT=44007` / `REDIS_PORT=46307` / `MINIO_PORT=9007` / `MINIO_CONSOLE_PORT=9008`
- 缓存：`REDIS_ADDR=redis:6379` / `REDIS_PASSWORD=`
- 对象存储：`MINIO_ENDPOINT=minio:9000` / `MINIO_ROOT_USER=minioadmin` / `MINIO_ROOT_PASSWORD=minioadmin`

## 本地开发（备选）

后端：

```bash
cd backend && go mod tidy && go run ./cmd/server
# 构建：go build ./...
# 测试：go test ./...
```

前端：

```bash
cd frontend && npm install && npm run dev
```

## 完整 API 清单（前缀 `/api/v1`，统一响应 `{"code":0,"message":"ok","data":...}`）

### 认证 Auth

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| POST | /auth/send-code | 发送邮箱验证码（返回 dev_code） | 否 |
| POST | /auth/register | 学号+邮箱+验证码注册 | 否 |
| POST | /auth/login | 学号/邮箱 + 密码登录 | 否 |

### 用户 User

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| GET | /users/me | 当前用户资料 | 是 |
| PUT | /users/me | 完善/更新资料 | 是 |
| PUT | /users/me/avatar | 更新头像 | 是 |
| GET | /users/me/stats | 当前用户统计（好评率/风险） | 是 |
| GET | /users | 用户列表（管理员） | 是 + admin |
| GET | /users/:id | 用户资料 | 是 |
| GET | /users/:id/stats | 用户统计（复用 userService.GetStats） | 否 |
| GET | /users/:id/evaluations | 用户收到的评价（复用 evaluationService.ListEvaluations） | 否 |

### 书籍 Book（核心状态机 on_sale → reserved → sold）

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| GET | /books | 搜索/筛选/排序（复用 bookService.ListBooks） | 否 |
| GET | /books/recommendations | 同院系推荐（复用 bookService.ListBooks 的 DTO 管线） | 是 |
| GET | /books/favorites | 我的收藏 | 是 |
| GET | /books/history | 我的浏览历史 | 是 |
| POST | /books | 发布书籍 | 是 |
| GET | /books/:id | 书籍详情（+浏览量 +浏览历史） | 否 |
| PUT | /books/:id | 编辑（仅卖家、在售） | 是 |
| DELETE | /books/:id | 下架（仅卖家） | 是 |
| POST | /books/:id/reserve | 标记已预约（on_sale→reserved，FOR UPDATE 事务） | 是 |
| POST | /books/:id/cancel-reserve | 取消预约（reserved→on_sale） | 是 |
| POST | /books/:id/sold | 确认售出（reserved/on_sale→sold） | 是 |
| POST | /books/:id/favorite | 收藏 | 是 |
| DELETE | /books/:id/favorite | 取消收藏 | 是 |

### 求购 Wish

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| GET | /wishes | 求购列表（关键词/学科/状态） | 否 |
| POST | /wishes | 发布求购 | 是 |
| GET | /wishes/:id | 求购详情 | 否 |
| PUT | /wishes/:id | 编辑（仅本人、进行中） | 是 |
| DELETE | /wishes/:id | 删除（仅本人） | 是 |
| POST | /wishes/:id/close | 关闭求购 | 是 |
| POST | /wishes/:id/contact | 卖家联系求购者（复用 conversationService.createConversation） | 是 |

### 会话与消息 Conversation / Message

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| GET | /conversations | 我的会话列表（含未读数） | 是 |
| POST | /conversations | 从书籍/求购创建会话 | 是 |
| GET | /conversations/:id | 会话详情 | 是 |
| POST | /conversations/:id/messages | 发送消息（文字/图片） | 是 |
| GET | /conversations/:id/messages | 消息列表 | 是 |
| PUT | /conversations/:id/read | 标记已读 | 是 |

### 评价 Evaluation

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| POST | /evaluations | 交易完成后互评（好评/中评/差评） | 是 |
| GET | /users/:id/evaluations | 用户收到的评价 | 否 |

### 文件与审计 File / Audit

| 方法 | 路径 | 说明 | 认证 |
| --- | --- | --- | --- |
| POST | /uploads | 上传图片（头像/书籍图，MinIO，≤5MB，jpg/png/webp/gif） | 是 |
| GET | /files/:key | 读取已上传文件 | 否 |
| GET | /audit-logs | 操作审计日志（管理员） | 是 + admin |

### 复用关系标注

- `bookService.ListBooks` 被「书籍列表」与「同院系推荐」复用（推荐页复用同一 DTO 转换管线）。
- `conversationService.createConversation` 被「从书籍发起会话」「从求购发起会话（wish contact）」两个接口复用。
- `userService.GetStats` 被「我的统计」与「他人主页统计」两个接口复用。
- `evaluationService.ListEvaluations` 被「评价接口响应」与「用户主页评价历史」复用。

## 横切关注点

1. **JWT 认证 + RBAC 权限**：数据库 `users.role` 字段 → `internal/middleware/auth.go`、`internal/middleware/rbac.go`、`internal/util/jwt.go` → 前端 `router/index.ts` 路由守卫（`requiresAuth` / `requiresAdmin`）与 `hooks/useAuth.ts` 按钮显隐（书籍详情「编辑/下架」仅卖家可见、「审计日志」仅 admin 可见）。
2. **操作审计日志**：`audit_logs` 表 → `internal/middleware/audit.go` 记录所有写操作 → service 埋点（登录、发布、状态流转等）→ 前端 `pages/Audit.vue` 审计页面（admin）。
3. **全局错误处理与请求追踪**：`internal/middleware/request_id.go`（X-Request-ID）→ `internal/middleware/error_handler.go` → `internal/util/app_error.go` → `internal/constants/error_codes.go` → 前端 `utils/request.ts` 统一拦截器（code!=0 提示、401 清理令牌）。

## 共享枚举 / 常量出现位置清单

以下业务枚举在后端 `internal/constants/enums.go` 统一定义，并贯穿模型 / DTO / service 状态机 / handler 校验 / 日志模板 / formatters 与前端 `frontend/src/constants/enums.ts`、筛选与状态徽标组件。

### 1. 书籍状态 BookStatus（on_sale / reserved / sold）

| 出现位置（后端） | 文件 |
| --- | --- |
| 枚举定义 | `internal/constants/enums.go` |
| 模型默认值 | `internal/model/book.go`（Status 字段默认 on_sale） |
| DTO 状态文本 | `internal/dto/book_dto.go`（FromBook） |
| 状态机 | `internal/service/book_service.go`（transition 事务 + FOR UPDATE） |
| handler 校验/路由动作 | `internal/handler/book_handler.go`（reserve/cancel/sold） |
| 日志模板 | `internal/constants/log_templates.go`（LogBookStatusChanged 等） |
| 错误码 | `internal/constants/error_codes.go`（CodeBookStatusInvalid/Conflict） |
| 格式化 | `internal/util/formatters.go`（FormatBookStatusText） |
| 种子数据 | `internal/database/seed.go` |
| 数据库脚本 | `backend/migrations/001_init.sql` |
| 前端枚举/文案/徽标 | `frontend/src/constants/enums.ts`（BookStatusText / BookStatusBadge） |
| 前端状态徽标 | `frontend/src/components/StatusBadge.vue` |
| 前端筛选/页面 | `frontend/src/pages/MyBooks.vue`（在售/已预约/已售出 tabs）、`BookDetail.vue`（操作按钮显隐） |

### 2. 新旧程度 Condition（brand_new / nine_new / seven_new / five_new）

| 出现位置（后端） | 文件 |
| --- | --- |
| 枚举定义 | `internal/constants/enums.go` |
| 模型 | `internal/model/book.go`、`internal/model/wish.go` |
| DTO 校验 | `internal/dto/book_dto.go`、`internal/dto/wish_dto.go`（oneof 标签） |
| 格式化 | `internal/util/formatters.go`（FormatConditionText） |
| 日志/种子 | `internal/database/seed.go` |
| 前端枚举 | `frontend/src/constants/enums.ts`（ConditionText / ConditionOptions） |
| 前端组件 | `frontend/src/components/ConditionTag.vue` |
| 前端表单/筛选 | `frontend/src/pages/PublishBook.vue`、`PublishWish.vue`、`BookList.vue` |

### 3. 学科分类 SubjectCategory（science / humanities / econ_management / art / other）

| 出现位置（后端） | 文件 |
| --- | --- |
| 枚举定义 | `internal/constants/enums.go` |
| 模型 | `internal/model/book.go`、`internal/model/wish.go` |
| DTO 校验 | `internal/dto/book_dto.go`、`internal/dto/wish_dto.go` |
| 格式化 | `internal/util/formatters.go`（FormatSubjectText） |
| 前端枚举 | `frontend/src/constants/enums.ts`（SubjectCategoryText / SubjectCategoryOptions） |
| 前端组件 | `frontend/src/components/SubjectPicker.vue` |
| 前端页面 | `Home.vue`（分类入口）、`BookList.vue`、`WishList.vue` |

### 4. 交易方式 TradeType（in_person / mail）、评价类型 EvaluationType（good / neutral / bad）、角色 Role（student / admin）、求购状态 WishStatus（open / closed）

- 后端定义：`internal/constants/enums.go`；格式化：`internal/util/formatters.go`（FormatTradeTypeText / FormatEvaluationTypeText / FormatRoleText / FormatWishStatusText）；DTO 校验：`internal/dto/*.go`；日志/种子：`internal/database/seed.go`。
- 前端定义：`frontend/src/constants/enums.ts`；组件：`StatusBadge.vue`（求购状态）、`BookCard.vue`（交易方式 tag）、`Profile.vue`（评价类型颜色）；RBAC 角色用于 `middleware/rbac.go` 与前端路由守卫。

## API 调用示例（curl）

```bash
# 1. 登录获取 JWT
TOKEN=$(curl -s -X POST http://localhost:3011/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"account":"zhang@campusbooks.local","password":"student123"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['data']['token'])")

# 2. 获取当前用户资料
curl -s http://localhost:3011/api/v1/users/me -H "Authorization: Bearer $TOKEN"

# 3. 发布书籍
curl -s -X POST http://localhost:3011/api/v1/books \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"计算机网络（第7版）","author":"谢希仁","price":18,"condition":"nine_new","subject_category":"science","trade_type":"in_person","images":[]}'

# 4. 搜索书籍（含 JWT 可选）
curl -s "http://localhost:3011/api/v1/books?keyword=高等数学&subject_category=science&sort=price_asc&page=1&page_size=10"

# 5. 标记书籍为已售出
curl -s -X POST http://localhost:3011/api/v1/books/1/sold -H "Authorization: Bearer $TOKEN"

# 6. 交易后评价
curl -s -X POST http://localhost:3011/api/v1/evaluations \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"to_user_id":3,"book_id":1,"type":"good","content":"交易顺利"}'

# 7. 管理员查看审计日志
ADMIN=$(curl -s -X POST http://localhost:3011/api/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"account":"admin@campusbooks.local","password":"admin123"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['data']['token'])")
curl -s "http://localhost:3011/api/v1/audit-logs?page=1&page_size=5" -H "Authorization: Bearer $ADMIN"
```

## Docker 部署说明

- 端口映射：前端 `${FRONTEND_PORT:-8011}:80`，后端 `${BACKEND_PORT:-3011}:8080`，MySQL `${DB_PORT:-44007}:3306`，Redis `${REDIS_PORT:-46307}:6379`，MinIO `${MINIO_PORT:-9007}:9000` + `${MINIO_CONSOLE_PORT:-9008}:9001`。
- 数据卷：`db_data`（MySQL）、`redis_data`（Redis AOF）、`minio_data`（对象存储）均为命名卷，`docker compose down` 不丢数据；`docker compose down -v` 会清空数据。
- 健康检查：MySQL `mysqladmin ping`、Redis `redis-cli ping`、MinIO TCP 探测、后端 `wget /healthz`；后端 `depends_on` 数据库/Redis/MinIO `service_healthy`。
- 常见问题：
  - 端口冲突：修改 `.env` 中 `FRONTEND_PORT`/`BACKEND_PORT`/`DB_PORT` 等即可。
  - 忘记演示账号：登录页有提示，或查看 `backend/internal/database/seed.go`。
  - 邮箱验证码：开发环境 `POST /auth/send-code` 直接在响应返回 `dev_code`，前端注册页会自动填入。
  - 镜像拉取慢：配置 Docker 镜像加速后重试 `docker compose up -d --build`。

## License

MIT
