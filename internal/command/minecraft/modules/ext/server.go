package ext

import (
	"context"

	"github.com/renevo/mcutils/pkg/minecraft"
)

const (
	minecraftContext = contextKey("server")
)

func WithMinecraft(ctx context.Context, server *minecraft.Server) context.Context {
	return context.WithValue(ctx, minecraftContext, server)
}

func Minecraft(ctx context.Context) *minecraft.Server {
	server := ctx.Value(minecraftContext)
	if server == nil {
		return minecraft.Default()
	}

	return server.(*minecraft.Server)
}
