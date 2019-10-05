package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
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

// PollReleases uses helm client to poll releasesCache
func PollReleases(releasesChan chan *services.ListReleasesResponse, tillerReachableChan chan bool, settings *Settings) {
	pollSleep := 6 * time.Second

	conn, err := grpc.Dial(*settings.TillerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create connection: %v\n", err)
	}
	client := services.NewReleaseServiceClient(conn)

	for {
		listReleaseRequest := &services.ListReleasesRequest{
			Limit:  9999999999,
			Offset: "",
			// Reverse chronological
			SortBy: services.ListSort_LAST_RELEASED,
			Filter: "",
			// Reverse chronological
			SortOrder: services.ListSort_DESC,
			// Releases with any status
			StatusCodes: releaseStatuses,
			// Releases with any namespace
			Namespace: "",
		}

		ctx := NewContext()
		listReleasesClient, err := client.ListReleases(ctx, listReleaseRequest)

		if err != nil {
			log.Printf("failed to create listReleasesClient: %v\n", err)
			onError(tillerReachableChan)
			continue
		}

		resp, err := listReleasesClient.Recv()
		if err == io.EOF {
			log.Println("Received EOF, no releases exist")
			// EOF if no releases exist
			emptyResp := &services.ListReleasesResponse{
				Count:                0,
				Next:                 "",
				Total:                0,
				Releases:             nil,
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			}
			err = nil
			resp = emptyResp
		}
		if err != nil {
			log.Printf("failed to list releases: %v\n", err)
			onError(tillerReachableChan)
			continue
		}

		releasesChan <- resp
		tillerReachableChan <- true

		time.Sleep(pollSleep)
	}
}

func onError(tillerReachableChan chan bool) {
	tillerReachableChan <- false
	time.Sleep(3 * time.Second)
}