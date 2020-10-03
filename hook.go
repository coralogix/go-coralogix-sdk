package coralogix

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "path/filepath"
    "runtime"
    "strings"
)

// Hook is presenting Coralogix Logger for Logrus library
type Hook struct {
    Writer LoggerManager
    Options HookOptions
}


type HookOptions struct {
    // Set to true to extract error message from error type
    ExtractError bool

    // Set to true to use the global ReportCaller value
    UseReportCaller bool
}

// NewCoralogixHook build Coralogix logger hook
func NewCoralogixHook(PrivateKey string, ApplicationName string, SubsystemName string, options ...HookOptions) *Hook {
    var op HookOptions;
    if options != nil && len(options) > 0{
        op = options[0]
    }

    NewHookInstance := &Hook{
        Writer: *NewLoggerManager(
            PrivateKey,
            ApplicationName,
            SubsystemName,
            true,
        ),
        Options: op,
    }

    go NewHookInstance.Writer.Run()

    return NewHookInstance
}

// Extract error message
func extractError(entry *logrus.Entry) {
    if entry != nil && entry.Data != nil{
        for key, val := range entry.Data {
            switch e := val.(type) {
            case error:
                entry.Data[key] = e.Error()
            }
        }
    }
}

// Fire send message to Coralogix
func (hook *Hook) Fire(entry *logrus.Entry) error {
    var (
        Level      uint
        Text       interface{}
        Category   string
        ClassName  string
        MethodName string
        ThreadId   string
    )

    switch entry.Level {
    case logrus.TraceLevel:
        Level = 1
    case logrus.DebugLevel:
        Level = 1
    case logrus.InfoLevel:
        Level = 3
    case logrus.WarnLevel:
        Level = 4
    case logrus.ErrorLevel:
        Level = 5
    case logrus.FatalLevel:
        Level = 6
    case logrus.PanicLevel:
        Level = 6
    }

    // If user is using the global log.SetReportCaller then lets get the data
    if hook.Options.UseReportCaller && entry.Caller != nil {
        Category, ClassName, MethodName = getCallerInformation(entry.Caller)
    }

    if Value, Exist := entry.Data["Category"]; Exist {
        Category = Value.(string)
        delete(entry.Data, "Category")
    } else {
        Category = ""
    }

    if Value, Exist := entry.Data["ClassName"]; Exist {
        ClassName = Value.(string)
        delete(entry.Data, "ClassName")
    } else {
        ClassName = ""
    }

    if Value, Exist := entry.Data["MethodName"]; Exist {
        MethodName = Value.(string)
        delete(entry.Data, "MethodName")
    } else {
        MethodName = ""
    }

    if Value, Exist := entry.Data["ThreadId"]; Exist {
        ThreadId = Value.(string)
        delete(entry.Data, "ThreadId")
    } else {
        ThreadId = ""
    }

    if hook.Options.ExtractError {
        extractError(entry)
    }

    if len(entry.Data) > 0 {
        Text = map[string]interface{}{
            "message": entry.Message,
            "fields":  entry.Data,
        }
    } else {
        Text = entry.Message
    }

    hook.Writer.LogsBuffer = append(
        hook.Writer.LogsBuffer,
        Log{
            float64(entry.Time.Unix()) * 1000.0,
            Level,
            MessageToString(Text),
            Category,
            ClassName,
            MethodName,
            ThreadId,
        },
    )

    if entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel {
        hook.Writer.Flush()
    }

    return nil
}

func getCallerInformation(frame *runtime.Frame) (string, string, string) {
    if frame == nil { return "", "", "" }
    file := filepath.Base(frame.File)
    sections := strings.Split(frame.Func.Name(), "%")
    len := len(sections)
    if len <= 1 {
        return file, "", frame.Func.Name()
    }
    class :=  fmt.Sprintf("%s",sections[len-2:len-1])
    function :=  fmt.Sprintf("%s",sections[len-1:len])

    return file, class, function
}

// Levels return levels which can be sent with this hook
func (hook *Hook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// Close is a defer function for buffer cleanup before exit
func (hook *Hook) Close() {
    hook.Writer.Stop()
}
