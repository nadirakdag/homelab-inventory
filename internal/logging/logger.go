package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init(debug bool) {
	var base *zap.Logger
	var err error

	if debug {
		base, err = zap.NewDevelopment()
	} else {
		base, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	Logger = base.Sugar()
}
