package rest

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Достаем токен их хедера
		token, err := getTokenFromRequest(r)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//Парсим токен (достаем айди)
		userId, err := h.userService.ParseToken(r.Context(), token)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//добавляем айди в контекст, чтобы передать вниз по цепочке
		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		//Вызываем следующий метод
		next.ServeHTTP(w, r)
	})	
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	//разделяем хедер на 2 части (Bearer <token>)
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
