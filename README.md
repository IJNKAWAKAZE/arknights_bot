# Arknights Telegram Bot

[![Go](https://github.com/IJNKAWAKAZE/arknights_bot/actions/workflows/go.yml/badge.svg)](https://github.com/IJNKAWAKAZE/arknights_bot/actions/workflows/go.yml)

一个基于 [Telegram Bot API](https://core.telegram.org/bots/api) 的明日方舟（Arknights）群组机器人，使用 Go 语言编写。

本机器人设计用于群组 [明日方舟 / Arknights 中文交流](https://t.me/ArknightsZH)，包含**入群验证、森空岛签到、游戏数据查询、群组小游戏、抽奖、寻访模拟**等丰富功能，并可通过 Web 模板 + Playwright 截图生成精美的查询图片。

## 功能特性

- **入群验证与群管理**：入群申请验证、自定义群规/欢迎信息、请求验证模式切换、自定义群标签、消息举报。
- **森空岛账号服务**：绑定/解绑角色、每日签到与自动签到、签到通知模式、理智提醒（阈值可调、动态调度）。
- **游戏数据查询**：干员/皮肤/敌人查询（支持中文、拼音、多音字模糊搜索）、角色名片、基建状态、仓库、抽卡记录导入导出、CDK 兑换。
- **群组小游戏**：猜干员（立绘/剪影/看脚模式）、寻访模拟、公招计算（图片识别）。
- **抽奖功能**：发起/停止/结束抽奖、报名与中奖详情。
- **定时任务**：B站动态推送、生日祝福、自动签到、理智检查、数据源每周自动更新、延迟删消息、每日寻访次数重置。
- **图片渲染**：基于本地 Web 服务（Gin）渲染模板，由 Playwright 截图生成查询图片。

## 指令列表

### 私聊指令

| 指令 | 说明 |
| --- | --- |
| `/start` | 查看使用说明 |
| `/bind` | 绑定森空岛账号（支持多账号与多角色） |
| `/unbind` | 解绑角色 |
| `/reset_token` | 重设森空岛 Token |
| `/import_gacha` | 导入抽卡记录 |
| `/export_gacha` | 导出抽卡记录 |
| `/cancel` | 取消当前操作 |

### 普通指令

| 指令 | 说明 |
| --- | --- |
| `/help` | 查看使用说明 |
| `/ping` | 存活测试 |
| `/tag` | 自定义群标签 |
| `/sign` | 森空岛签到（`auto` 开启自动签到 / `stop` 关闭 / `notify_all` / `notify_fail` / `notify_success` 设置通知模式） |
| `/ap` | 理智提醒（`on` / `off` / `thr [1-100]` 设置阈值） |
| `/state` | 查看账号当前状态 |
| `/box` | 查看干员列表（默认 6 星，支持 `all` / `1,2` 等参数） |
| `/box_detail` | 干员详情 |
| `/box_summary` | 干员信息汇总 |
| `/missing` | 查看未获取干员 |
| `/card` | 查看角色名片 |
| `/base` | 查看基建状态 |
| `/gacha` | 查看抽卡分析 |
| `/operator` | 干员信息查询（不输入名称显示搜索按钮） |
| `/skin` | 干员皮肤查询 |
| `/enemy` | 敌人信息查询 |
| `/report` | 回复群消息进行举报（自动 @ 所有管理员） |
| `/quiz` | 猜干员小游戏（`h` 黑白剪影 / `ex` 看脚模式） |
| `/redeem [code]` | 兑换 CDK |
| `/headhunt` | 寻访模拟 |
| `/calendar` | 活动日历 |
| `/depot` | 查看我的仓库 |
| `/join_lottery` | 报名抽奖 |
| `/lottery_detail` | 查看抽奖详情 |

### 群管理员指令

| 指令 | 说明 |
| --- | --- |
| `/news` | 开启/关闭 B站动态推送 |
| `/birthday` | 开启/关闭生日推送 |
| `/request_mode` | 切换群验证模式 |
| `/quiz on/off` | 开启/关闭猜干员小游戏 |
| `/headhunt on/off` | 开启/关闭寻访模拟 |
| `/reg` | 回复一条消息设置为群规 |
| `/welcome` | 设置入群欢迎信息 |
| `/clear [key]` | 根据 key 删除 Redis 缓存 |
| `/start_lottery` | 发起抽奖 |
| `/stop_lottery` | 停止报名 |
| `/end_lottery` | 结束抽奖 |
| `/lottery` | 查看抽奖信息 |

### 机器人拥有者指令

| 指令 | 说明 |
| --- | --- |
| `/update` | 更新数据源（首次运行需执行一次，之后每周五凌晨自动更新）；开启本地素材缓存后会同时增量下载干员素材 |
| `/sign_all` | 全部账号签到 |
| `/kill` | 关闭机器人 |

## 部署文档

### 环境要求

- [Go](https://go.dev) 1.22+
- [MySQL](https://www.mysql.com/) 5.7+ 或 [MariaDB](https://mariadb.org)（推荐）
- [Redis](https://redis.io) 5.0+
- [FFmpeg](https://ffmpeg.org/)（用于 B站动态视频转码）
- 中文字体（用于生成图片，Docker 镜像已内置 `fonts-noto-*`）

> 图片查询功能依赖 Playwright，首次启动会自动安装浏览器内核，请确保网络可访问 Playwright 下载源。

### 方式一：直接编译运行

1. 克隆仓库并进入项目目录：

   ```shell
   git clone https://github.com/IJNKAWAKAZE/arknights_bot
   cd arknights_bot
   ```

2. 复制配置文件并修改为你的配置：

   ```shell
   cp arknights.example.yaml arknights.yaml
   ```

   至少需要修改以下内容：

   - `bot.token`：Telegram Bot Token（通过 [@BotFather](https://t.me/BotFather) 创建）
   - `bot.owner`：机器人拥有者的 Telegram userId
   - `mysql.dsn`：数据库连接串
   - `redis.addr` / `redis.pwd`：Redis 地址与密码

3. 导入数据库结构：

   ```shell
   mysql -u root -p arknights < arknights.sql
   ```

4. 编译并运行：

   ```shell
   cd src
   go build -o ../arknights_bot .
   cd ..
   ./arknights_bot
   ```

   需要指定其他位置的配置文件时，使用 `-config` 参数：

   ```shell
   ./arknights_bot -config /path/to/arknights.yaml
   ```

   启动成功后，向机器人发送 `/update` 初始化数据源。

### 方式二：Docker Compose 部署

仓库自带的 `docker-compose.yaml` 提供 `mysql`（MariaDB）、`redis`、`arkbot` 三个服务。Dockerfile 只负责打包运行环境，需要先编译好 Linux 版二进制放到仓库根目录，再一条命令部署：

```shell
cd src
GOOS=linux GOARCH=amd64 go build -o ../arknights_bot .
cd ..
docker compose up -d --build
```

> 在 Linux 服务器上部署时，也可以直接在服务器上编译：`cd src && go build -o ../arknights_bot .`，再执行上面的 `docker compose` 命令。

> 部署前请修改 `docker-compose.yaml` 中 mysql/redis 的环境变量（`DATABASE_NAME`、`DATABASE_USER`、`DATABASE_PASSWORD`、`ROOT_PASSWORD`、`REDIS_PASSWORD`），并准备 `arknights.yaml`。

容器挂载了以下目录/文件，修改后重启容器即可生效：

- `./arknights.yaml`：配置文件
- `./assets`、`./template`：静态资源与图片模板
- `./assets_local`：本地素材缓存目录（开启 `local_assets.enable` 后使用，全量约 1GB 出头，随皮肤数量增长）
- `./logs`：入群审计日志（`join.log`）
- Playwright 浏览器缓存（持久化卷）

### 配置文件说明

`arknights.yaml` 主要配置项：

| 配置项 | 说明 |
| --- | --- |
| `bot.name` | 机器人用户名 |
| `bot.token` | 机器人 Token |
| `bot.owner` | 拥有者 userId |
| `bot.msg_del_delay` | 延迟删除消息的秒数 |
| `bot.debug` | 是否开启调试模式 |
| `api.*` | 各数据源接口地址（PRTS Wiki、B站、森空岛等） |
| `mysql.dsn` | MySQL/MariaDB 连接串 |
| `redis.addr` / `pwd` / `db` / `pool_size` | Redis 配置 |
| `http.host` / `port` | Web 服务监听地址与端口（用于生成查询图片，默认仅监听 `127.0.0.1`，如需公网访问将 `host` 改为 `0.0.0.0`） |
| `telegraph.token` | Telegraph 页面 token（用于举报内容） |
| `headhunt.*` | 寻访模拟卡池配置（每日次数、UP 与各星级干员池） |
| `recruit.*` | 公招标签与缺失干员映射 |
| `enemy_name` | 敌人名称映射 |
| `ignore_birthday` | 忽略生日推送的干员列表 |
| `ad` | bio 广告词过滤，命中即拒绝申请/踢出（去空格、转小写后模糊匹配） |
| `proxy` | 森空岛请求代理，不填则不使用 |
| `local_assets.*` | 本地素材缓存（开关、目录、缺失回退等），详见下文「本地素材缓存」 |

### 本地素材缓存（可选）

图床（`media.prts.wiki`、`web.hycdn.cn`）偶发不可用或网络抖动时，截图会因为素材没加载出来而显示占位图。开启本地素材缓存后：

- 执行 `/update`（含每周定时更新）时，会把干员的头像、半身像、立绘、皮肤等素材**增量**下载到本地磁盘（图床素材 + 游戏内皮肤头像/半身像），已存在的文件不会重复下载，下载失败会自动重试；
- 渲染图片时优先使用本地文件，本地缺失时才回退到远程图床；渲染过程中遇到缺失的素材也会自动补齐到本地。

配置示例：

```yaml
local_assets:
  enable: true
  dir: ./assets_local
  remote_fallback: true
  concurrency: 4
  retry: 3
```

> **注意**：开启后会占用一定磁盘空间，全量素材约 1~2GB（干员立绘与皮肤占大头，随皮肤数量增长），首次下载耗时也较长（视网络情况从十几分钟到一小时不等），请确认磁盘剩余空间充足。

- `enable`：开关，默认 `false`，修改后重启或热更新配置即可生效；
- `dir`：素材目录，相对路径以程序运行目录为基准，Docker 部署对应容器内 `/root/assets_local`；建议不要放在 `assets` 里，避免和仓库自带的静态资源混在一起；
- `remote_fallback`：本地缺失且补齐失败时是否回退远程图床，默认 `true`；关闭后缺图会直接显示占位图；
- `concurrency` / `retry`：下载并发数与单张图片的最大尝试次数，默认 `4` / `3`。

素材按 `<dir>/<图床域名>/<原始路径>` 存放，需要清理时直接删掉目录即可，下次 `/update` 会重新下载。

## 开发文档

### 项目结构

```
.
├── src/                   # Go 源码
│   ├── cmd/               # 入口，初始化与启动
│   ├── config/            # 配置加载与 DB/Redis/机器人初始化
│   ├── core/
│   │   ├── bot/           # Telegram 指令/回调注册
│   │   ├── cron/          # 定时任务
│   │   └── web/           # Gin Web 服务（帮助页、查询图片）
│   ├── plugins/           # 功能插件
│   │   ├── account/       # 账号绑定
│   │   ├── apremind/      # 理智提醒
│   │   ├── arknightsnews/ # B站动态推送
│   │   ├── birthday/      # 生日推送
│   │   ├── datasource/    # 数据源更新
│   │   ├── enemy/         # 敌人查询
│   │   ├── gatekeeper/    # 入群验证
│   │   ├── lottery/       # 抽奖
│   │   ├── operator/      # 干员查询
│   │   ├── player/        # 玩家数据（box/state/gacha 等）
│   │   ├── sign/          # 签到
│   │   ├── skin/          # 皮肤查询
│   │   ├── skland/        # 森空岛 API 封装
│   │   └── system/        # 系统功能（help/quiz/recruit 等）
│   └── utils/             # 基础工具库
│       ├── model/         # 数据模型
│       ├── cache/         # Redis 封装
│       ├── repo/          # 数据库查询
│       ├── media/         # 截图/图片/OCR
│       ├── search/        # 干员/敌人搜索
│       ├── localassets/   # 图床素材本地缓存
│       └── hashutil/      # 哈希/随机字符串
├── assets/                # 静态资源
├── template/              # 图片生成模板（Go template + HTML）
├── arknights.sql          # 数据库初始化脚本
├── arknights.example.yaml # 示例配置文件
└── docker-compose.yaml    # Docker 编排
```

### 技术栈

- **语言**：Go 1.22
- **Web**：Gin（本地查询图片渲染）
- **配置**：Viper（支持热更新）
- **数据库**：GORM + MySQL/MariaDB
- **缓存**：go-redis
- **Telegram**：`github.com/ijnkawakaze/telegram-bot-api`
- **截图**：playwright-go
- **视频**：ffmpeg-go

### 本地开发

**Visual Studio Code**

1. 安装 [Go 扩展](https://marketplace.visualstudio.com/items?itemName=golang.Go)。
2. 安装语言服务器与调试器：

   ```shell
   go install -v golang.org/x/tools/gopls@latest
   go install -v github.com/go-delve/delve/cmd/dlv@latest
   ```

3. 打开项目，在 `src` 目录下运行 `go build` 或直接调试 `src/arknights_bot.go`。

**GoLand**

1. 安装 [GoLand](https://www.jetbrains.com/go/)。
2. 打开项目，将 `src` 设为模块根目录后直接运行。

### 定时任务

| 执行时间 | 任务 |
| --- | --- |
| 每 30 秒 | B站动态推送检查 |
| 每周五 02:33 | 更新数据源 |
| 每日 08:00 | 生日祝福 |
| 每日 01:00 | 自动签到 |
| 每日 01:30 | 理智检查 |
| 每 1 秒 | 清理延迟删除的消息 |
| 每日 00:00 | 重置寻访模拟次数 |
| 每 1 分钟 | 检查抽奖报名是否停止 |

消息处理限制为普通任务同时执行 6 个、验证和部分管理任务同时执行 2 个；普通任务每个会话最多保留 32 个（含正在执行的任务）、全局上限 256 个。空闲队列自动回收，超限的指令会尽量回复繁忙提示。抽奖管理与报名保持在同一队列，保证操作顺序。入群申请、成员变更这类丢失后无法补救的更新不做容量拒绝，只限制并发数，积压超过 256 条时记录日志，避免验证事件被静默丢弃。

延迟删消息使用 Redis 有序集合 `msgObjects:due`，每秒最多处理 100 条到期任务。升级时会自动分批迁移旧的 `msgObjects` 列表，无需手动清理 Redis。临时删除失败会保留任务重试；Telegram 明确表示消息已不存在或不能再删时按完成处理；到期超过 24 小时仍未删掉的任务会被放弃，避免永久失败的任务被无限重试。领取后未确认的任务会在 5 分钟后重新可用。

退出时先停止接收新任务、取消待执行任务和定时等待，再等待运行中的处理完成，最后关闭 Web、浏览器和数据库连接。清理步骤不会因为某一步失败或超时而中断，避免浏览器、数据库被跳过；退出最多等待 30 秒，之后仍会给资源释放留 10 秒。超时会记录错误并以非零状态退出。批量签到期间再次启动批量签到会提示已有任务在执行。

### 数据源

- 干员、皮肤、敌人等游戏数据来自 [PRTS Wiki](https://prts.wiki/)，由 `/update` 或定时任务抓取并缓存至 Redis。
- 玩家数据来自森空岛 API，需要在私聊中绑定账号。

## 致谢

- 游戏数据来源：[PRTS Wiki](https://prts.wiki/)

## License

本项目使用 [GPL-3.0](LICENSE) 许可证。
