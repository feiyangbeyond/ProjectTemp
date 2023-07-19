package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"deviceback/v3/pkg/config"
	"deviceback/v3/pkg/log"
	"deviceback/v3/pkg/server"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	mode       string
	configPath string
	port       int
)

func init() {
	flag.StringVar(&mode, "mode", "debug", "运行模式 debug/release")
	flag.StringVar(&configPath, "config", "./config/config.yaml", "配置文件路径")
	flag.IntVar(&port, "port", 8080, "端口号")
	flag.Parse()
}

type App struct {
	s *server.Server
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.s.Stop(ctx)
}

func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(a.s.Start)

	c := make(chan os.Signal, 1)
	// go 不允许监听 kill stop 信号
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func newApp(logger log.Logger, s *server.Server) *App {
	return &App{s: s}
}

func main() {
	gin.SetMode(mode)

	conf, err := config.Scan(configPath)
	if err != nil {
		panic(err)
	}

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	app, f, err := wireApp(conf, logger)
	if err != nil {
		panic(err)
	}

	defer f()

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
