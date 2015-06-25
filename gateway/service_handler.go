package gateway

import (
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/gateway/middleware"
	"github.com/backstage/maestro/gateway/transformer"
)

// ServiceHandler registers the handler, transformers and middlewares for the given
// service.
type ServiceHandler struct {
	handler      http.Handler
	service      *account.Service
	transformers []transformer.Transformer
	middlewares  []middleware.Middleware
}

// func (s *ServiceHandler) addMiddleware(m middleware.Middleware, mc *account.Plugin) {
// 	marshal, err := json.Marshal(mc.Config)
// 	if err != nil {
// 		log.Printf("Failed to register middleware `%s`. Error: %s", mc.Name, err)
// 		return
// 	}
// 	m.Configure(string(marshal))
// 	s.middlewares = append(s.middlewares, m)
// 	log.Printf("Middleware `%s` added successfully for service `%s`.", mc.Name, s.service.Subdomain)
// }

// func (s *ServiceHandler) addTransformer(name string, t transformer.Transformer) {
// 	s.transformers = append(s.transformers, t)
// 	log.Printf("Transformer `%s` added successfully for service `%s`.", name, s.service.Subdomain)
// }
