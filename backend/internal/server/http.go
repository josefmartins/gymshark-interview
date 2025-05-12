package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Server struct {
	server          *http.Server
	productService  ProductsService
	packagesService PackagesService
	api             huma.API
}

func (s Server) Start() {
	log.Println("server started. listening on port " + s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("server failed, ", err)
	}
	log.Println("server stopped")
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func New(port int, productService ProductsService, packagesService PackagesService) *Server {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Product Package Sizes API", "1.0.0"))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	s := &Server{
		api:             api,
		server:          httpServer,
		productService:  productService,
		packagesService: packagesService,
	}

	s.api.UseMiddleware(allowCORS)

	s.declareRoutes()

	return s
}

// allow server to be called by an external browser
func allowCORS(ctx huma.Context, next func(huma.Context)) {
	ctx.SetHeader("Access-Control-Allow-Origin", "*") // or specific origin
	ctx.SetHeader("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
	ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type")

	if ctx.Method() == http.MethodOptions {
		ctx.SetStatus(http.StatusNoContent)
		return
	}

	next(ctx)
}
