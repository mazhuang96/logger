# logger

Simple to use zap log with file-rotatelogs

## Usage

### install

```shell
go get -u github.com/mazhuang96/logger
```

### example

```go
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
```

output:

```log
test 2021/06/21 - 16:30:04.963  DEBUG   hello world!
test 2021/06/21 - 16:30:04.963  INFO    today is holiday!
test 2021/06/21 - 16:30:04.963  INFO    /home/mazhuang/zaplogger/example/main.go:27     I came to the Old Summer Palace!
test 2021/06/21 - 16:30:04.964  ERROR   /home/mazhuang/zaplogger/example/main.go:29     BUt i didn't bring cash!        {"error": "Fuck"}
main.main
        /home/mazhuang/zaplogger/example/main.go:29
runtime.main
        /usr/local/go/src/runtime/proc.go:204
test 2021/06/21 - 16:30:04.964  INFO    /home/mazhuang/zaplogger/example/main.go:34     I received a message on my phone!
main.main
        /home/mazhuang/zaplogger/example/main.go:34
runtime.main
        /usr/local/go/src/runtime/proc.go:204
test 2021/06/21 - 16:30:04.964  WARN    /home/mazhuang/zaplogger/example/main.go:37     It's going to rain in the afternoon !   {"start from": "14:00", "probability ": "60%"}
{"level":"FATAL","time":"test 2021/06/21 - 16:30:04.964","caller":"/home/mazhuang/zaplogger/example/main.go:40","message":"But I forgot to bring an umbrella!","error":"not found"}
exit status 1
```
