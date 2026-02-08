# zlog 性能优化指南

## 性能特点分析

根据基准测试，zlog 在不同使用场景下表现出不同的性能特征：

### 1. 格式性能差异显著

- **JSON 格式**: 最佳性能，适合生产环境
- **Console 格式**: 由于格式化开销，性能相对较低

## 性能优化建议

### 1. 生产环境最佳实践

```go
// 推荐：生产环境使用 JSON 格式
logger := zlog.New(
    zlog.WithFormat(zlog.JSONFormat),  // JSON 格式性能最优
    zlog.WithLevel(hertzlog.LevelInfo), // 合适的日志级别
    zlog.WithOutput(os.Stdout),         // 直接输出到 stdout
)

// 不推荐：生产环境使用 Console 格式（性能开销大）
logger := zlog.New(
    zlog.WithFormat(zlog.ConsoleFormat), // 性能较低
    zlog.WithOutput(file),
)
```

### 2. 开发环境优化

```go
// 开发环境可使用 Console 格式，但建议禁用过多字段
logger := zlog.New(
    zlog.WithFormat(zlog.ConsoleFormat),
    zlog.WithLevel(hertzlog.LevelDebug),
)

// 使用上下文日志时注意字段数量
logger.CtxInfof(ctx, "Request processed",
    "userID", userID,
    "action", action,
    // 避免添加过多字段，影响性能
)
```

### 3. 高并发场景优化

```go
// 对于高并发场景，考虑预创建 logger
var sharedLogger = zlog.New(
    zlog.WithFormat(zlog.JSONFormat),
    zlog.WithLevel(hertzlog.LevelInfo),
)

// 避免频繁创建 logger 实例
func handleRequest() {
    sharedLogger.Info("Handling request") // 重复使用实例
}
```

### 4. 批量日志优化

```go
// 减少不必要的日志记录
func processItems(items []Item) {
    // 避免循环内记录大量日志
    if len(items) > 1000 {
        logger.Info("Starting batch processing",
            "count", len(items))
    }

    for i, item := range items {
        // 只在必要时记录详细信息
        if i%1000 == 0 { // 每1000条记录一次进度
            logger.Info("Processing progress",
                "processed", i,
                "total", len(items))
        }

        processItem(item)
    }
}
```

### 5. 配置优化

```go
// 针对不同场景选择合适配置
func createProductionLogger() *zlog.ZLogger {
    return zlog.New(
        zlog.WithFormat(zlog.JSONFormat),
        zlog.WithLevel(hertzlog.LevelInfo),
        zlog.WithRotation(&zlog.RotateConfig{
            MaxSize:    100,    // 较大的单文件大小
            MaxBackups: 5,      // 适量备份文件
            MaxAge:     30,     // 适当的保留天数
            Compress:   true,   // 启用压缩
        }),
    )
}

// 开发环境可启用更多调试信息
func createDevelopmentLogger() *zlog.ZLogger {
    return zlog.New(
        zlog.WithFormat(zlog.ConsoleFormat),
        zlog.WithLevel(hertzlog.LevelDebug),
    )
}
```

## 性能监控

```go
// 实现日志性能监控
type PerformanceLogger struct {
    logger *zlog.ZLogger
    metrics MetricsCollector
}

func (p *PerformanceLogger) Info(msg string, fields ...interface{}) {
    start := time.Now()
    p.logger.Info(msg, fields...)
    duration := time.Since(start)

    // 记录日志性能指标
    p.metrics.RecordLogLatency(duration)
}
```

## 最佳实践摘要

| 场景 | 推荐配置 | 理由 |
|------|----------|------|
| 生产环境 | JSON 格式 | 最佳性能 |
| 开发环境 | Console 格式 | 可读性强 |
| 高并发 | 预创建实例 | 避免重复初始化 |
| 微服务 | 集成 OTel | 便于追踪 |
| 日志分析 | JSON 格式 | 易解析 |
| 低延迟 | 避免复杂格式化 | 减少开销 |

遵循这些优化建议，可以充分发挥 zlog 的性能优势，同时获得丰富的功能特性。