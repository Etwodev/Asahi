package main

import (
	"fmt"
	"github.com/SpeedSlime/Covalence/server"
	"github.com/SpeedSlime/Covalence/server/middleware"
	"github.com/SpeedSlime/Covalence/server/router"
)

type NotFoundError struct {
    Name string
    Msg  string
}

func (e *NotFoundError) Error() string { 
    return e.Msg
}

type AlreadyExistsError struct {
    Name string
	Port string
    Msg  string
}

func (e *AlreadyExistsError) Error() string { 
    return e.Msg
}

type Network struct {
	servers		[]*server.Server
}

func (n *Network) Find(name string) (*server.Server, error) {
	for _, server := range n.servers {
		if ( server.Name() == name  ) {
			return server, nil
		}
	}
	return nil, &NotFoundError{
		Msg: fmt.Sprintf("server name '%s' could not be found", name),
		Name: name,
	}
}

func (n *Network) Check(name string, port string) (error) {
	for _, server := range n.servers {
		if ( server.Name() == name || server.Port() == port) {
			return &AlreadyExistsError{
				Msg: fmt.Sprintf("server name '%s' or port '%s' already in use", name, port),
				Name: name, Port: port,
			}
		}
	}
	return nil 
}

func (n *Network) Create(version string, port string, name string, address string, connection string, experimental bool) (error) {
	if err := n.Check(name, port); err != nil {
		return fmt.Errorf("Create: could not create server: %w", err)
	}
	n.servers = append(n.servers, server.Create(version, port, name, address, connection, experimental))
	return nil
}


func (n *Network) LoadRouters(name string, routers ...router.Router) (error) {
	s, err := n.Find(name)
	if err != nil {
		return fmt.Errorf("LoadRouters: failed to find server: %w", err)
	}
	s.LoadRouters(routers...)
	return nil
}

func (n *Network) LoadMiddlewares(name string, middlewares ...middleware.Middleware) (error) {
	s, err := n.Find(name)
	if err != nil {
		return fmt.Errorf("LoadMiddlewares: failed to find server: %w", err)
	}
	s.LoadMiddleware(middlewares...)
	return nil
}

func (n *Network) Start(name string, invoke func()) (error) {
	if invoke != nil {
		invoke()
	}
	s, err := n.Find(name)
	if err != nil || s.Status() {
		return fmt.Errorf("Start: failed to find valid server: %w", err)
	}
	s.Start()
	return nil
}

func (n *Network) Stop(name string) (error) {
	s, err := n.Find(name)
	if err != nil {
		return fmt.Errorf("Stop: failed to find server: %w", err)
	}
	err = s.Stop()
	if err != nil {
		return fmt.Errorf("Stop: failed to stop server: %w", err)
	}
	return nil
}