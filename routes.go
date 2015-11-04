package isumm

import "net/http"

func init() {
	http.HandleFunc("/app", App)
	http.HandleFunc("/inv", Inv)
	http.HandleFunc("/op", Op)
	http.HandleFunc("/delop", DelOp)
}
