package router

import (
	"github.com/avito-shop-service/internal/handlers"
	"github.com/avito-shop-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(transactionHandler *handlers.TransactionHandler, userHandler *handlers.UserHandler, buyHandler *handlers.BuyHandler, infoHandler *handlers.InformationHandler) http.Handler {
	r := chi.NewRouter()
	r.With(middleware.AuthMiddleware).Get("/api/info", infoHandler.GetInfo)
	r.With(middleware.AuthMiddleware).Get("/api/buy/{item}", buyHandler.Buy)
	r.With(middleware.AuthMiddleware).Post("/api/sendCoin", transactionHandler.SendCoin)
	r.Post("/api/auth", userHandler.Auth)
	return r
}
