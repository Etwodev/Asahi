package main

import (
	"github.com/SpeedSlime/Covalence/server"
	"github.com/SpeedSlime/Covalence/server/router"
	"github.com/SpeedSlime/Covalence/server/middleware"
)

type Covalence struct {
	servers		[]*server.Server
}

func New() {
	return Covalence{}
}

func (c *Covalence) Create(version string, port string, name string, address string, connection string, experimental bool) (error) {
	for _, server := range c.servers {
		if ( server.Name() == name || server.Port() == port ) {
			return // raise error
		}
	}
	c.servers = append(c.servers, server.Create(version, port, name, address, connection, experimental))
}

func (c *Covalence) LoadRouters(name string, routers ...router.Router) {
	for i, server := range c.servers {
		if ( server.Name() == name ) {
			c.servers[i].Load(routers...)
		}
	}
}

func (c *Covalence) LoadMiddleware(name string, middlewares ...middleware.Middleware) {
	for i, server := range c.servers {
		if ( server.Name() == name ) {
			c.servers[i].Load(middlewares...)
		}
	}
}

