package logger

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func integerLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(int(l))
}

// InitLogger sets up a standardized logger according to standards setup by payments.
// Based on the `appenv` argument, the debugging level will be set to Info ("production") or Debug (anything else).
func InitLogger(appenv string) (*zap.Logger, error) {
	logConf := zap.NewProductionConfig()

	switch appenv {
	case "production", "live":
		logConf.Level.SetLevel(zap.InfoLevel) // be explicit.
	default:
		logConf.Level.SetLevel(zap.DebugLevel)
	}

	logConf.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",
		LineEnding: zapcore.DefaultLineEnding,

		TimeKey:        "timestamp",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,

		LevelKey:    "level",
		EncodeLevel: integerLevelEncoder,

		StacktraceKey: "stacktrace",

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	viewHook, err := newViewHook()
	if err != nil {
		return nil, err
	}

	logger, err := logConf.Build(zap.Hooks(viewHook))
	if err != nil {
		return nil, err
	}

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return logger, nil
}

var (
	logLevelTag     = tag.MustNewKey("log_level")
	messagesCounter = stats.Int64(
		"kit/zap/message_count",
		"Amount of messages logged for each log level",
		stats.UnitDimensionless,
	)
	messagesView = view.View{
		Name:        "kit/zap/message_count",
		Description: "Amount of messages logged for each log level",
		TagKeys: []tag.Key{
			logLevelTag,
		},
		Measure:     messagesCounter,
		Aggregation: view.Sum(),
	}

	// initialisingLogLevels are the log levels we will initialise the metric for on Hook creation
	initialisingLogLevels = []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}
)

func newViewHook() (func(zapcore.Entry) error, error) {
	err := view.Register(&messagesView) // registering a view multiple times is fine, it will ignore subsequent calls
	if err != nil {
		return nil, err
	}

	// initialise view with all log levels
	for _, logLevel := range initialisingLogLevels {

		err = stats.RecordWithTags(
			context.Background(),
			[]tag.Mutator{
				tag.Upsert(logLevelTag, logLevel.String()),
			},
			messagesCounter.M(0),
		)
		if err != nil {
			return nil, err
		}
	}

	hook := func(entry zapcore.Entry) error {

		err = stats.RecordWithTags(
			context.Background(),
			[]tag.Mutator{
				tag.Upsert(logLevelTag, entry.Level.String()),
			},
			messagesCounter.M(1),
		)

		return err
	}

	return hook, nil
}
