package actions

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func RunReleaseTesting(releaseName, namespace string, conf *action.Configuration) (*release.Release, error) {

	cmd := action.NewReleaseTesting(conf)
	cmd.Namespace = namespace
	release, err := cmd.Run(releaseName)
	if err != nil {
		return nil, err
	}
	return release, nil
}
