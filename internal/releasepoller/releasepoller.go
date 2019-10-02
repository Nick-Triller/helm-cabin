package releasepoller

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/tlsutil"
	"log"
	"os"
	"time"
)

var (
	tillerTunnel *kube.Tunnel
	settings     helm_env.EnvSettings
)

// PollReleases uses helm client to poll releases
func PollReleases(releasesChan chan *rls.ListReleasesResponse) {
	pollSleep := 10 * time.Second

	flags := (&cobra.Command{}).PersistentFlags()
	settings.AddFlags(flags)
	settings.AddFlagsTLS(flags)
	settings.Init(flags)
	settings.InitTLS(flags)

	settings.TillerConnectionTimeout = 10
	// TODO use tiller svc when running in cluster
	// settings.TillerHost = "tiller-deploy.svc.kube-system.cluster.local"

	trySetupConnection()
	client := newClient()

	for {
		// Get all releases
		limitOption := helm.ReleaseListLimit(9999999999)
		// Include all releases irrespective of status except superseded revisions
		listStatuses := helm.ReleaseListStatuses([]release.Status_Code{
			release.Status_UNKNOWN,
			release.Status_DEPLOYED,
			release.Status_DELETED,
			release.Status_SUPERSEDED,
			release.Status_FAILED,
			release.Status_DELETING,
			release.Status_PENDING_INSTALL,
			release.Status_PENDING_UPGRADE,
			release.Status_PENDING_ROLLBACK,
		})
		sort := helm.ReleaseListSort(int32(rls.ListSort_LAST_RELEASED))
		sortOrder := helm.ReleaseListOrder(int32(rls.ListSort_DESC))
		releases, err := client.ListReleases(limitOption, listStatuses, sort, sortOrder)
		if err != nil {
			log.Println(err)
		} else {
			releasesChan <- releases
		}
		time.Sleep(pollSleep)
	}
}

func newClient() helm.Interface {
	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(settings.TillerConnectionTimeout)}

	if settings.TLSVerify || settings.TLSEnable {
		tlsopts := tlsutil.Options{
			ServerName:         settings.TLSServerName,
			KeyFile:            settings.TLSKeyFile,
			CertFile:           settings.TLSCertFile,
			InsecureSkipVerify: true,
		}
		if settings.TLSVerify {
			tlsopts.CaCertFile = settings.TLSCaCertFile
			tlsopts.InsecureSkipVerify = false
		}
		tlscfg, err := tlsutil.ClientConfig(tlsopts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		options = append(options, helm.WithTLS(tlscfg))
	}
	return helm.NewClient(options...)
}

func trySetupConnection() {
	for {
		err := setupConnection()
		if err != nil {
			log.Println("Failed to connect with Tiller. Retrying in 10 seconds.", err)
		} else {
			return
		}
		time.Sleep(10 * time.Second)
	}
}

func setupConnection() error {
	if settings.TillerHost == "" {
		config, client, err := getKubeClient(settings.KubeContext, settings.KubeConfig)
		if err != nil {
			return err
		}

		tillerTunnel, err = portforwarder.New(settings.TillerNamespace, client, config)
		if err != nil {
			return err
		}

		settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
	}

	return nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}
