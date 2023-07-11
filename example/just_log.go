package example

import (
	"fmt"

	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
)

type LogFileEventHandler struct {
}

func (l *LogFileEventHandler) Create(ev *filer_pb.SubscribeMetadataResponse) {
	// System.out.println("created entry " + event.getDirectory() + "/" + notification.getNewEntry().getName());
	fmt.Printf("Created entry %s%s\n", ev.GetDirectory(), ev.GetEventNotification().GetNewEntry().GetName())
}

func (l *LogFileEventHandler) Delete(ev *filer_pb.SubscribeMetadataResponse) {
	// 		System.out.println("deleted entry " + event.getDirectory() + "/" + notification.getOldEntry().getName())
	fmt.Printf("Deleted entry %s%s\n", ev.GetDirectory(), ev.GetEventNotification().GetOldEntry().GetName())
}

func (l *LogFileEventHandler) Update(ev *filer_pb.SubscribeMetadataResponse) {
	// System.out.println("updated entry " + event.getDirectory() + "/" + notification.getNewEntry().getName())
	fmt.Printf("Updated entry %s%s\n", ev.GetDirectory(), ev.GetEventNotification().GetNewEntry().GetName())
}

func (l *LogFileEventHandler) Move(ev *filer_pb.SubscribeMetadataResponse) {
	fmt.Printf("Moved %s%s to %s%s", ev.GetDirectory(), ev.GetEventNotification().GetOldEntry().GetName(), ev.GetEventNotification().GetNewParentPath(), ev.GetEventNotification().GetNewEntry().GetName())
}
