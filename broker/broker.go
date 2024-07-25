package broker

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/server"
)

const (
	port = 8080
)

var (
	logger        domain.ILogger
	serverFactory domain.IServerFactory = &server.ServerFactory{}
)

func Init(ctx context.Context) error {
	environment := ctx.Value("environment").(*domain.Environment)
	logger = environment.Logger
	svr := serverFactory.MakeServer(ctx)
	return svr.Init(intializeRoutes, port)
}
