package logging

import (
	"fmt"
	"os"

	"github.com/ems-go/logging/zlog"
)

const (
	LOG_FILE_NAME       string = "./tmp.log" // 日志文件名称
	LOG_FILE_SIZE       int    = 10          // 日志文件大小
	LOG_FILE_MAXAGE     int    = 30          // 日志文件最大保存时间
	LOG_FILE_MAXBACKUPS        = 7           // 日志文件最大保存个数
)

const (
	// DefaultLogPath           = "."     // 【默认】日志文件路径
	DefaultLogLevel          = "debug" // 【默认】日志打印级别 debug  info  warning  error
	DefaultLogFileMaxsize    = 5       // 【日志分割】  【默认】单个日志文件最多存储量 单位(mb)
	DefaultLogFileMaxBackups = 3       // 【日志分割】  【默认】日志备份文件最多数量
	LogMaxage                = 28      // 【默认】日志保留时间，单位: 天 (day)
	LogCompress              = false   // 【默认】是否压缩日志
	LogStdout                = true    // 【默认】是否输出到文件
)

var DefaultLogFilename = fmt.Sprintf("%s.log", os.Args[0]) // 【默认】日志文件名称

type Zlog struct {
	logger *Zlogger
}

var zstd Zlog

func Init() {

	conf := LogConfigs{
		LogLevel:          DefaultLogLevel,          // 输出日志级别 "debug" "info" "warn" "error"
		LogFormat:         "logfmt",                 // 输出日志格式 logfmt, json
		LogPath:           string(*DefaultLogPath),  // 输出日志文件位置
		LogFileName:       DefaultLogFilename,       // 输出日志文件名称
		LogFileMaxSize:    DefaultLogFileMaxsize,    // 输出单个日志文件大小，单位MB
		LogFileMaxBackups: DefaultLogFileMaxBackups, // 输出最大日志备份个数
		LogMaxAge:         LogMaxage,                // 日志保留时间，单位: 天 (day)
		LogCompress:       LogCompress,              // 是否压缩日志
		LogStdout:         bool(*dbgFileValue),      // 是否输出到控制台
	}

	zstd = NewZlog(conf)

}

func NewZlog(conf LogConfigs) (z Zlog) {
	logger, _ := initlogger(&conf)
	z.logger = logger
	z.logger = logger.WithOptions(zlog.AddCallerSkip(1))
	return z
}

func (l *Zlog) Debugf(fmt string, arg ...interface{}) {
	l.logger.Sugar().Debugf(fmt, arg...)
}

func (l *Zlog) Infof(fmt string, args ...interface{}) {
	l.logger.Sugar().Infof(fmt, args...)
}

func (l *Zlog) Warningf(fmt string, arg ...interface{}) {
	l.logger.Sugar().Warnf(fmt, arg...)
}

func (l *Zlog) Errorf(fmt string, arg ...interface{}) {
	l.logger.Sugar().Errorf(fmt, arg...)
}

func (l *Zlog) Fatal(arg ...interface{}) {
	l.logger.Sugar().Fatal(arg...)
}

func (l *Zlog) Debug(arg ...interface{}) {
	l.logger.Sugar().Debug(arg...)
}

func (l *Zlog) Info(arg ...interface{}) {
	l.logger.Sugar().Info(arg...)
}

func (l *Zlog) Warning(arg ...interface{}) {
	l.logger.Sugar().Warn(arg...)
}

func (l *Zlog) Error(arg ...interface{}) {
	l.logger.Sugar().Error(arg...)
}

// 开箱即用
func FlushLogs() {
	zstd.logger.Sync()
}

// func ZlogDefault() *Zlogger {
// 	return zstd.logger
// }

/* debug */
func Debugf(fmt string, args ...interface{}) {
	zstd.Debugf(fmt, args...)
}
func Debug(args ...interface{}) {
	zstd.Debug(args...)
}

/* infof */
func Infof(fmt string, args ...interface{}) {
	zstd.Infof(fmt, args...)
}
func Info(args ...interface{}) {
	zstd.Info(args...)
}

/* Warning */
func Warningf(fmt string, args ...interface{}) {
	zstd.Warningf(fmt, args...)
}
func Warning(args ...interface{}) {
	zstd.Warning(args...)
}

/* Error */
func Errorf(fmt string, args ...interface{}) {
	zstd.Errorf(fmt, args...)
}
func Error(args ...interface{}) {
	zstd.Error(args...)
}

func Fatal(args ...interface{}) {
	zstd.Fatal(args...)
}
