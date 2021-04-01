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
package checks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var pyxisBaseUrl = "https://catalog.redhat.com/api/containers/v1/repositories"

type RepositoriesBody struct {
	PyxisRepositories []PyxisRepository `json:"data"`
}

type PyxisRepository struct {
	Id          string `json:"_id"`
	Repository  string `json:"repository"`
	VendorLabel string `json:"vendor_label"`
	Registry    string `json:"registry"`
}

type RegistriesBody struct {
	PyxisRegistries []PyxisRegistry `json:"data"`
}

type PyxisRegistry struct {
	Id           string               `json:"_id"`
	Repositories []RegistryRepository `json:"repositories"`
}

type RegistryRepository struct {
	Registry   string          `json:"registry"`
	Repository string          `json:"repository"`
	Tags       []RepositoryTag `json:"tags"`
}

type RepositoryTag struct {
	Name string `json:"name"`
}

func getImageRegistries(repository string) ([]string, error) {
	var err error
	var registries []string

	requestFilter := fmt.Sprintf("filter=repository==%s", repository)
	requestUrl := fmt.Sprintf("%s?%s", pyxisBaseUrl, requestFilter)

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("X-API-KEY", "RedHatChartVerifier")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error getting repository %s : %v\n", repository, err))
	} else {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var repositoriesBody RepositoriesBody
			json.Unmarshal(body, &repositoriesBody)

			if len(repositoriesBody.PyxisRepositories) > 0 {
				for _, repo := range repositoriesBody.PyxisRepositories {
					registries = append(registries, repo.Registry)
				}
			} else {
				err = errors.New(fmt.Sprintf("Respository not found: %s", repository))
			}
		} else {
			err = errors.New(fmt.Sprintf("Bad response code from Pyxis: %d : %s", resp.StatusCode, requestUrl))
		}
	}

	return registries, err
}

func checkImageInRegistry(repository string, version string, registry string) (bool, error) {

	var err error
	found := false

	requestFilter := fmt.Sprintf("filter=repositories=em=(repository==%s;registry==%s)", repository, registry)
	requestUrl := fmt.Sprintf("%s/registry/%s/repository/%s/images?%s", pyxisBaseUrl, registry, repository, requestFilter)
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("X-API-KEY", "RedHatChartVerifier")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var registriesBody RegistriesBody
			json.Unmarshal(body, &registriesBody)

			if len(registriesBody.PyxisRegistries) > 0 {
				var tags []string
				found = false
				for _, reg := range registriesBody.PyxisRegistries {
					for _, repo := range reg.Repositories {
						if repo.Repository == repository && repo.Registry == registry {
							for _, tag := range repo.Tags {
								if tag.Name == version {
									found = true
									break
								} else {
									tags = append(tags, tag.Name)
								}
							}
						}
						if found {
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					err = errors.New(fmt.Sprintf("Version %s not found. Found : %s", version, strings.Join(tags, ", ")))
				}
			} else {
				err = errors.New(fmt.Sprintf("Registry not found: %s", registry))
			}
		} else {
			err = errors.New(fmt.Sprintf("Bad response code %d from pyxis request : %s", resp.StatusCode, requestUrl))
		}
	}

	return found, err
}
