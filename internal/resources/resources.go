package resources

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

// File represents a helm chart file
type File struct {
	TypeURL string
	Value   []byte
}

// ReleaseResource represents a helm2 or helm3 release
type ReleaseResource struct {
	Name      string
	Namespace string
	Templates []Template
	Files     []File
	Values    string
	Chart     *ChartMetadata
	Info      *ReleaseInfo
	Manifest  string
	Version   int32
}

// Template represents a helm chart template
type Template struct {
	Name string
	Data []byte
}

// Maintainer represents a maintainer (specified in chart.yaml)
type Maintainer struct {
	Name  string
	Email string
	URL   string
}

// ChartMetadata represents data about a chart
type ChartMetadata struct {
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
	Maintainers []Maintainer
	// The name of the template engine to use. Defaults to 'gotpl'.
	Engine string
	// The URL to an icon file.
	Icon string
	// The API Version of this chart.
	APIVersion string
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

// ReleaseInfo contains information about a Helm Release
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
	StatusID string
	// Contains the rendered templates/NOTES.txt if available
	Notes string
	// LastTestSuiteRun provides results on the last test run on a release
	// LastTestSuiteRun *TestSuite `protobuf:"bytes,5,opt,name=last_test_suite_run,json=lastTestSuiteRun,proto3" json:"last_test_suite_run,omitempty"`
}
