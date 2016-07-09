package isumm

import "net/http"

func init() {
	http.HandleFunc("/app", App)
	http.HandleFunc("/app/ops", AppOps)
	http.HandleFunc("/inv", Inv)
	http.HandleFunc("/op", Op)
}
