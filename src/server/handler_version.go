package server

import (
	response "github.com/hwchiu/vortex/src/net/http"
	"github.com/hwchiu/vortex/src/version"
	"github.com/hwchiu/vortex/src/web"
)

func versionHandler(ctx *web.Context) {
	_, _, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	resp.WriteEntity(response.ActionResponse{
		Error:   false,
		Message: version.GetVersion(),
	})
}
