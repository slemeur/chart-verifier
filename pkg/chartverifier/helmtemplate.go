/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package chartverifier

import (
	"bufio"
	"io/ioutil"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/redhat-certification/chart-verifier/pkg/helm/actions"
)

func getImages(chartUri string) ([]string, error) {

	actionConfig := &action.Configuration{
		Releases:     nil,
		KubeClient:   &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Capabilities: chartutil.DefaultCapabilities,
		Log:          func(format string, v ...interface{}) {},
	}
	mem := driver.NewMemory()
	mem.SetNamespace("TestNamespace")
	actionConfig.Releases = storage.Init(mem)

	var m map[string]interface{}
	var images []string

	txt, err := actions.RenderManifests("testRelease", chartUri, m, actionConfig)
	if err == nil {
		r := strings.NewReader(txt)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ReplaceAll(line, " ", "")
			if strings.HasPrefix(line, "image:") {
				images = append(images, strings.Trim(strings.TrimLeft(line, "image:"), "\""))
			}
		}
	}

	return images, err
}
