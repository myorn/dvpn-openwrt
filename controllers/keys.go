package controllers

import (
	"encoding/json"
	"github.com/audi70r/dvpn-openwrt/services/keys"
	"net/http"
)

func ListKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	keys, err := keys.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(keys)
}
