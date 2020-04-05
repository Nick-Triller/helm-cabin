package helm3

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/Nick-Triller/helm-cabin/internal/resources"
	"github.com/Nick-Triller/helm-cabin/internal/settings"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"encoding/json"
	rspb "helm.sh/helm/v3/pkg/release"
)

func PollReleases(releasesChan chan []resources.ReleaseListResource, settings *settings.Settings) {
	connectKubernetes()
	pollSleep := 6 * time.Second
	for {
		resp, err := clientset.CoreV1().Secrets("").List(metav1.ListOptions{
			FieldSelector: "type=helm.sh/release.v1",
			Watch:         false,
		})

		if err != nil {
			log.Warningf("Failed to retrieve releases: %v", err)
		} else {
			releasesChan <- convertResponseToReleaseListResources(resp)
		}

		time.Sleep(pollSleep)
	}
}

func convertResponseToReleaseListResources(resp *v1.SecretList) []resources.ReleaseListResource {
	releaseResources := make([]resources.ReleaseListResource, len(resp.Items))
	for i, secret := range resp.Items {
		helm3Release, err := decodeRelease(string(secret.Data["release"]))
		if err != nil {
			log.Warn("Failed to decode helm release secret, skipping")
			continue
		}
		releaseResources[i] = releaseListResourceFrom(helm3Release)
	}
	return releaseResources
}

func releaseListResourceFrom(r *rspb.Release) resources.ReleaseListResource {
	files := make([]resources.File, len(r.Chart.Files))
	for i, helm3ChartFile := range r.Chart.Files {
		template := resources.File{
			TypeUrl: helm3ChartFile.Name,
			Value:   helm3ChartFile.Data,
		}
		files[i] = template
	}

	templates := make([]resources.Template, len(r.Chart.Templates))
	for i, helm3Template := range r.Chart.Templates {
		template := resources.Template{
			Name: helm3Template.Name,
			Data: helm3Template.Data,
		}
		templates[i] = template
	}

	maintainers := make([]resources.Maintainer, len(r.Chart.Metadata.Maintainers))
	for i, helm3Mantainer := range r.Chart.Metadata.Maintainers {
		maintainer := resources.Maintainer{
			Name:  helm3Mantainer.Name,
			Email: helm3Mantainer.Email,
			Url:   helm3Mantainer.URL,
		}
		maintainers[i] = maintainer
	}

	return resources.ReleaseListResource{
		Name:      r.Name,
		Namespace: r.Namespace,
		Templates: templates,
		Files:     files,
		Values:    mapToJson(r.Chart.Values),
		Chart: &resources.ChartMetadata{
			Name:          r.Chart.Metadata.Name,
			Home:          r.Chart.Metadata.Home,
			Sources:       r.Chart.Metadata.Sources,
			Version:       r.Chart.Metadata.Version,
			Description:   r.Chart.Metadata.Description,
			Keywords:      r.Chart.Metadata.Keywords,
			Maintainers:   maintainers,
			Engine:        "",
			Icon:          r.Chart.Metadata.Icon,
			APIVersion:    r.Chart.Metadata.APIVersion,
			Condition:     r.Chart.Metadata.Condition,
			Tags:          r.Chart.Metadata.Tags,
			AppVersion:    r.Chart.Metadata.AppVersion,
			Deprecated:    r.Chart.Metadata.Deprecated,
			// Helm 3 has no tiller
			TillerVersion: "",
			Annotations:   r.Chart.Metadata.Annotations,
			KubeVersion:   r.Chart.Metadata.KubeVersion,
		},
		Info: &resources.ReleaseInfo{
			Status: &resources.Status{
				StatusID:  r.Info.Status.String(),
				Notes:     r.Info.Notes,
			},
			FirstDeployed: &timestamp.Timestamp{
				Seconds:              r.Info.FirstDeployed.Unix(),
				Nanos:                0,
			},
			LastDeployed:  &timestamp.Timestamp{
				Seconds:              r.Info.LastDeployed.Unix(),
				Nanos:                0,
			},
			Deleted:       &timestamp.Timestamp{
				Seconds:              r.Info.Deleted.Unix(),
				Nanos:                0,
			},
			Description:   r.Info.Description,
		},
		Manifest: r.Manifest,
		Version:  int32(r.Version),
	}
}

func mapToJson(structured map[string]interface{}) string {
	data, err := json.MarshalIndent(structured, "", "  ")
	if err != nil {
		log.Warn("Failed to map chart values to json")
		return ""
	}
	return string(data)
}

var magicGzip = []byte{0x1f, 0x8b, 0x08}
var b64 = base64.StdEncoding

func decodeRelease(data string) (*rspb.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls rspb.Release
	// unmarshal release object bytes
	if err := json.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}

