package grace

import (
	"context"
	"fmt"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/cloudflare/tableflip"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Options struct {
	UpgradeTimeout  time.Duration
	ShutdownTimeout time.Duration
	Network         string
}

type Server interface {
	Addr() string
	Serve(ln net.Listener) error
	Shutdown(ctx context.Context) error
}

type Grace interface {
	Serve(ctx context.Context, server Server)
	Stop()
	ListenForUpgrade(sigToResets ...os.Signal)
}

type gc struct {
	opt      Options
	Upgrader *tableflip.Upgrader
	Logger   logger.Logger
}

func (g *gc) Serve(ctx context.Context, server Server) {
	ln, err := g.Upgrader.Listen(g.opt.Network, server.Addr())
	if err != nil {
		g.Logger.Panic(fmt.Sprintf("Can't listen: %s", err.Error()))
	}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			g.Logger.Panic(err.Error())
		}
	}()

	g.Logger.Info("Ready")
	if err := g.Upgrader.Ready(); err != nil {
		err = errors.PropagateWithCode(err, errors.ErrCodeReadyStateFiled, "[Error] Ready state failed.")
		g.Logger.Panic(err.Error())
	}

	g.Logger.Info(fmt.Sprintf("Listening to %s", server.Addr()))

	<-g.Upgrader.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	time.AfterFunc(g.opt.ShutdownTimeout, func() {
		g.Logger.Info("Graceful shutdown timed out")
		os.Exit(1)
	})

	if err = server.Shutdown(ctx); err != nil {
		g.Logger.Warn(err.Error())
	}
}

// ListenForUpgrade listen for signals to do an upgrade
func (g *gc) ListenForUpgrade(sigToResets ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sigToResets...)
	for range c {
		if err := g.Upgrader.Upgrade(); err != nil {
			err = errors.PropagateWithCode(err, errors.ErrCodeUpgradeFailed, "[Error] Grace Upgrade Failed.")
			g.Logger.Error(err.Error())
		}
	}
}

func (g *gc) Stop() {
	g.Upgrader.Stop()
}

func New(Logger logger.Logger, opt Options) (Grace, error) {
	upg, err := tableflip.New(tableflip.Options{UpgradeTimeout: opt.UpgradeTimeout})
	if err != nil {
		return nil, err
	}

	if opt.Network == "" {
		opt.Network = "tcp"
	}

	return &gc{
		Upgrader: upg,
		opt:      opt,
		Logger:   Logger,
	}, nil
}
