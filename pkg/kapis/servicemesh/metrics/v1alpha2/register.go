/*
Copyright 2020 KubeSphere Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	"kubesphere.io/kubesphere/pkg/apiserver/runtime"
	"kubesphere.io/kubesphere/pkg/simple/client/cache"
	"kubesphere.io/kubesphere/pkg/simple/client/servicemesh"
)

const groupName = "servicemesh.kubesphere.io"

var GroupVersion = schema.GroupVersion{Group: groupName, Version: "v1alpha2"}

func AddToContainer(o *servicemesh.Options, c *restful.Container, client kubernetes.Interface, cache cache.Interface) error {

	tags := []string{"ServiceMesh"}

	webservice := runtime.NewWebService(GroupVersion)

	h := NewHandler(o, client, cache)

	// Get service metrics
	// GET /namespaces/{namespace}/services/{service}/metrics
	webservice.Route(webservice.GET("/namespaces/{namespace}/services/{service}/metrics").
		To(h.GetServiceMetrics).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get service metrics from a specific namespace").
		Param(webservice.PathParameter("namespace", "name of the namespace")).
		Param(webservice.PathParameter("service", "name of the service")).
		Param(webservice.QueryParameter("filters[]", "type of metrics type, fetch all metrics when empty, e.g. request_count, request_duration, request_error_count").DefaultValue("[]")).
		Param(webservice.QueryParameter("queryTime", "from which UNIX time to extract metrics")).
		Param(webservice.QueryParameter("duration", "duration of the query period, in seconds").DefaultValue("1800")).
		Param(webservice.QueryParameter("step", "step between graph data points, in seconds.").DefaultValue("15")).
		Param(webservice.QueryParameter("rateInterval", "metrics rate intervals, e.g. 20s").DefaultValue("1m")).
		Param(webservice.QueryParameter("direction", "traffic direction: 'inbound' or 'outbound'").DefaultValue("outbound")).
		Param(webservice.QueryParameter("quantiles[]", "list of quantiles to fetch, fetch no quantiles when empty. eg. 0.5, 0.9, 0.99").DefaultValue("[]")).
		Param(webservice.QueryParameter("byLabels[]", "list of labels to use for grouping metrics(via Prometheus 'by' clause), e.g. source_workload, destination_service_name").DefaultValue("[]")).
		Param(webservice.QueryParameter("requestProtocol", "request protocol for the telemetry, e.g. http/tcp/grpc").DefaultValue("all protocols")).
		Param(webservice.QueryParameter("reporter", "istio telemetry reporter, 'source' or 'destination'").DefaultValue("source")).
		Produces(restful.MIME_JSON))

	// Get app metrics
	// Get /namespaces/{namespace}/apps/{app}/metrics
	webservice.Route(webservice.GET("/namespaces/{namespace}/apps/{app}/metrics").
		To(h.GetAppMetrics).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get app metrics from a specific namespace").
		Param(webservice.PathParameter("namespace", "name of the namespace")).
		Param(webservice.PathParameter("app", "name of the app")).
		Param(webservice.QueryParameter("filters[]", "type of metrics type, fetch all metrics when empty, e.g. request_count, request_duration, request_error_count").DefaultValue("[]")).
		Param(webservice.QueryParameter("queryTime", "from which UNIX time to extract metrics")).
		Param(webservice.QueryParameter("duration", "duration of the query period, in seconds").DefaultValue("1800")).
		Param(webservice.QueryParameter("step", "step between graph data points, in seconds.").DefaultValue("15")).
		Param(webservice.QueryParameter("rateInterval", "metrics rate intervals, e.g. 20s").DefaultValue("1m")).
		Param(webservice.QueryParameter("direction", "traffic direction: 'inbound' or 'outbound'").DefaultValue("outbound")).
		Param(webservice.QueryParameter("quantiles[]", "list of quantiles to fetch, fetch no quantiles when empty. eg. 0.5, 0.9, 0.99").DefaultValue("[]")).
		Param(webservice.QueryParameter("byLabels[]", "list of labels to use for grouping metrics(via Prometheus 'by' clause), e.g. source_workload, destination_service_name").DefaultValue("[]")).
		Param(webservice.QueryParameter("requestProtocol", "request protocol for the telemetry, e.g. http/tcp/grpc").DefaultValue("all protocols")).
		Param(webservice.QueryParameter("reporter", "istio telemetry reporter, 'source' or 'destination'").DefaultValue("source")).
		Produces(restful.MIME_JSON))

	// Get workload metrics
	// Get /namespaces/{namespace}/workloads/{workload}/metrics
	webservice.Route(webservice.GET("/namespaces/{namespace}/workloads/{workload}/metrics").
		To(h.GetWorkloadMetrics).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get workload metrics from a specific namespace").
		Param(webservice.PathParameter("namespace", "name of the namespace").Required(true)).
		Param(webservice.PathParameter("workload", "name of the workload").Required(true)).
		Param(webservice.QueryParameter("filters[]", "type of metrics type, fetch all metrics when empty, e.g. request_count, request_duration, request_error_count").DefaultValue("[]")).
		Param(webservice.QueryParameter("queryTime", "from which UNIX time to extract metrics")).
		Param(webservice.QueryParameter("duration", "duration of the query period, in seconds").DefaultValue("1800")).
		Param(webservice.QueryParameter("step", "step between graph data points, in seconds.").DefaultValue("15")).
		Param(webservice.QueryParameter("rateInterval", "metrics rate intervals, e.g. 20s").DefaultValue("1m")).
		Param(webservice.QueryParameter("direction", "traffic direction: 'inbound' or 'outbound'").DefaultValue("outbound")).
		Param(webservice.QueryParameter("quantiles[]", "list of quantiles to fetch, fetch no quantiles when empty. eg. 0.5, 0.9, 0.99").DefaultValue("[]")).
		Param(webservice.QueryParameter("byLabels[]", "list of labels to use for grouping metrics(via Prometheus 'by' clause), e.g. source_workload, destination_service_name").DefaultValue("[]")).
		Param(webservice.QueryParameter("requestProtocol", "request protocol for the telemetry, e.g. http/tcp/grpc").DefaultValue("all protocols")).
		Param(webservice.QueryParameter("reporter", "istio telemetry reporter, 'source' or 'destination'").DefaultValue("source")).
		Produces(restful.MIME_JSON))

	// Get namespace metrics
	// Get /namespaces/{namespace}/metrics
	webservice.Route(webservice.GET("/namespaces/{namespace}/metrics").
		To(h.GetNamespaceMetrics).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get metrics from a specific namespace").
		Param(webservice.PathParameter("namespace", "name of the namespace").Required(true)).
		Param(webservice.QueryParameter("filters[]", "type of metrics type, fetch all metrics when empty, e.g. request_count, request_duration, request_error_count").DefaultValue("[]")).
		Param(webservice.QueryParameter("queryTime", "from which UNIX time to extract metrics")).
		Param(webservice.QueryParameter("duration", "duration of the query period, in seconds").DefaultValue("1800")).
		Param(webservice.QueryParameter("step", "step between graph data points, in seconds.").DefaultValue("15")).
		Param(webservice.QueryParameter("rateInterval", "metrics rate intervals, e.g. 20s").DefaultValue("1m")).
		Param(webservice.QueryParameter("direction", "traffic direction: 'inbound' or 'outbound'").DefaultValue("outbound")).
		Param(webservice.QueryParameter("quantiles[]", "list of quantiles to fetch, fetch no quantiles when empty. eg. 0.5, 0.9, 0.99").DefaultValue("[]")).
		Param(webservice.QueryParameter("byLabels[]", "list of labels to use for grouping metrics(via Prometheus 'by' clause), e.g. source_workload, destination_service_name").DefaultValue("[]")).
		Param(webservice.QueryParameter("requestProtocol", "request protocol for the telemetry, e.g. http/tcp/grpc").DefaultValue("all protocols")).
		Param(webservice.QueryParameter("reporter", "istio telemetry reporter, 'source' or 'destination'").DefaultValue("source")).
		Produces(restful.MIME_JSON))

	// Get namespace graph
	// Get /namespaces/{namespace}/graph
	webservice.Route(webservice.GET("/namespaces/{namespace}/graph").
		To(h.GetNamespaceGraph).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get service graph for a specific namespace").
		Param(webservice.PathParameter("namespace", "name of a namespace").Required(true)).
		Param(webservice.QueryParameter("duration", "duration of the query period, in seconds").DefaultValue("10m")).
		Param(webservice.QueryParameter("graphType", "type of the generated service graph. Available graph types: [app, service, versionedApp, workload].").DefaultValue("workload")).
		Param(webservice.QueryParameter("groupBy", "app box grouping characteristic. Available groupings: [app, none, version].").DefaultValue("none")).
		Param(webservice.QueryParameter("queryTime", "from which time point in UNIX timestamp, default now")).
		Param(webservice.QueryParameter("injectServiceNodes", "flag for injecting the requested service node between source and destination nodes.").DefaultValue("false")).
		Returns(http.StatusBadRequest, "bad request", BadRequestError{}).
		Returns(http.StatusNotFound, "not found", NotFoundError{}).
		Produces(restful.MIME_JSON))

	// Get namespace health
	webservice.Route(webservice.GET("/namespaces/{namespace}/health").
		To(h.GetNamespaceHealth).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get app/service/workload health of a namespace").
		Param(webservice.PathParameter("namespace", "name of a namespace").Required(true)).
		Param(webservice.QueryParameter("rateInterval", "the rate interval used for fetching error rate").DefaultValue("10m").Required(true)).
		Param(webservice.QueryParameter("queryTime", "the time to use for query")).
		Returns(http.StatusBadRequest, "bad request", BadRequestError{}).
		Returns(http.StatusNotFound, "not found", NotFoundError{}).
		Produces(restful.MIME_JSON))

	// Get workloads health
	webservice.Route(webservice.GET("/namespaces/{namespace}/workloads/{workload}/health").
		To(h.GetWorkloadHealth).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get workload health").
		Param(webservice.PathParameter("namespace", "name of a namespace").Required(true)).
		Param(webservice.PathParameter("workload", "workload name").Required(true)).
		Param(webservice.QueryParameter("rateInterval", "the rate interval used for fetching error rate").DefaultValue("10m").Required(true)).
		Param(webservice.QueryParameter("queryTime", "the time to use for query")).
		Produces(restful.MIME_JSON))

	// Get app health
	webservice.Route(webservice.GET("/namespaces/{namespace}/apps/{app}/health").
		To(h.GetAppHealth).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get app health").
		Param(webservice.PathParameter("namespace", "name of a namespace").Required(true)).
		Param(webservice.PathParameter("app", "app name").Required(true)).
		Param(webservice.QueryParameter("rateInterval", "the rate interval used for fetching error rate").DefaultValue("10m").Required(true)).
		Param(webservice.QueryParameter("queryTime", "the time to use for query")).
		Produces(restful.MIME_JSON))

	// Get service health
	webservice.Route(webservice.GET("/namespaces/{namespace}/services/{service}/health").
		To(h.GetServiceHealth).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("Get service health").
		Param(webservice.PathParameter("namespace", "name of a namespace").Required(true)).
		Param(webservice.PathParameter("service", "service name").Required(true)).
		Param(webservice.QueryParameter("rateInterval", "the rate interval used for fetching error rate").DefaultValue("10m").Required(true)).
		Param(webservice.QueryParameter("queryTime", "the time to use for query")).
		Produces(restful.MIME_JSON))

	// Get service tracing
	webservice.Route(webservice.GET("/namespaces/{namespace}/services/{service}/traces").
		To(h.GetServiceTracing).
		Doc("Get tracing of a service, should have servicemesh enabled first").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(webservice.PathParameter("namespace", "namespace of service").Required(true)).
		Param(webservice.PathParameter("service", "name of service queried").Required(true)).
		Param(webservice.QueryParameter("start", "start of time range want to query, in unix timestamp")).
		Param(webservice.QueryParameter("end", "end of time range want to query, in unix timestamp")).
		Param(webservice.QueryParameter("limit", "maximum tracing entries returned at one query, default 10").DefaultValue("10")).
		Param(webservice.QueryParameter("loopback", "loopback of duration want to query, e.g. 30m/1h/2d")).
		Param(webservice.QueryParameter("maxDuration", "maximum duration of a request")).
		Param(webservice.QueryParameter("minDuration", "minimum duration of a request")).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))

	c.Add(webservice)

	return nil
}
