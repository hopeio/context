package context

import (
	"github.com/hopeio/utils/log"
	"go.uber.org/zap"
)

func (c *Context) ErrorLog(args ...any) {
	log.Errorw(log.Sprintln(args...), zap.String(log.FieldTraceId, c.traceID))
}

func (c *Context) RespErrorLog(err, originErr error, funcName string) error {
	// caller 用原始logger skip刚好
	log.GetCallerSkipLogger(1).Errorw(originErr.Error(), zap.String(log.FieldTraceId, c.traceID), zap.String(log.FieldPosition, funcName))
	return err
}
