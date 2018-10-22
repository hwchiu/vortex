package server

import (
	"github.com/hwchiu/vortex/src/entity"
	response "github.com/hwchiu/vortex/src/net/http"
	"github.com/hwchiu/vortex/src/net/http/query"
	pc "github.com/hwchiu/vortex/src/prometheuscontroller"
	"github.com/hwchiu/vortex/src/web"
)

func getContainerMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	podId := req.PathParameter("pod")
	containerId := req.PathParameter("container")

	rs := pc.RangeSetting{}
	query := query.New(req.Request.URL.Query())

	var err error
	rs.Interval, err = query.TimeDuration("interval", 2)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Resolution, err = query.TimeDuration("resolution", 10)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Rate, err = query.TimeDuration("rate", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	container, err := pc.GetContainer(sp, podId, containerId, rs)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(container)
}

func getPodMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	id := req.PathParameter("pod")

	rs := pc.RangeSetting{}
	query := query.New(req.Request.URL.Query())

	var err error
	rs.Interval, err = query.TimeDuration("interval", 2)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Resolution, err = query.TimeDuration("resolution", 10)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Rate, err = query.TimeDuration("rate", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	pod, err := pc.GetPod(sp, id, rs)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(pod)
}

func getServiceMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	id := req.PathParameter("service")

	service, err := pc.GetService(sp, id)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(service)
}

func getControllerMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	id := req.PathParameter("controller")

	controller, err := pc.GetController(sp, id)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(controller)
}

func getNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	id := req.PathParameter("node")

	rs := pc.RangeSetting{}
	query := query.New(req.Request.URL.Query())

	var err error
	rs.Interval, err = query.TimeDuration("interval", 2)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Resolution, err = query.TimeDuration("resolution", 10)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Rate, err = query.TimeDuration("rate", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	node, err := pc.GetNode(sp, id, rs)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(node)
}

func listPodMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	queryLabels := map[string]string{}

	if node, ok := query.Str("node"); ok {
		queryLabels["node"] = node
	}

	if namespace, ok := query.Str("namespace"); ok {
		queryLabels["namespace"] = namespace
	}

	if controller, ok := query.Str("controller"); ok {
		queryLabels["created_by_kind"] = "ReplicaSet"
		queryLabels["created_by_name"] = controller + ".*"
	}

	rs := pc.RangeSetting{}

	var err error
	rs.Interval, err = query.TimeDuration("interval", 2)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Resolution, err = query.TimeDuration("resolution", 10)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Rate, err = query.TimeDuration("rate", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	podNameList, err := pc.ListPodName(sp, queryLabels)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	podList := map[string]entity.PodMetrics{}
	for _, podName := range podNameList {
		pod, err := pc.GetPod(sp, podName, rs)
		if err != nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		podList[podName] = pod
	}

	resp.WriteEntity(podList)
}

func listServiceMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	queryLabels := map[string]string{}

	if namespace, ok := query.Str("namespace"); ok {
		queryLabels["namespace"] = namespace
	}

	serviceNameList, err := pc.ListServiceName(sp, queryLabels)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	serviceList := map[string]entity.ServiceMetrics{}
	for _, serviceName := range serviceNameList {
		service, err := pc.GetService(sp, serviceName)
		if err != nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		serviceList[serviceName] = service
	}

	resp.WriteEntity(serviceList)
}

func listControllerMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	query := query.New(req.Request.URL.Query())
	queryLabels := map[string]string{}

	if namespace, ok := query.Str("namespace"); ok {
		queryLabels["namespace"] = namespace
	}

	controllerNameList, err := pc.ListControllerName(sp, queryLabels)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	controllerList := map[string]entity.ControllerMetrics{}
	for _, controllerName := range controllerNameList {
		controller, err := pc.GetController(sp, controllerName)
		if err != nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		controllerList[controllerName] = controller
	}

	resp.WriteEntity(controllerList)
}

func listNodeMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	nodeNameList, err := pc.ListNodeName(sp, map[string]string{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	rs := pc.RangeSetting{}
	query := query.New(req.Request.URL.Query())

	rs.Interval, err = query.TimeDuration("interval", 2)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Resolution, err = query.TimeDuration("resolution", 10)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	rs.Rate, err = query.TimeDuration("rate", 1)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	nodeList := map[string]entity.NodeMetrics{}
	for _, nodeName := range nodeNameList {
		node, err := pc.GetNode(sp, nodeName, rs)
		if err != nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
		nodeList[nodeName] = node
	}

	resp.WriteEntity(nodeList)
}

func listNodeNicsMetricsHandler(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	id := req.PathParameter("node")

	nicList, err := pc.ListNodeNICs(sp, id)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(nicList)
}
