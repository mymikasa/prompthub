package ctxutil

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
)

type contextKey struct{}

func WithActor(ctx context.Context, actor *domain.Actor) context.Context {
	return context.WithValue(ctx, contextKey{}, actor)
}

func ActorFromCtx(ctx context.Context) *domain.Actor {
	actor, _ := ctx.Value(contextKey{}).(*domain.Actor)
	return actor
}
