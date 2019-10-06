package server

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
)

type releaseListResource struct {
	Name      string
	Namespace string
	Chart     *chartMetadata
	Info      *ReleaseInfo
	Version   int32
}

type chartMetadata struct {
	// The name of the chart
	Name string
	// The URL to a relevant project page, git repo, or contact person
	Home string
	// Source is the URL to the source code of this chart
	Sources []string
	// A SemVer 2 conformant version string of the chart
	Version string
	// A one-sentence description of the chart
	Description string
	// A list of string keywords
	Keywords []string
	// A list of name and URL/email address combinations for the maintainer(s)
	Maintainers []*chart.Maintainer
	// The name of the template engine to use. Defaults to 'gotpl'.
	Engine string
	// The URL to an icon file.
	Icon string
	// The API Version of this chart.
	ApiVersion string
	// The condition to check to enable chart
	Condition string
	// The tags to check to enable chart
	Tags string
	// The version of the application enclosed inside of this chart.
	AppVersion string
	// Whether or not this chart is deprecated
	Deprecated bool
	// TillerVersion is a SemVer constraints on what version of Tiller is required.
	// See SemVer ranges here: https://github.com/Masterminds/semver#basic-comparisons
	TillerVersion string
	// Annotations are additional mappings uninterpreted by Tiller,
	// made available for inspection by other applications.
	Annotations map[string]string
	// KubeVersion is a SemVer constraint specifying the version of Kubernetes required.
	KubeVersion string
}

type ReleaseInfo struct {
	Status        *Status
	FirstDeployed *timestamp.Timestamp
	LastDeployed  *timestamp.Timestamp
	// Deleted tracks when this object was deleted.
	Deleted *timestamp.Timestamp
	// Description is human-friendly "log entry" about this release.
	Description string
}

// Status defines the status of a release.
type Status struct {
	// This field differs from Helm structs. Contains the enum string of status code.
	StatusId string
	// Cluster resources as kubectl would print them.
	Resources string
	// Contains the rendered templates/NOTES.txt if available
	Notes string
	// LastTestSuiteRun provides results on the last test run on a release
	// LastTestSuiteRun *TestSuite `protobuf:"bytes,5,opt,name=last_test_suite_run,json=lastTestSuiteRun,proto3" json:"last_test_suite_run,omitempty"`
}

func releaseToResourceList(r *release.Release) *releaseListResource {
	resource := &releaseListResource{
		Name:      r.Name,
		Namespace: r.Namespace,
		Chart: &chartMetadata{
			Name:          r.Chart.Metadata.Name,
			Home:          r.Chart.Metadata.Home,
			Sources:       r.Chart.Metadata.Sources,
			Version:       r.Chart.Metadata.Version,
			Description:   r.Chart.Metadata.Description,
			Keywords:      r.Chart.Metadata.Keywords,
			Maintainers:   r.Chart.Metadata.Maintainers,
			Engine:        r.Chart.Metadata.Engine,
			Icon:          r.Chart.Metadata.Icon,
			ApiVersion:    r.Chart.Metadata.ApiVersion,
			Condition:     r.Chart.Metadata.Condition,
			Tags:          r.Chart.Metadata.Tags,
			AppVersion:    r.Chart.Metadata.ApiVersion,
			Deprecated:    r.Chart.Metadata.Deprecated,
			TillerVersion: r.Chart.Metadata.TillerVersion,
			Annotations:   r.Chart.Metadata.Annotations,
			KubeVersion:   r.Chart.Metadata.KubeVersion,
		},
		Info: &ReleaseInfo{
			Status: &Status{
				StatusId:  r.Info.Status.Code.String(),
				Resources: r.Info.Status.Resources,
				Notes:     r.Info.Status.Notes,
			},
			FirstDeployed: r.Info.FirstDeployed,
			LastDeployed:  r.Info.LastDeployed,
			Deleted:       r.Info.Deleted,
			Description:   r.Info.Description,
		},
		Version: r.Version,
	}
	return resource
}
