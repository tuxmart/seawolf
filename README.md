## Sea Wolf

This is made for listening seaweedfs's file metadata changes.

## Install
```bash
go get github.com/tuxmart/seawolf
```

## Usage

```golang
package main

import (
	"github.com/tuxmart/seawolf/example"
	"github.com/tuxmart/seawolf/v1"
)

func main() {
	listener := &example.LogFileEventHandler{}
	if err := seawolf.Run("localhost:18888", seawolf.WithListener(listener)); err != nil {
		panic(err)
	}
}
```

You could implement your own file listeners by implementing `FileListener` interface.

```golang
type FileListener interface {
	Create(ev *filer_pb.SubscribeMetadataResponse)
	Delete(ev *filer_pb.SubscribeMetadataResponse)
	Update(ev *filer_pb.SubscribeMetadataResponse)
	Move(ev *filer_pb.SubscribeMetadataResponse)
}
```

## Reference
+ [Filer Change Data Capture](https://github.com/seaweedfs/seaweedfs/wiki/Filer-Change-Data-Capture)