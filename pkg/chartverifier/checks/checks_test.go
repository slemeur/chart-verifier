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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestIsHelmV3(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "valid tarball", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		config := viper.New()
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, Helm3Reason, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "invalid tarball", uri: "chart-0.1.0-v2.invalid.tgz"},
	}

	for _, tc := range negativeTestCases {
		config := viper.New()
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, NotHelm3Reason, r.Reason[0])
		})
	}
}

func TestHasReadme(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with README", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasReadme(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ReadmeExist, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "chart with README", uri: "chart-0.1.0-v3.without-readme.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasReadme(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ReadmeDoesNotExist, r.Reason[0])
		})
	}
}

func TestContainsTest(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "tarball contains at least one test", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsTest(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ChartTestFilesExist, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "tarball contains at least one test", uri: "chart-0.1.0-v3.valid.notest.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsTest(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ChartTestFilesDoesNotExist, r.Reason[0])
		})
	}
}

func TestHasValuesSchema(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with values", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValuesSchema(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ValuesSchemaFileExist, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "chart without values", uri: "chart-0.1.0-v3.no-values-schema.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValuesSchema(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ValuesSchemaFileDoesNotExist, r.Reason[0])
		})
	}
}

func TestHasValues(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with values", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValues(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ValuesFileExist, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "chart without values", uri: "chart-0.1.0-v3.no-values.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValues(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ValuesFileDoesNotExist, r.Reason[0])
		})
	}
}

func TestHasMinKubeVersion(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "minimum Kubernetes version specified", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasMinKubeVersion(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, MinKuberVersionSpecified, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "minimum Kubernetes version not specified", uri: "chart-0.1.0-v3.without-minkubeversion.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasMinKubeVersion(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, MinKuberVersionNotSpecified, r.Reason[0])
		})
	}

}

func TestNotContainCRDs(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Not contain CRDs", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCRDs(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ChartDoesNotContainCRDs, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "Contain CRDs", uri: "chart-0.1.0-v3.with-crd.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCRDs(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ChartContainCRDs, r.Reason[0])
		})
	}
}

func TestNotContainCSIObjects(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Not contain CSI objects", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCSIObjects(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, CSIObjectsDoesNotExist, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "Contain CRDs", uri: "chart-0.1.0-v3.with-csi.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCSIObjects(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, CSIObjectsExist, r.Reason[0])
		})
	}
}

func TestHelmLint(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Helm lint works for valid chart", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HelmLint(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, HelmLintSuccessful, r.Reason[0])
		})
	}

	negativeTestCases := []testCase{
		{description: "Helm lint fails for invalid chart", uri: "chart-0.1.0-v2.invalid.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HelmLint(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Contains(t, r.Reason[0], HelmLintHasFailedPrefix)
		})
	}

}

func TestImageCertify(t *testing.T) {

	type testCase struct {
		description string
		uri         string
		numErrors   int
	}

	negativeTestCases := []testCase{
		{description: "chart-0.1.0-v3.valid.tgz check images fails", uri: "chart-0.1.0-v3.valid.tgz", numErrors: 2},
		{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-crd.tgz", numErrors: 2},
		{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-csi.tgz", numErrors: 1},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ImagesAreCertified(tc.uri, config)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			for i := 0; i < tc.numErrors; i++ {
				require.Contains(t, r.Reason[i], ImageNotCertified)
			}
		})
	}

}

func TestImageParsing(t *testing.T) {

	type testCase struct {
		description      string
		image            string
		expectedVersion  string
		expectedRepo     string
		expectedRegistry string
	}

	testCases := []testCase{
		{"Single repo Default version 1", "repo", "latest", "repo", ""},
		{"Single repo Default version 2", "repo:", "latest", "repo", ""},
		{"Single repo with version", "repo:1.1.8", "1.1.8", "repo", ""},
		{"Double repo with version", "repo/product:1.1.8", "1.1.8", "repo/product", ""},
		{"Registry, double repo with version", "registry/repo/product:1.1.8", "1.1.8", "repo/product", "registry"},
		{"Registry with port, double repo with version", "registry:8080/repo/product:1.1.8", "1.1.8", "repo/product", "registry:8080"},
	}

	for _, testCase := range testCases {
		registries, repository, version := getImageParts(testCase.image)

		//fmt.Println("Image : " + testCase.image)
		//if len(registries) > 0  {
		///	fmt.Println("    Registry : " + registries[0])
		//}
		//fmt.Println("    Repository : " + repository)
		//fmt.Println("    Version : " + version)

		require.Equal(t, repository, testCase.expectedRepo)
		require.Equal(t, version, testCase.expectedVersion)
		if len(registries) == 0 {
			require.True(t, len(testCase.expectedRegistry) == 0)
		} else {
			require.Equal(t, registries[0], testCase.expectedRegistry)
		}
	}
}

func TestPartnerCharts(t *testing.T) {

	tarballs, err := WalkMatch("/Users/martinmulholland/helm/Partner Helm Charts", "*.tgz")
	require.NoError(t, err)

	numTests := 0
	numPasses := 0
	numNoImages := 0

	for _, tarball := range tarballs {
		config := viper.New()
		numTests++
		r, err := ImagesAreCertified(tarball, config)
		if err == nil {
			if r.Ok {
				fmt.Println("\n\nWINNNER!!!")
				numPasses++
				if strings.Contains(r.Reason[0], "No images to certify") {
					numNoImages++
				}
			}
			fmt.Println(fmt.Sprintf(" %t  : %s", r.Ok, tarball))
			for _, reason := range r.Reason {
				fmt.Println("     Reason : " + reason)
			}
			if r.Ok {
				fmt.Print("\n\n")
			}
		} else {
			fmt.Println(fmt.Sprintf(" FAIL  : %s", tarball))
			fmt.Println(fmt.Sprintf("      error : %v", err))
		}
	}

	fmt.Println(fmt.Sprintf("Tests : %d, Passes : %d, No Images : %d", numTests, numPasses, numNoImages))

}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
