package server

import (
	"io"

	"google.golang.org/grpc"
	"k8s.io/helm/pkg/proto/hapi/services"
	log "k8s.io/klog"
)

var conn *grpc.ClientConn
var client services.ReleaseServiceClient

func connectTiller(settings *Settings) {
	var err error = nil
	conn, err = grpc.Dial(*settings.TillerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create connection: %v\n", err)
	}
	client = services.NewReleaseServiceClient(conn)
}

func getReleasesList(listReleaseRequest *services.ListReleasesRequest) (*services.ListReleasesResponse, error) {
	ctx := NewContext()

	listReleasesClient, err := client.ListReleases(ctx, listReleaseRequest)
	if err != nil {
		return nil, err
	}
	resp, err := listReleasesClient.Recv()
	if err != nil {
		if err == io.EOF {
			log.V(1).Info("Received EOF, no releases exist")
			// EOF if no releases exist
			emptyResp := &services.ListReleasesResponse{
				Count:    0,
				Next:     "",
				Total:    0,
				Releases: nil,
			}
			resp = emptyResp
		} else {
			return nil, err
		}
	}

	return resp, nil
}
