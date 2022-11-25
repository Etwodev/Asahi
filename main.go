package main

import (
	"github.com/SpeedSlime/Covalence/server"
	"github.com/SpeedSlime/Covalence/server/middleware"
	"github.com/SpeedSlime/Covalence/server/router"
)

type Covalence struct {
	servers		[]*server.Server
}

func (c *Covalence) Find(name string) (*server.Server) {
	for _, server := range c.servers {
		if ( server.Name() == name  ) {
			return server
		}
	}
	return nil
}

func (c *Covalence) Create(version string, port string, name string, address string, connection string, experimental bool) (error) {
	if s := c.Find(name); s != nil {
		return // handle error
	}
	c.servers = append(c.servers, server.Create(version, port, name, address, connection, experimental))
}


func (c *Covalence) LoadRouters(name string, routers ...router.Router) {
	s := c.Find(name)
	if s == nil  {
		return // handle error
	}
	s.LoadRouters(routers...)
}

func (c *Covalence) LoadMiddlewares(name string, middlewares ...middleware.Middleware) {
	s := c.Find(name)
	if s == nil {
		return // handle error
	}
	s.LoadMiddleware(middlewares...)
}

func (c *Covalence) Start(name string, invoke func()) {
	if invoke != nil {
		invoke()
	}
	s := c.Find(name)
	if s == nil {
		return // handle error
	}
	s.Start()
}

func (c *Covalence) Stop(name string) {
	s := c.Find(name)
	if s == nil {
		return // handle error
	}
	s.Stop()
}