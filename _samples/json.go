// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
)

const IndexName = "go-test-index1"

func main() {
	if err := example(); err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
}

func example() error {
	client, err := opensearchapi.NewDefaultClient()
	if err != nil {
		return err
	}

	///

	infoRequest, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return err
	}

	infoResponse, err := client.Client.Perform(infoRequest)
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		return err
	}
	fmt.Printf("client info: %s\n", resBody)

	///

	var index_body = strings.NewReader(`{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 1
        }
    },
    "mappings": {
        "properties": {
            "title": {
                "type": "text"
            },
            "year": {
                "type": "integer"
            }
        }
    }
}`)

	createIndexRequest, err := http.NewRequest("PUT", "/movies", index_body)
	if err != nil {
		return err
	}
	createIndexRequest.Header["Content-Type"] = []string{"application/json"}
	createIndexResp, err := client.Client.Perform(createIndexRequest)
	if err != nil {
		return err
	}
	createIndexRespBody, err := ioutil.ReadAll(createIndexResp.Body)
	if err != nil {
		return err
	}
	fmt.Println("create index: ", string(createIndexRespBody))

	///

	query := strings.NewReader(`{
    "size": 5,
    "query": {
        "multi_match": {
        "query": "miller",
        "fields": ["title^2", "director"]
        }
    }
}`)
	searchRequest, err := http.NewRequest("POST", "/movies/_search", query)
	if err != nil {
		return err
	}
	searchRequest.Header["Content-Type"] = []string{"application/json"}
	searchResp, err := client.Client.Perform(searchRequest)
	if err != nil {
		return err
	}
	searchRespBody, err := ioutil.ReadAll(searchResp.Body)
	if err != nil {
		return err
	}
	fmt.Println("search: ", string(searchRespBody))

	///

	deleteIndexRequest, err := http.NewRequest("DELETE", "/movies", nil)
	if err != nil {
		return err
	}
	deleteIndexResp, err := client.Client.Perform(deleteIndexRequest)
	if err != nil {
		return err
	}
	deleteIndexRespBody, err := ioutil.ReadAll(deleteIndexResp.Body)
	if err != nil {
		return err
	}
	fmt.Println("delete index: ", string(deleteIndexRespBody))

	return nil
}
