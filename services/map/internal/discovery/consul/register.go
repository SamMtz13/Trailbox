package consul

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type Registrar struct {
	client *consulapi.Client
	id     string
}

func NewRegistrar() (*Registrar, error) {
	cfg := consulapi.DefaultConfig()
	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" {
		if u, err := url.Parse(addr); err == nil && u.Host != "" {
			cfg.Address = u.Host
		} else {
			cfg.Address = strings.TrimPrefix(strings.TrimPrefix(addr, "http://"), "https://")
		}
	}
	c, err := consulapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Registrar{client: c}, nil
}

func (r *Registrar) Register(serviceName, address string, port int, healthPath string) (string, error) {
	if healthPath == "" {
		healthPath = "/healthz"
	}
	if !strings.HasPrefix(healthPath, "/") {
		healthPath = "/" + healthPath
	}
	id := fmt.Sprintf("%s-%s-%d", serviceName, address, port)
	reg := &consulapi.AgentServiceRegistration{
		ID:      id,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", address, port, healthPath),
			Interval:                       "10s",
			Timeout:                        "2s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	if err := r.client.Agent().ServiceRegister(reg); err != nil {
		return "", err
	}
	log.Printf("[consul] registered service=%s id=%s addr=%s:%d", serviceName, id, address, port)
	r.id = id
	return id, nil
}

func (r *Registrar) Deregister() {
	if r.client == nil || r.id == "" {
		return
	}
	if err := r.client.Agent().ServiceDeregister(r.id); err != nil {
		log.Printf("[consul] deregister error: %v", err)
	} else {
		log.Printf("[consul] deregistered id=%s", r.id)
	}
}
