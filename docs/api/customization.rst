==============
Customizations
==============

Add custom private routes
-------------------------

Private routes have several middlewares that are applied for each request. If you want to create custom routes/handlers, you can do this in a very simple way:

.. code:: go

  //Custom handler
  type HiHandler struct {
    Handler
  }
  func (handler *HiHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
    return OK("Hi from custom route!")
  }

  func main() {
    var config = &Config{
      FilePath: "config.yaml",
      Port: ":8000",
    }
    var api = NewApi(config)
    api.AddPrivateRoute("GET", "/hi", &HiHandler{}, "Index")
  }


Use different log library
-------------------------

By default, if no library is provided, all the logs will be printed out on the stdout. But, it is pretty easy to use any library. It's just required to implement the following interface:

.. code:: go

  type Log interface {
    // Debug logs information that is diagnostically helpful to developers.
    Debug(format string, args ...interface{})

    // Info logs useful information to log.
    Info(format string, args ...interface{})

    // Error logs any error which is fatal to the operation.
    Error(format string, args ...interface{})

    // Warn logs message with severity "warn". Anything that can potentially cause error.
    Warn(format string, args ...interface{})

    // Disable will prevent the application to log anything.
    Disable()

    // SetLevel sets the error reporting level.
    SetLevel(level int32)
  }

Follow bellow an example using `Go-Logging <https://github.com/ccding/go-logging>`_:

.. code:: go

  package main

  import (
    "github.com/backstage/maestro/log"
    "github.com/ccding/go-logging/logging"
  )

  type CustomLogger struct {
    disabled bool
    logger   *logging.Logger
  }

  func NewCustomLogger() *CustomLogger {
    l := new(CustomLogger)
    l.logger, _ = logging.SimpleLogger("main")
    return l
  }

  func (l *CustomLogger) Debug(format string, args ...interface{}) {
    if !l.disabled {
      l.logger.Debugf(format, args...)
    }
  }

  func (l *CustomLogger) Info(format string, args ...interface{}) {
    if !l.disabled {
      l.logger.Infof(format, args...)
    }
  }

  func (l *CustomLogger) Warn(format string, args ...interface{}) {
    if !l.disabled {
      l.logger.Warnf(format, args...)
    }
  }

  func (l *CustomLogger) Error(format string, args ...interface{}) {
    if !l.disabled {
      l.logger.Errorf(format, args...)
    }
  }

  func (l *CustomLogger) Disable() {
    l.disabled = true
  }

  func (l *CustomLogger) SetLevel(level int32) {
    levelStr := log.GetLevelFlagName(level)
    l.logger.SetLevel(logging.GetLevelValue(levelStr))
  }


Then, you just need to configure the api to use that:

.. code:: go

  var api = NewApi(config)
  logger := NewCustomLogger()
  logger.SetLevel(log.DEBUG)
  api.Logger(logger)
