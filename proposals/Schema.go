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
    helmchart "helm.sh/helm/v3/pkg/chart"
}


// Proposed schema for verifier report

type CheckType string

const (
	MandatoryCheckType    CheckType = "Mandatory"
	OptionalCheckType     CheckType = "Optional"
	ExperimentalCheckType CheckType = "Experimental"
	ExceptionCheckType    CheckType = "Exception"
)


type Certificate struct {
	Apiversion string         `json:"apiversion" yaml:"apiversion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   CertificateMetadata `json:"metadata" yaml:"metadata"`
	Results    []CheckReport  `json:"results" yaml:"results"`
}

type CertificateMetadata struct {
	Tool      ToolMetadata       `json:"tool" yaml:"tool"`
	ChartData helmchart.Metadata `json:"chart" yaml:"chart"`
	Overrides string             `json: "chart-overrides" yaml:"chart-overrides"`
}

type ToolMetadata struct {
	Version                    string `json:"verifier-version" yaml:"verifier-version"`
	ChartUri                   string `json:"chart-uri" yaml:"chart-uri"`
	Digest                     string `json:"digest" yaml:"digest"`
	LastCertifiedTime          string `json:"lastCertifiedTime" yaml:"lastCertifiedTime"`
	CertifiedOpenShiftVersions string `json:"certifiedOpenShiftVersions" yaml:"certifiedOpenShiftVersions"`
}

type CheckReport struct {
	Check   string    `json:"apiversion" yaml:"apiversion"`
	Type    CheckType `json:"type" yaml:"type"`
	Outcome string    `json:"outcome" yaml:"outcome"`
	Reason  string    `json:"reason" yaml:"reason"`
}
