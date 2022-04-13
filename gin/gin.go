package gin

import (
	"fmt"
	cus_common "github.com/cheerUpPing/cus-common"
	"github.com/cheerUpPing/cus-common/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
	"time"
)

func traceIdMiddle(ctx *gin.Context) {
	uid := uuid.New()
	ctx.Set(cus_common.TRACE_ID, uid)
	ctx.Next()
}

func requestTime(ctx *gin.Context) {
	begin := time.Now()
	ctx.Next()
	delta := time.Since(begin)
	log.LogInfo(ctx.GetString(cus_common.TRACE_ID), fmt.Sprintf("cost: %s", delta.String()))
}

func handleException(ctx *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			log.LogError(ctx.GetString(cus_common.TRACE_ID), err.(error))
			ctx.JSON(http.StatusOK, cus_common.Error())
		}
	}()
	ctx.Next()
}

func AddTraceMiddle(engine *gin.Engine) {
	engine.Use(traceIdMiddle, requestTime, handleException)
}
