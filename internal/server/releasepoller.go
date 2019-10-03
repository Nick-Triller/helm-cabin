package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	"log"
	"time"
)

var releaseStatuses = []release.Status_Code{
	release.Status_UNKNOWN,
	release.Status_DEPLOYED,
	release.Status_DELETED,
	release.Status_SUPERSEDED,
	release.Status_FAILED,
	release.Status_DELETING,
	release.Status_PENDING_INSTALL,
	release.Status_PENDING_UPGRADE,
	release.Status_PENDING_ROLLBACK,
}

// NewContext creates a versioned context.
func NewContext() context.Context {
	md := metadata.Pairs("x-helm-api-client", version.GetVersion())
	return metadata.NewOutgoingContext(context.TODO(), md)
}

// PollReleases uses helm client to poll releases
func PollReleases(releasesChan chan *services.ListReleasesResponse, tillerReachableChan chan bool, tillerAddress *string) {
	pollSleep := 6 * time.Second

	// target := "tiller-deploy.svc.kube-system.cluster.local:44134"
	// target = "localhost:8888"

	for {
		conn, err := grpc.Dial(*tillerAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect: %v\n", err)
		}
		client := services.NewReleaseServiceClient(conn)
		listReleaseRequest := &services.ListReleasesRequest{
			Limit:                9999999999,
			Offset:               "",
			// Reverse chronological
			SortBy:               services.ListSort_LAST_RELEASED,
			Filter:               "",
			// Reverse chronological
			SortOrder:            services.ListSort_DESC,
			// Releases with any status
			StatusCodes:          releaseStatuses,
			// Releases with any namespace
			Namespace:            "",
		}

		ctx := NewContext()
		listReleasesClient, err := client.ListReleases(ctx, listReleaseRequest)

		if err != nil {
			log.Printf("failed to create listReleasesClient: %v\n", err)
			tillerReachableChan <- false
			time.Sleep(2 * time.Second)
			continue
		}

		resp, err := listReleasesClient.Recv()
		if err != nil {
			log.Printf("failed to list releases: %v\n", err)
			tillerReachableChan <- false
			time.Sleep(2 * time.Second)
			continue
		}
		releasesChan <- resp
		tillerReachableChan <- true

		_ = conn.Close()
		time.Sleep(pollSleep)
	}
}
