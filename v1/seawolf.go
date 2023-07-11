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
	Listeners []internal.FileListener
}

type Option func(*SeaWolf)

func WithListener(listener internal.FileListener) Option {
	return func(s *SeaWolf) {
		s.Listeners = append(s.Listeners, listener)
	}
}

func Run(address string, opts ...Option) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	wolf := &SeaWolf{
		Listeners: make([]internal.FileListener, 0),
	}

	for _, opt := range opts {
		opt(wolf)
	}

	client := filer_pb.NewSeaweedFilerClient(conn)

	eventListeners := make([]internal.FileListener, 0)

	fmt.Println("Listening file metadata...")
	for {
		req, err := client.SubscribeMetadata(context.TODO(), &filer_pb.SubscribeMetadataRequest{
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
					for _, listener := range eventListeners {
						listener.Move(event)
					}
				}
			} else if notification.NewEntry != nil && notification.OldEntry == nil {
				for _, listener := range eventListeners {
					listener.Create(event)
				}
			} else if notification.NewEntry == nil && notification.OldEntry != nil {
				for _, listener := range eventListeners {
					listener.Delete(event)
				}
			} else if notification.NewEntry != nil && notification.OldEntry != nil {
				for _, listener := range eventListeners {
					listener.Update(event)
				}
			} else {
				return fmt.Errorf("unexpected event %v", event)
			}
		}
	}
}
