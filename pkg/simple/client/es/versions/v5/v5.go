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

package v5

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v5"
	"github.com/elastic/go-elasticsearch/v5/esapi"
	jsoniter "github.com/json-iterator/go"

	"kubesphere.io/kubesphere/pkg/simple/client/es/versions"
)

type Elastic struct {
	client *elasticsearch.Client
	index  string
}

func New(address string, basicAuth bool, username, password, index string) (*Elastic, error) {
	var client *elasticsearch.Client
	var err error

	if !basicAuth {
		username = ""
		password = ""
	}

	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{address},
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	return &Elastic{client: client, index: index}, err
}

func (e *Elastic) Search(indices string, body []byte, scroll bool) ([]byte, error) {
	opts := []func(*esapi.SearchRequest){
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(indices),
		e.client.Search.WithIgnoreUnavailable(true),
		e.client.Search.WithBody(bytes.NewBuffer(body)),
	}
	if scroll {
		opts = append(opts, e.client.Search.WithScroll(time.Minute))
	}

	response, err := e.client.Search(opts...)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, parseError(response)
	}

	return io.ReadAll(response.Body)
}

func (e *Elastic) Scroll(id string) ([]byte, error) {
	body, err := jsoniter.Marshal(map[string]string{
		"scroll_id": id,
	})
	if err != nil {
		return nil, err
	}

	response, err := e.client.Scroll(
		e.client.Scroll.WithContext(context.Background()),
		e.client.Scroll.WithBody(bytes.NewBuffer(body)),
		e.client.Scroll.WithScroll(time.Minute))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, parseError(response)
	}

	return io.ReadAll(response.Body)
}

func (e *Elastic) ClearScroll(scrollId string) {
	response, _ := e.client.ClearScroll(
		e.client.ClearScroll.WithContext(context.Background()),
		e.client.ClearScroll.WithScrollID(scrollId))
	defer response.Body.Close()
}

func (e *Elastic) GetTotalHitCount(v interface{}) int64 {
	f, _ := v.(float64)
	return int64(f)
}

func parseError(response *esapi.Response) error {
	var e versions.Error
	if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
		return err
	} else {
		// Print the response status and error information.
		if len(e.Details.RootCause) != 0 {
			return fmt.Errorf("type: %v, reason: %v", e.Details.Type, e.Details.RootCause[0].Reason)
		} else {
			return fmt.Errorf("type: %v, reason: %v", e.Details.Type, e.Details.Reason)
		}
	}
}
