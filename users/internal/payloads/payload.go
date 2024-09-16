package payloads

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type Payload struct {
	Id string
}

func ExtractPayload(ctx context.Context) *Payload {
	payload := &Payload{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	for key, values := range md {
		switch key {
		case "id":
			payload.Id = values[0]
		}
	}

	return payload
}
