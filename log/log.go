package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	if err := InitLogger(); err != nil {
		panic("init logger failed")
	}
}

func GetLogger() *zap.Logger {
	if _logger == nil {
		return nil
	}

	return _logger.Desugar()
}

type Level string

const (
	InfoLevel  Level = "info"
	DebugLevel Level = "debug"
)

func (l Level) toZapLevel() zapcore.Level {
	return map[Level]zapcore.Level{
		InfoLevel:  zap.InfoLevel,
		DebugLevel: zap.DebugLevel,
	}[l]
}

type Logger struct {
	*zap.Logger
}

var _logger *zap.SugaredLogger

type Option interface {
	apply(c *Config)
}

type optionFunc func(c *Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

func withDefaultOption() Option {
	return optionFunc(func(c *Config) {
		c.level = InfoLevel
		c.useConsole = true
		c.maxSize = 10
		c.maxAge = 30
		c.maxBackups = 10
	})
}

func WithLevel(level Level) Option {
	return optionFunc(func(c *Config) {
		c.level = level
	})
}

func WithUseConsole(useConsole bool) Option {
	return optionFunc(func(c *Config) {
		c.useConsole = useConsole
	})
}

func WithPath(path string) Option {
	return optionFunc(func(c *Config) {
		c.path = path
	})
}

func WithMaxSize(maxSize int) Option {
	return optionFunc(func(c *Config) {
		c.maxSize = maxSize
	})
}

func WithMaxAge(maxAge int) Option {
	return optionFunc(func(c *Config) {
		c.maxAge = maxAge
	})
}

func WithMaxBackups(maxBackups int) Option {
	return optionFunc(func(c *Config) {
		c.maxBackups = maxBackups
	})
}

type Config struct {
	level      Level
	useConsole bool
	path       string
	maxSize    int // MB
	maxAge     int // Day
	maxBackups int
}

func InitLogger(opts ...Option) error {
	c := &Config{}
	withDefaultOption().apply(c)
	for _, opt := range opts {
		opt.apply(c)
	}

	zapOpts := []zap.Option{
		zap.WithCaller(true),
		zap.AddCallerSkip(1),
	}

	var logger *zap.Logger
	var err error
	if c.useConsole {
		zc := zap.NewDevelopmentConfig()
		zc.Level = zap.NewAtomicLevelAt(c.level.toZapLevel())

		logger, err = zc.Build(zapOpts...)
		if err != nil {
			return err
		}
	} else {
		lumberjackLogger := lumberjack.Logger{
			Filename:   c.path,
			MaxSize:    c.maxSize,
			MaxAge:     c.maxAge,
			MaxBackups: c.maxBackups,
			Compress:   true,
		}
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder()), zapcore.AddSync(&lumberjackLogger), c.level.toZapLevel())
		//core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder()), zapcore.AddSync(&lumberjackLogger), c.level.toZapLevel())
		logger = zap.New(core, zapOpts...)
	}

	_logger = logger.Sugar()

	return nil
}

func encoder() zapcore.EncoderConfig {
	return zap.NewDevelopmentEncoderConfig()
}

func NewLogger(cfg Config) (*Logger, error) {
	return nil, nil
}

func SetLogger(logger *Logger) {
	if logger != nil {
		_logger = logger.Sugar()
	}
}

func Info(args ...interface{}) {
	_logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	_logger.Infof(format, args)
}

func Debug(args ...interface{}) {
	_logger.Debug(args)
}

func Debugf(template string, args ...interface{}) {
	_logger.Debugf(template, args)
}

func Error(args ...interface{}) {
	_logger.Error(args)
}

func Errorf(template string, args ...interface{}) {
	_logger.Errorf(template, args)
}

func Warn(args ...interface{}) {
	_logger.Warn(args)
}

func Warnf(template string, args ...interface{}) {
	_logger.Warnf(template, args)
}
