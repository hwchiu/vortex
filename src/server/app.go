package server

import (
	"log"
	"net"
	"net/http"

	"github.com/linkernetworks/logger"
	"github.com/hwchiu/vortex/src/config"
	"github.com/hwchiu/vortex/src/serviceprovider"
)

// App is the structure to set config & service provider of APP
type App struct {
	Config          config.Config
	ServiceProvider *serviceprovider.Container
}

// LoadConfig consumes a string of path to the json config file and read config file into Config.
func (a *App) LoadConfig(configPath string) *App {
	if configPath == "" {
		log.Fatal("-config option is required.")
	}

	a.Config = config.MustRead(configPath)
	return a
}

// Start consumes two strings, host and port, invoke service initilization and serve on desired host:port
func (a *App) Start(host, port string) error {

	a.InitilizeService()

	bind := net.JoinHostPort(host, port)

	return http.ListenAndServe(bind, a.AppRoute())
}

// InitilizeService weavering services with global variables inside server package
func (a *App) InitilizeService() {
	logger.Setup(a.Config.Logger)

	a.ServiceProvider = serviceprovider.New(a.Config)
}
