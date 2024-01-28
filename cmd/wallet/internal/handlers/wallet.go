package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cubawheeler.io/pkg/cubawheeler"
)

type WalletHandler struct {
	service cubawheeler.WalletService
}

func NewWalletHandler(service cubawheeler.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) error {
	client := cubawheeler.ClientFromContext(r.Context())
	if client == nil {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	owner := r.FormValue("owner")
	_, err := h.service.Create(r.Context(), owner)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *WalletHandler) Balance(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	user := cubawheeler.UserFromContext(r.Context())
	balance, err := h.service.Balance(r.Context(), user.ID)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(balance)
}

func (h *WalletHandler) Transactions(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	user := cubawheeler.UserFromContext(r.Context())
	txs, err := h.service.Transactions(r.Context(), user.ID)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(txs)
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	to := r.FormValue("to")
	if to == "" {
		return cubawheeler.NewMissingParameter("to")
	}
	strAmount := r.FormValue("amount")
	if strAmount == "" {
		return cubawheeler.NewMissingParameter("amount")
	}
	amount, err := strconv.ParseInt(strAmount, 10, 64)
	if err != nil {
		return cubawheeler.NewInvalidParameter("amount", strAmount)
	}
	currency := r.FormValue("currency")
	if currency == "" {
		return cubawheeler.NewMissingParameter("currency")
	}

	tx, err := h.service.Transfer(r.Context(), to, amount, currency)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(tx)
}

func (h *WalletHandler) ConfirmTransfer(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	id := r.FormValue("id")
	if id == "" {
		return cubawheeler.NewMissingParameter("id")
	}
	pin := r.FormValue("pin")
	if pin == "" {
		return cubawheeler.NewMissingParameter("pin")
	}
	return h.service.ConfirmTransfer(r.Context(), id, pin)
}

type topupBody struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	To       string `json:"to"`
}

func (b topupBody) Valid() error {
	if b.Amount <= 0 {
		return cubawheeler.NewInvalidParameter("amount", strconv.FormatInt(b.Amount, 10))
	}

	if b.Currency == "" {
		return cubawheeler.NewMissingParameter("currency")
	}

	if b.To == "" {
		return cubawheeler.NewMissingParameter("to")
	}

	return nil
}

func (h *WalletHandler) TopUp(w http.ResponseWriter, r *http.Request) (err error) {
	if !canDo(r, cubawheeler.RoleAdmin) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}
	body := topupBody{
		To: r.FormValue("to"),
	}
	if amount := r.FormValue("amount"); amount != "" {
		body.Amount, err = strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return cubawheeler.NewInvalidParameter("amount", amount)
		}
	}
	body.Currency = r.FormValue("currency")
	if err := body.Valid(); err != nil {
		return err
	}

	tx, err := h.service.Deposit(r.Context(), body.To, body.Amount, body.Currency)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(tx)
}