/**************************************
 * @Author: mazhuang
 * @Date: 2021-06-21 10:35:38
 * @LastEditTime: 2021-06-21 16:25:27
 * @Description:
 **************************************/

package main

import (
	"errors"

	"github.com/mazhuang96/logger"

	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger("debug", "logs", "test")
	if err != nil {
		panic(err)
	}
	log.Debug("hello world!")
	log.Info("today is holiday!")
	// Show call line number
	log.SetShowline()
	log.Info("I came to the Old Summer Palace!")
	// Print stack above error level
	log.Error("BUt i didn't bring cash!", zap.Error(errors.New("Fuck")))
	// Turn off print colors
	log.CloseColor()
	// Set the print stack level to debug
	log.SetStacktraceLevel("debug")
	log.Info("I received a message on my phone!")
	// Turn off print stack
	log.CloseStacktrace()
	log.Warn("It's going to rain in the afternoon !", zap.String("start from", "14:00"), zap.String("probability ", "60%"))
	// Print as json format
	log.SetJSONStyle()
	log.Fatal("But I forgot to bring an umbrella!", zap.Error(errors.New("not found")))
}
