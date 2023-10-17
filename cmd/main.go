package main

import (
	"flag"
	"os"

	"template/internal/server/http"
	"template/pkg/app"
	"template/pkg/config"
	"template/pkg/log"

	"github.com/gin-gonic/gin"
)

var (
	id, _      = os.Hostname() //nolint:errcheck
	Name       string
	mode       string
	configPath string
	port       int
	Version    string
)

func init() {
	flag.StringVar(&mode, "mode", "debug", "运行模式 debug/release")
	flag.StringVar(&configPath, "config", "./config/config.yaml", "配置文件路径")
	flag.IntVar(&port, "port", 8080, "端口号")
	flag.Parse()
}

func newApp(logger log.Logger, h *http.Server) *app.App {
	return app.New(
		app.ID(id),           // app id
		app.Name(Name),       // app name
		app.Version(Version), // app version
		app.Logger(logger),   // app logger
		app.Server(h),
	)
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

	a, f, err := wireApp(conf, logger)
	if err != nil {
		panic(err)
	}

	defer f()

	err = a.Run()
	if err != nil {
		panic(err)
	}
}
