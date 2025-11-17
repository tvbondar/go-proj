package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(debug bool) error {
	var z *zap.Logger
	var err error
	if debug {
		z, err = zap.NewDevelopment()
	} else {
		z, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	sugar := z.Sugar()
	log = sugar
	return nil
}

func L() *zap.SugaredLogger {
	return log
}
