/**************************************
 * @Author: mazhuang
 * @Date: 2021-06-21 10:35:38
 * @LastEditTime: 2021-06-22 11:22:16
 * @Description:
 **************************************/

package main

import (
	"errors"

	"github.com/mazhuang96/logger"

	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewDefault()
	if err != nil {
		panic(err)
	}
	log.Debug("hello world!")
	log.Info("today is holiday!")
	// Set the stack trace level to debug
	log.SetStacktraceLevel("debug").Info("I came to the Old Summer Palace!")
	// Show call line number
	log.Showline().Error("BUt i didn't bring cash!", zap.Error(errors.New("Fuck")))
	// Turn off print colors
	log.NoColor().Info("I received a message on my phone!")
	// Print as json format
	log.SetJSONStyle().Warn("It's going to rain in the afternoon!", zap.String("start from", "14:00"), zap.String("probability ", "60%"))
	// Turn off print stack
	log.CloseStacktrace().Fatal("But I forgot to bring an umbrella!", zap.Error(errors.New("not found")))
}
