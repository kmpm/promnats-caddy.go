package promnatscaddy

import (
	"strings"
	"time"

	"github.com/kmpm/promnats.go"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type PromNats struct {
	ContextName string
	ServerURL   string
	Subject     string
	Interval    time.Duration
	logger      *zap.Logger
	nc          *nats.Conn
	routes      map[string][]byte
}

func (m *PromNats) refresh() {
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				m.logger.Info("tick")
				err := m.request()
				if err != nil {
					m.logger.Error("error requesting data", zap.Error(err))
				}
			case <-quit:
				m.logger.Info("stop")
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *PromNats) request() error {
	ctx := context.TODO()
	msgs, err := doReq(ctx, nil, "metrics", 0, m.nc)
	if err != nil {
		return err
	}
	keys := []string{}
	routes := map[string][]byte{}
	folders := map[string][]string{}
	for _, msg := range msgs {
		id := strings.Trim(msg.Header.Get(promnats.HeaderPnID), ". ")
		if id == "" {
			m.logger.Warn("response without header", zap.Any("header", msg.Header), zap.Int("length", len(msg.Data)))
			continue
		}
		// TODO: make something out of it,
		// sort by id, split by dot and make available by path
		// also add some kind of index page to ease discovery
		parts := strings.Split(id, ".")
		if len(parts) < 3 {
			m.logger.Warn("id must have at least 3 parts", zap.String("id", id))
			continue
		}

		p := "/" + strings.Join(parts, "/")
		routes[p] = msg.Data
		keys = append(keys, p)
		for i := 0; i < len(parts); i++ {
			p = "/" + strings.Join(parts[:i], "/")
			if v, ok := folders[p]; ok {
				folders[p] = append(v, parts[i])
			} else {
				folders[p] = []string{parts[i]}
			}
		}
	}
	for k, v := range folders {
		routes[k] = []byte(strings.Join(v, ","))
		keys = append(keys, k)
	}
	m.routes = routes
	m.logger.Info("routes", zap.Strings("routes", keys), zap.Any("folders", folders))
	return nil
}
