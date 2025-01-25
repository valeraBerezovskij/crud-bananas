package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"valerii/crudbananas/internal/domain"
)

func (h *Handler) routeAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.SignUp(w, r)
	case http.MethodGet:
		h.SignIn(w, r)
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	//Читаем тело запроса
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Преобразуем тело запроса из JSON в структуру domain.SignUpInput
	var inp domain.SignUpInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверяем данные структуры на валидность
	if err := inp.Validate(); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Отправляем на уровень сервис
	err = h.userService.SignUp(r.Context(), inp)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	//Читаем тело запроса
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Преобразуем тело запроса из JSON в структуру domain.SignInInput
	var inp domain.SignInInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверяем данные структуры на валидность
	if err := inp.Validate(); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Отправляем на уровень сервиса
	token, err := h.userService.SignIn(r.Context(), inp)
	if err != nil{
		if errors.Is(err, domain.ErrUserNotFound){
			handleNotFoundError(w, err)
			return
		}

		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//записываем token в ответ
	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil{
		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Отправляем ответ на клиент
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func handleNotFoundError(w http.ResponseWriter, err error){
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}