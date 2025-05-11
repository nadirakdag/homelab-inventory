package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init() {
	raw, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = raw.Sugar()
}
