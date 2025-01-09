package toast

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "error"
)

type Toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func New(level string, message string) Toast {
	return Toast{level, message}
}

func Success(c echo.Context, message string) Toast {
	return New(SUCCESS, message).SetHXTriggerHeader(c)
}

func Info(message string) Toast {
	return New(INFO, message)
}

func Warning(message string) Toast {
	return New(WARNING, message)
}

func Danger(message string) Toast {
	return New(DANGER, message)
}

func (t Toast) Error() string {
	return fmt.Sprintf("%s: %s", t.Level, t.Message)
}

func (t Toast) jsonify() (string, error) {
	t.Message = t.Error()
	eventMap := map[string]Toast{}
	eventMap["makeToast"] = t
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (t Toast) SetHXTriggerHeader(c echo.Context) Toast {
	jsonData, _ := t.jsonify()
	c.Response().Header().Set("HX-Trigger", jsonData)
	return t
}

func ToastMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		handleToast(err, c)
		return err
	}
}

func handleToast(err error, c echo.Context) {
	if err == nil {
		return
	}

	he, _ := err.(*echo.HTTPError)
	te, ok := he.Unwrap().(Toast)

	if !ok {
		fmt.Println(err)
		te = Danger("there has been an unexpected error")
	}

	if te.Level != SUCCESS && te.Level != INFO {
		c.Response().Header().Set("HX-Reswap", "none")
	}

	_ = te.SetHXTriggerHeader(c)
}
