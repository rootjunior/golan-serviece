package interfaces

import "context"

type IWorkerPool interface {
	StartProcessEvents(ctx context.Context)
}
