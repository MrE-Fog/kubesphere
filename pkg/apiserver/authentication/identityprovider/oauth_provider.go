/*
Copyright 2020 The KubeSphere Authors.

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

package identityprovider

import (
	"net/http"

	"kubesphere.io/kubesphere/pkg/server/options"
)

type OAuthProvider interface {
	// IdentityExchangeCallback handle oauth callback, exchange identity from remote server
	IdentityExchangeCallback(req *http.Request) (Identity, error)
}

type OAuthProviderFactory interface {
	// Type unique type of the provider
	Type() string
	// Create Apply the dynamic options
	Create(options options.DynamicOptions) (OAuthProvider, error)
}
