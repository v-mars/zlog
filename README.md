# zlog

## 介绍
zlog是一个灵活且高性能的Go日志库，支持与Hertz框架的hlog集成，并提供日志轮转功能。该库基于zerolog构建，提供了丰富的日志功能和良好的性能。

## 功能特性
- 兼容Hertz的hlog接口
- 支持多种日志级别（Trace, Debug, Info, Notice, Warn, Error, Fatal）
- 格式化日志输出
- 上下文日志支持（context-aware logging）
- 日志轮转（基于lumberjack）
- 动态调整日志级别和输出目标
- 高性能（基于zerolog）

## 安装

```bash
go mod init your-project
go get github.com/cloudwego/hertz/pkg/common/hlog
go get github.com/v-mars/oceanlog
go get gopkg.in/natefinch/lumberjack.v2
```

## 快速开始

### 基本用法

```go
package main

import (
    "zlog"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
    logger := zlog.New(
        zlog.WithLevel(hlog.LevelDebug),
    )

    logger.Info("This is an info message")
    logger.Debugf("Formatted message: %s, %d", "hello", 42)
}
```

### 使用日志轮转

```go
package main

import (
    "github.com/v-mars/zlog"
)

func main() {
    // 创建默认的轮转配置
    config := zlog.GetDefaultRotateConfig("app.log")

    // 或自定义配置
    config = &zlog.RotateConfig{
        Filename:   "app.log",
        MaxSize:    50,  // 50MB
        MaxBackups: 5,   // 保留5个备份
        MaxAge:     30,  // 保留30天
        Compress:   true, // 压缩备份文件
    }

    rotatingLogger := zlog.NewRotatingLogger(config)
    rotatingLogger.Info("This goes to a rotating log file")
}
```

### 与Hertz hlog集成

```go
package main

import (
    "github.com/v-mars/zlog"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
    zLogger := zlog.New(zlog.WithLevel(hlog.LevelDebug))

    // 设置为默认的hlog后端
    zlog.SetAsHlogDefault(zLogger)

    // 现在所有的hlog调用都会使用zlog
    hlog.Info("This message uses zlog as backend")
    hlog.Errorf("Error occurred: %v", err)
}
```

### 上下文日志

```go
package main

import (
    "context"
    "github.com/v-mars/zlog"
)

func main() {
    logger := zlog.New()
    ctx := context.WithValue(context.Background(), "request_id", "12345")

    logger.CtxInfof(ctx, "Processing request: %s", "12345")
}
```

## 接口兼容性

zlog完全兼容以下接口：
- `hlog.Logger` - 基础日志方法
- `hlog.FormatLogger` - 格式化日志方法
- `hlog.CtxLogger` - 上下文日志方法
- `hlog.Control` - 控制方法（设置级别、输出）

## 配置选项

### 日志轮转配置

| 选项 | 描述 |
|------|------|
| Filename | 日志文件名 |
| MaxSize | 单个文件最大大小（MB） |
| MaxBackups | 保留的最大备份文件数 |
| MaxAge | 保留日志文件的最大天数 |
| Compress | 是否压缩备份文件 |
| LocalTime | 是否使用本地时间戳 |

### 创建不同类型的logger

```go
// 创建控制台logger
consoleLogger := zlog.New(zlog.WithOutput(os.Stdout))

// 创建带级别的logger
levelLogger := zlog.New(zlog.WithLevel(hlog.LevelDebug))

// 创建同时具有轮转功能的logger
rotatingLogger := zlog.New(zlog.WithRotation(config))
```

## 性能优化

zlog基于zerolog构建，在性能方面有以下特点：
- 零分配日志记录（在热路径上）
- 结构化日志输出
- 高效的JSON序列化
- 并发安全

## 错误处理

大多数zlog操作都是无错误的，但在某些情况下（如手动轮转时），可能需要检查错误：

```go
rotatingLogger := zlog.NewRotatingLogger(config)
err := rotatingLogger.Rotate() // 手动触发轮转
if err != nil {
    // 处理轮转错误
}
```

## 维护者

Ocean

## 贡献

欢迎提交Issue和PR！

## 许可证

MIT