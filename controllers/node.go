package controllers

import (
	"github.com/audi70r/dvpn-openwrt/services/node"
	"fmt"
	"net/http"
)

func StartNodeStreamStd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	resp, err := node.StartNodeStd(w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(resp)
}

func StartNode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START NODE")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	resp, err := node.StartNode()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(resp)
}

func GetNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	resp, err := node.GetNode()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
