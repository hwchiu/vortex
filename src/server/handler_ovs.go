package server

import (
	"fmt"

	response "github.com/hwchiu/vortex/src/net/http"
	"github.com/hwchiu/vortex/src/net/http/query"
	"github.com/hwchiu/vortex/src/ovscontroller"
	"github.com/hwchiu/vortex/src/web"
)

func getOVSPortInfoHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	//Get the parameter
	query := query.New(req.Request.URL.Query())
	nodeName, exist := query.Str("nodeName")
	if !exist {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The nodeName must not be empty"))
		return
	}

	bridgeName, exist := query.Str("bridgeName")
	if !exist {
		response.BadRequest(req.Request, resp.ResponseWriter, fmt.Errorf("The bridgeName must not be empty"))
		return
	}

	portStats, err := ovscontroller.DumpPorts(sp, nodeName, bridgeName)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(portStats)
}
