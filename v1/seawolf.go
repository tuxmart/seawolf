package seawolf

import (
	"context"
	"fmt"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
	"github.com/tuxmart/seawolf/internal"
	"google.golang.org/grpc"
)

type SeaWolf struct {
	client    filer_pb.SeaweedFilerClient
	Listeners []internal.FileListener
}

type Option func(*SeaWolf)

func WithListener(listener internal.FileListener) Option {
	return func(s *SeaWolf) {
		s.Listeners = append(s.Listeners, listener)
	}
}

func New(address string, opts ...Option) *SeaWolf {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := filer_pb.NewSeaweedFilerClient(conn)

	wolf := &SeaWolf{
		client:    client,
		Listeners: make([]internal.FileListener, 0),
	}

	for _, opt := range opts {
		opt(wolf)
	}

	return wolf
}

func (wolf *SeaWolf) Run() error {
	fmt.Println("Listening file metadata...")
	for {
		req, err := wolf.client.SubscribeMetadata(context.TODO(), &filer_pb.SubscribeMetadataRequest{
			SinceNs: time.Now().UnixNano(),
		})
		if err != nil {
			return err
		}

		event, err := req.Recv()
		if err != nil {
			return err
		}

		if notification := event.EventNotification; notification != nil {
			if event.GetDirectory() != notification.GetNewParentPath() {
				if notification.OldEntry != nil && notification.NewEntry != nil {
					for _, listener := range wolf.Listeners {
						listener.Move(event)
					}
				}
			} else if notification.NewEntry != nil && notification.OldEntry == nil {
				for _, listener := range wolf.Listeners {
					listener.Create(event)
				}
			} else if notification.NewEntry == nil && notification.OldEntry != nil {
				for _, listener := range wolf.Listeners {
					listener.Delete(event)
				}
			} else if notification.NewEntry != nil && notification.OldEntry != nil {
				for _, listener := range wolf.Listeners {
					listener.Update(event)
				}
			} else {
				return fmt.Errorf("unexpected event %v", event)
			}
		}
	}
}

func (wolf *SeaWolf) Client() filer_pb.SeaweedFilerClient {
	return wolf.client
}
