package isumm

import (
	"fmt"
	"net/http"
)

func Inv(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Foo")
}
