package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apiModel"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/auth"
	"log"
	"net/http"
)

func (h *LoyaltyHandler) Login(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		msg := fmt.Sprintf("Unsupported content type [%s]", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	cred := apiModel.CredentialsJSON{}

	dec := json.NewDecoder(r.Body)
	if errDec := dec.Decode(&cred); errDec != nil {
		msg := fmt.Sprintf("Cannot parse credentials: %s", errDec.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	u, errLogin := h.authSrv.Login(cred.ToCredentials())
	if errors.Is(errLogin, auth.ErrWrongCredentials) {
		http.Error(w, "Incorrect login/password", http.StatusUnauthorized)
		return
	}
	if errLogin != nil {
		log.Printf("Cannot log user in [%s]: %s", cred.Login, errLogin.Error())
		http.Error(w, "Cannot log in because of server error", http.StatusInternalServerError)
		return
	}

	authCookie := makeAuthCookie(u)
	http.SetCookie(w, &authCookie)

	w.WriteHeader(http.StatusOK)
}
