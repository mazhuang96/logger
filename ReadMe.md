# logger

Simple to use [zap](https://github.com/uber-go/zap) log with [file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)

## Usage

### install

```shell
go get -u github.com/mazhuang96/logger
```

## example

```go
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
```

The default configuration is to print with color level and log level above "info". When the level is higher than "error", the stack will be traced and the call line will not be displayed.

log levels from low to high: "debug", "info", "warn", "error", "dpanic", "panic", "fatal",

### config

When the default configuration does not meet the requirements, please use a custom configuration `logger.New(Config)`

```go
type Config struct {
    Level      string // log print level, default "info"
    // Please make sure the folder exists or has permission to read and write 
    Dir        string // log output folder
    Prefix     string // prefix on each line to identify the logger
    TimeFormat string // the time format of each line in the log
    MaxAge     int    // max age (days) of each log file, default 7 days
    Color      bool   // the color of log level
    ShowLine   bool   // show log call line number
    Stacktrace string // stack trace log level
    Encoder    string // log encoding format, divided into "json" and "console", default "console"
}
```

### output

The following is the console output of the default configuration

```log
$ go run example/main.go
2021/06/22 - 10:14:13.067       INFO    today is holiday!
2021/06/22 - 10:14:13.068       INFO    I came to the Old Summer Palace!
main.main
        /home/mazhuang/zaplogger/example/main.go:26
runtime.main
        /usr/local/go/src/runtime/proc.go:204
2021/06/22 - 10:14:13.068       ERROR   /home/mazhuang/zaplogger/example/main.go:28     BUt i didn't bring cash!        {"error": "Fuck"}
main.main
        /home/mazhuang/zaplogger/example/main.go:28
runtime.main
        /usr/local/go/src/runtime/proc.go:204
2021/06/22 - 10:14:13.068       INFO    I received a message on my phone!
{"level":"\u001b[33mWARN\u001b[0m","time":"2021/06/22 - 10:14:13.068","message":"It's going to rain in the afternoon!","start from":"14:00","probability ":"60%"}
2021/06/22 - 10:14:13.068       FATAL   But I forgot to bring an umbrella!      {"error": "not found"}
exit status 1
```
