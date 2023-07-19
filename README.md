# ProjectTemp

## 1. Introduction

使用wire依赖注入代码生成，提供http、websocket服务。

## 2. Quick Start

```bash
# 1. 安装依赖
make init

# 2. 代码生成
make wire

# 3. 编译可执行文件
make build
```

## 3. Project Structure

```bash
.
├── Makefile    # makefile
├── README.md   # readme
├── bin         # 编译的可执行文件生成目录
├── cmd         # 项目代码入口
│   ├── main.go # main文件
│   ├── wire.go
│   └── wire_gen.go
├── config      # 配置文件
│   └── config.yaml
├── go.mod      # go mod 文件
├── go.sum
├── internal    # 内部包
│   ├── handler # handler，入口
│   │   ├── handler.go
│   │   └── test.go
│   ├── router # 路由
│   │   ├── http.go
│   │   ├── router.go
│   │   └── ws.go
│   └── service # service
│       ├── service.go
│       └── test.go
└── pkg           # 公共包
    ├── config    # 配置文件解析
    │   └── config.go
    ├── es        # es相关
    │   └── es.go
    ├── gorm      # mysql相关
    │   ├── gorm.go
    │   └── mysql.go
    ├── kafka     # kafka相关
    │   └── kafka.go
    ├── log       # 日志相关
    ├── middeware # 中间件
    │   └── cors
    │       └── cors.go
    ├── redis     # redis相关
    │   └── redis.go
    ├── server    # 服务相关
    │   └── server.go
    ├── util      # 工具类
    │   └── util.go
    └── ws        # websocket相关

```