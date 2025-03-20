package handlers

import (
	"fmt"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "🦋 I.M.A.G.I.N.A.T.I.O.N 🦋")
}
