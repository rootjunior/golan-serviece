package interfaces

import "context"

type WorkerPool interface {
	StartProcessEvents(ctx context.Context)
}
