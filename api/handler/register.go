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

func (h *LoyaltyHandler) Register(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		msg := fmt.Sprintf("Unsupported content type [%s]", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	cred := apiModel.CredentialsJSON{}

	dec := json.NewDecoder(r.Body)
	if errDec := dec.Decode(&cred); errDec != nil {
		msg := fmt.Sprintf("Cannot parse registration credentials: %s", errDec.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if cred.Login == "" || cred.Password == "" {
		http.Error(w, "Empty login/password not allowed", http.StatusBadRequest)
		return
	}

	errReg := h.authSrv.Register(cred.ToCredentials())
	if errors.Is(errReg, auth.ErrLoginAlreadyTaken) {
		http.Error(w, "Login already taken", http.StatusConflict)
		return
	}
	if errReg != nil {
		log.Printf("Cannot register user [%v]: %s", cred, errReg.Error())
		http.Error(w, "Cannot register user because of error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/api/user/login", http.StatusTemporaryRedirect)
}
