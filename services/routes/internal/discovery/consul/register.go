package consul

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type Registrar struct {
	client *consulapi.Client
	id     string
}

func NewRegistrar() (*Registrar, error) {
	cfg := consulapi.DefaultConfig()

	// CONSUL_HTTP_ADDR puede venir como "http://127.0.0.1:8500" o "127.0.0.1:8500"
	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" {
		if u, err := url.Parse(addr); err == nil && u.Host != "" {
			cfg.Address = u.Host // nos quedamos con host:port
		} else {
			// fallback: quitamos posible esquema manualmente
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
		healthPath = "/health"
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

// Helpers para leer env y convertir PORT
func EnvPort() int {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8000"
	}
	n, _ := strconv.Atoi(p)
	return n
}

// Mantener vivo (opcional)
func SleepTillExit() {
	for {
		time.Sleep(24 * time.Hour)
	}
}
