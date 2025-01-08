package rest

import (
	"encoding/json"
	"net/http"
	"valerii/crudbananas/internal/domain"
	"io"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request){
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