package isumm

import "time"

type Operation struct {
	Value float32   `datastore:"value,noindex"`
	Date  time.Time `datastore:"date,noindex"`
}

type Investiment struct {
	Name string `datastore:"name"`
	Op   []Operation
}
