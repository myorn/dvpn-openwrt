package controllers

import (
	"github.com/solarlabsteam/dvpn-openwrt/services/dvpnconf"
	"io/ioutil"
	"net/http"
)

func Config(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		config, err := dvpnconf.GetConfigs()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(config)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		config, err := dvpnconf.ValidateAndUnmarshal(body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		resp, err := dvpnconf.PostConfig(config)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
