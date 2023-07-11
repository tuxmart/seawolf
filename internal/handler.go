package internal

import "github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"

type FileListener interface {
	Create(ev *filer_pb.SubscribeMetadataResponse)
	Delete(ev *filer_pb.SubscribeMetadataResponse)
	Update(ev *filer_pb.SubscribeMetadataResponse)
	Move(ev *filer_pb.SubscribeMetadataResponse)
}
