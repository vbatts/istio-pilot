// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
)

type mongo struct {
	*infra
	//s *mgo.Session
}

func (t *mongo) String() string {
	return "mongodb"
}

func (t *mongo) setup() error {
	if !t.Mongo {
		return nil
	}

	return nil
}

func (t *mongo) teardown() {
}

func (t *mongo) run() error {
	if !t.Mongo {
		return nil
	}

	f := func() status {
		// once to check the ratings collection
		response := t.clientRequest(
			"t",
			fmt.Sprintf("mongodb://mongo.%s:27017", t.IstioNamespace),
			1,
			"",
		)

		if len(response.code) == 0 || response.code[0] != httpOk {
			//return errAgain
			return errors.New(response.body)
		}

		// then iterate to insert into a iteration collection
		for i := 0; i < 10; i++ {
			response := t.clientRequest(
				"t",
				fmt.Sprintf("mongodb://mongo.%s:27017", t.IstioNamespace),
				1,
				fmt.Sprintf("%d", i),
			)

			if len(response.code) == 0 || response.code[0] != httpOk {
				//return errAgain
				return errors.New(response.body)
			}

			glog.Infof("%#v", response)
		}
		return nil
	}

	return parallel(map[string]func() status{
		"iterate on mongo inserts": f,
	})

	//resp := t.clientRequest(src, url, 1, "")
	// Soo, the clientRequest needs to be given the url where to make a request to,
	// and then we check that it actually succeeded
}
