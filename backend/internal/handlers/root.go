package handlers

import (
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
	//fmt.Fprintf(w, "🦋 I.M.A.G.I.N.A.T.I.O.N 🦋")
}
