package c_gin

import (
	"fmt"
	c_common "github.com/cheerUpPing/cus-common"
	c_log "github.com/cheerUpPing/cus-common/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
	"time"
)

func traceIdMiddle(ctx *gin.Context) {
	uid := uuid.New()
	ctx.Set(c_common.TRACE_ID, uid)
	ctx.Header(c_common.TRACE_ID, uid)
	ctx.Next()
}

func requestTime(ctx *gin.Context) {
	begin := time.Now()
	ctx.Next()
	delta := time.Since(begin)
	c_log.LogInfo(ctx.GetString(c_common.TRACE_ID), fmt.Sprintf("cost: %s", delta.String()))
}

func handleException(ctx *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			c_log.LogError(ctx.GetString(c_common.TRACE_ID), err.(error))
			ctx.JSON(http.StatusOK, c_common.Error())
		}
	}()
	ctx.Next()
}

func AddTraceMiddle(engine *gin.Engine) {
	engine.Use(traceIdMiddle, requestTime, handleException)
}
