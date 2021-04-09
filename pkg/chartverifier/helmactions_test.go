package chartverifier

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTemplate(t *testing.T) {

	type testCase struct {
		description string
		uri         string
		images      []string
	}

	TestCases := []testCase{
		{description: "chart-0.1.0-v3.valid.tgz images ", uri: "checks/chart-0.1.0-v3.valid.tgz", images: []string{"nginx:1.16.0", "busybox"}},
		{description: "chart-0.1.0-v3.with-crd.tgz", uri: "checks/chart-0.1.0-v3.with-crd.tgz", images: []string{"nginx:1.16.0", "busybox"}},
		{description: "chart-0.1.0-v3.with-csi.tgz", uri: "checks/chart-0.1.0-v3.with-csi.tgz", images: []string{"nginx:1.16.0"}},
	}

	for _, tc := range TestCases {
		t.Run(tc.description, func(t *testing.T) {
			images, err := getImages(tc.uri)
			require.NoError(t, err)
			require.Equal(t, len(images), len(tc.images))
			for i := 0; i < len(tc.images); i++ {
				require.Contains(t, images, tc.images[i])
			}
		})
	}
}
