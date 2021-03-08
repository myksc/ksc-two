package common

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

const (
	LogNameAccess = "access"
	LogNameServer = "server"
)

const (
	txtLogNormal    = "normal"
	txtLogWarnFatal = "warnfatal"
	txtLogStdout    = "stdout"
)

// 业务日志输出logger
var (
	ServerLogger *zap.Logger
	AccessLogger *zap.Logger
)

type loggerConfig struct {
	ZapLevel zapcore.Level

	// 以下变量仅对开发环境生效
	Stdout   bool
	Log2File bool
	Path     string

	RotateUnit  string
	RotateCount int
	RotateSwitch bool
}

// 全局配置 仅限Init函数进行变更
var logConfig = loggerConfig{
	ZapLevel: zapcore.InfoLevel,
	Log2File: true,
	Path:     "./log",
	RotateUnit:   "h",
	RotateCount:  24,
	RotateSwitch: false,
}

func InitZap(){
	ServerLogger = GetLogger()
	AccessLogger = GetAccessLogger()
}

// 业务日志(分为成功和失败)
func GetLogger() (l *zap.Logger) {
	if ServerLogger == nil {
		ServerLogger = newLogger(LogNameServer).WithOptions(zap.AddCallerSkip(1))
	}
	return ServerLogger
}

// 成功日志
func GetAccessLogger() (l *zap.Logger) {
	if AccessLogger == nil {
		AccessLogger = newLogger(LogNameAccess)
	}
	return AccessLogger
}

// NewLogger 新建Logger，每一次新建会同时创建x.log与x.log.wf (access.log 不会生成wf)
func newLogger(name string) *zap.Logger {
	var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl <= zapcore.InfoLevel
	})

	var errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= zapcore.WarnLevel
	})

	var zapCore []zapcore.Core

	if logConfig.Log2File {
		c := zapcore.NewCore(
			getEncoder(),
			zapcore.AddSync(getLogWriter(name, txtLogNormal)),
			infoLevel)
		zapCore = append(zapCore, c)

		if name != LogNameAccess {
			c := zapcore.NewCore(
				getEncoder(),
				zapcore.AddSync(getLogWriter(name, txtLogWarnFatal)),
				errorLevel)
			zapCore = append(zapCore, c)
		}
	}

	// core
	core := zapcore.NewTee(zapCore...)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields()
	// 构造日志
	logger := zap.New(core, filed, caller, development)
	return logger
}


func getLogLevel(lv string) (level zapcore.Level) {
	str := strings.ToUpper(lv)
	switch str {
	case "DEBUG":
		level = zap.DebugLevel
	case "INFO":
		level = zap.InfoLevel
	case "WARN":
		level = zap.WarnLevel
	case "ERROR":
		level = zap.ErrorLevel
	case "FATAL":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

func getEncoder() zapcore.Encoder {
	// 公用编码器
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//LineEnding:    "tp=&tc=xxx.logger\n",
		EncodeCaller: zapcore.FullCallerEncoder, // 全路径编码器
		//EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter(name, logType string) (wr io.Writer) {
	// stdOut
	if logType == txtLogStdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

	// 写日志到文件filename中
	filename := genFilename(name, logType)
	return  NewTimeFileLogWriter(filename)
}

func genFilename(appName, logType string) string {
	var tailFixed string
	switch logType {
	case txtLogNormal:
		tailFixed = ".log"
	case txtLogWarnFatal:
		tailFixed = ".log.wf"
	default:
		tailFixed = ".log"
	}

	return appName+tailFixed
}


func CloseLogger() {
	if ServerLogger != nil {
		_ = ServerLogger.Sync()
	}

	if AccessLogger != nil {
		_ = AccessLogger.Sync()
	}
}

// 避免用户改动过大，以下为封装的之前的Entry打印field的方法
type Fields map[string]interface{}
type entry struct {
	s *zap.SugaredLogger
}

func NewEntry(s *zap.SugaredLogger) *entry {
	x := s.Desugar().WithOptions(zap.AddCallerSkip(+1)).Sugar()
	return &entry{s: x}
}

// 注意这种使用方式固定头的顺序会变
func (e entry) WithFields(f Fields) *zap.SugaredLogger {
	var fields []interface{}
	for k, v := range f {
		fields = append(fields, k, v)
	}

	return e.s.With(fields...)
}

func DebugLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Debug(msg, fields...)
}
func InfoLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Info(msg, fields...)
}

func WarnLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Warn(msg, fields...)
}

func ErrorLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Error(msg, fields...)
}

func PanicLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Panic(msg, fields...)
}

func FatalLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Fatal(msg, fields...)
}

func zapLogger(ctx *gin.Context) *zap.Logger {
	m := GetLogger()
	if ctx == nil {
		return m
	}
	return m.With(
		zap.String("logId", GetLogID(ctx)),
		zap.String("requestId", GetRequestID(ctx)),
	)
}
