package main

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/v-mars/zlog"
)

var Dlog hlog.FullLogger

func InitLog(LogFileName, logFormat string, level hlog.Level) {
	var rc = zlog.GetDefaultRotateConfig(
		LogFileName,
		zlog.WithMaxSize(10),    // A file can be up to 10M.
		zlog.WithMaxAge(7),      // A file can exist for a maximum of 7 days.
		zlog.WithMaxBackups(10), // Save up to 10 files at the same time.
	)
	var lo = zlog.New(
		zlog.WithLevel(level),
		zlog.WithFormat(zlog.GetLogFormat(logFormat)),
		zlog.WithRotation(rc, nil),
		zlog.CallerWithSkipFrameCount(3),
		//zlog.WithOutput(os.Stdout),
	)
	hlog.SetLogger(lo)
	hlog.SetLevel(level)
	Dlog = lo
	lo.Debug("Init log")
}

func main() {
	InitLog(
		"log.log",
		"text",
		hlog.LevelInfo,
	)
	hlog.Info("Hello, World!")
}
