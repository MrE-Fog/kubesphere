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

package v1alpha1

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sinformers "k8s.io/client-go/informers"

	"kubesphere.io/kubesphere/pkg/api"
	"kubesphere.io/kubesphere/pkg/apiserver/runtime"
	kubesphere "kubesphere.io/kubesphere/pkg/client/clientset/versioned"
	"kubesphere.io/kubesphere/pkg/client/informers/externalversions"
	"kubesphere.io/kubesphere/pkg/constants"
)

const (
	GroupName = "cluster.kubesphere.io"
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

func AddToContainer(container *restful.Container,
	ksclient kubesphere.Interface,
	k8sInformers k8sinformers.SharedInformerFactory,
	ksInformers externalversions.SharedInformerFactory,
	proxyService string,
	proxyAddress string,
	agentImage string) error {

	webservice := runtime.NewWebService(GroupVersion)
	h := newHandler(ksclient, k8sInformers, ksInformers, proxyService, proxyAddress, agentImage)

	// returns deployment yaml for cluster agent
	webservice.Route(webservice.GET("/clusters/{cluster}/agent/deployment").
		Doc("Return deployment yaml for cluster agent.").
		Param(webservice.PathParameter("cluster", "Name of the cluster.").Required(true)).
		To(h.generateAgentDeployment).
		Returns(http.StatusOK, api.StatusOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.MultiClusterTag}))

	webservice.Route(webservice.POST("/clusters/validation").
		Doc("").
		Param(webservice.BodyParameter("cluster", "cluster specification")).
		To(h.validateCluster).
		Returns(http.StatusOK, api.StatusOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.MultiClusterTag}))

	webservice.Route(webservice.PUT("/clusters/{cluster}/kubeconfig").
		Doc("Update cluster kubeconfig.").
		Param(webservice.PathParameter("cluster", "Name of the cluster.").Required(true)).
		To(h.updateKubeConfig).
		Returns(http.StatusOK, api.StatusOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.MultiClusterTag}))

	container.Add(webservice)

	return nil
}
