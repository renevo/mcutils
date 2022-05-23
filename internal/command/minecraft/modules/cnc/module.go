package cnc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/internal/control"
	"github.com/renevo/rpc"
)

type cfg struct {
	ControlAddr  string `hcl:"control_address,optional"`
	ControlToken string `hcl:"control_token,optional"`
}

type module struct {
	cfg      *cfg
	server   *rpc.Server
	listener net.Listener
	wg       sync.WaitGroup
}

func New() application.Module {
	m := &module{
		cfg: &cfg{},
	}

	return m
}

func (m *module) Config() (interface{}, error) {
	return m.cfg, nil
}

func (m *module) Start(ctx context.Context) error {
	log := ext.Logger(ctx)

	if m.cfg.ControlAddr == "" {
		log.Debug("Control Server Not Enabled")
		return nil
	}

	srv := ext.Minecraft(ctx)

	m.server = rpc.NewServer()
	if err := m.server.Register(&control.Minecraft{Server: srv}); err != nil {
		return errors.Wrap(err, "failed to initialize minecraft control server")
	}

	ln, err := net.Listen("tcp", m.cfg.ControlAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to open control server: %q", m.cfg.ControlAddr)
	}
	m.listener = ln

	log.Infof("Control Server Running on %v", ln.Addr().String())

	// logging middleware
	m.server.Use(func(next rpc.MiddlewareHandler) rpc.MiddlewareHandler {
		return func(ctx context.Context, rw rpc.ResponseWriter, req *rpc.Request) {
			start := time.Now()
			log.Infof("Control Server executing %q", req.ServiceMethod)
			next(ctx, rw, req)
			if rw.Err() == nil {
				log.Infof("Control Server executed %q in %v", req.ServiceMethod, time.Since(start))
			} else {
				log.Warnf("Control Server executed %q with error %v in %v", req.ServiceMethod, rw.Err(), time.Since(start))
			}
		}
	})

	// auth middleware
	m.server.Use(func(next rpc.MiddlewareHandler) rpc.MiddlewareHandler {
		return func(ctx context.Context, rw rpc.ResponseWriter, req *rpc.Request) {
			token := rpc.ContextHeader(ctx).Get(ext.HeaderToken)
			if token != m.cfg.ControlToken {
				rw.WriteError(errors.New("invalid minecraft control token"))
				return
			}

			next(ctx, rw, req)
		}
	})

	m.wg.Add(1)
	go m.acceptConnections(ctx)

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	if m.listener != nil {
		log := ext.Logger(ctx)
		log.Info("Control Server shutting down")

		defer log.Info("Control Server shutdown")

		_ = m.listener.Close()
	}

	m.wg.Wait()

	return nil
}

func (m *module) acceptConnections(ctx context.Context) {
	m.server.Accept(ctx, m.listener)
	m.wg.Done()
}
