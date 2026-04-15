package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//var logger *zap.sugaredlogger // 定义日志打印全局变量
//var logger *zap.Logger // 句柄, 对内不对外

type LogConfigs struct {
	LogLevel          string // 日志打印级别 debug  info  warning  error
	LogFormat         string // 输出日志格式	logfmt, json
	LogPath           string // 输出日志文件路径
	LogFileName       string // 输出日志文件名称
	LogFileMaxSize    int    // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    // 【日志分割】日志备份文件最多数量

	LogMaxAge   int  // 日志保留时间，单位: 天 (day)
	LogCompress bool // 是否压缩日志
	LogStdout   bool // 是否输出到控制台
}

type Zlogger = zap.Logger

// 初始化 logger
func initlogger(conf *LogConfigs) (logger *Zlogger, err error) {

	LogLevel := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
	writesyncer, err := getLogWriter(conf) // 日志文件配置 文件位置和切割
	if err != nil {
		return
	}
	encoder := getEncoder(conf)          // 获取日志输出编码
	level, ok := LogLevel[conf.LogLevel] // 日志打印级别
	if !ok {
		level = LogLevel["info"]
	}
	core := zapcore.NewCore(encoder, writesyncer, level)
	logger = zap.New(core, zap.AddCaller()) //  zap.addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33

	logger = logger.WithOptions(zap.AddCallerSkip(1)) //为了封装外部访问接口,跳到上一层

	//logger = logger.sugar()
	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	zap.ReplaceGlobals(logger)

	//zap.L().Debug("") //输出带级别的
	//zap.S().Debugf("")//输出不带级别的
	return logger, nil
}

// 编码器(如何写入日志)
func getEncoder(conf *LogConfigs) zapcore.Encoder {

	encoderconfig := zap.NewProductionEncoderConfig()
	encoderconfig.EncodeTime = zapcore.ISO8601TimeEncoder   // looger 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderconfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 info debug error

	//	encoderconfig.EncodeCaller = zapcore.FullCallerEncoder
	//	encoderconfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if conf.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderconfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderconfig) // 以logfmt格式写入
}

// 获取日志输出方式  日志文件 控制台
func getLogWriter(conf *LogConfigs) (zapcore.WriteSyncer, error) {
	// 判断日志路径是否存在，如果不存在就创建
	if exist := isexist(conf.LogPath); !exist {
		if conf.LogPath == "" {
			conf.LogPath = string(*DefaultLogPath)
		}
		if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
			conf.LogPath = string(*DefaultLogPath)
			if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
				conf.LogStdout = false
			}
		}
	}

	// 日志文件 与 日志切割 配置
	lumberjacklogger := &lumberjack.Logger{
		Filename:   filepath.Join(conf.LogPath, conf.LogFileName), // 日志文件路径
		MaxSize:    conf.LogFileMaxSize,                           // 单个日志文件最大多少 mb
		MaxBackups: conf.LogFileMaxBackups,                        // 日志备份数量
		MaxAge:     conf.LogMaxAge,                                // 日志最长保留时间
		Compress:   conf.LogCompress,                              // 是否压缩日志
	}
	if conf.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberjacklogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到控制台
		return zapcore.AddSync(os.Stdout), nil
	}
}

// 判断文件或者目录是否存在
func isexist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
