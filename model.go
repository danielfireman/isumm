package isumm

import (
	"time"

	"appengine"
	"appengine/datastore"
)

const (
	opEntityName         = "Operation"
	investmentEntityName = "Investment"
)

type Operation struct {
	Value       float32    `datastore:"value,noindex"`
	Date        time.Time  `datastore:"date`
	Investiment Investment `datastore:"investment"`
}

type Investment struct {
	Key  string `datastore:"-"`
	Name string `datastore:"name"`
}

func investmentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, investmentEntityName, "singleton", 0, nil)
}

func PutInvestment(c appengine.Context, i *Investment) error {
	var k *datastore.Key
	if i.Key == "" {
		k = datastore.NewIncompleteKey(c, investmentEntityName, investmentKey(c))
	} else {
		decoded, err := datastore.DecodeKey(i.Key)
		if err != nil {
			return err
		}
		k = decoded
	}
	_, err := datastore.Put(c, k, i)
	return err
}

func GetInvestments(c appengine.Context) ([]*Investment, error) {
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Greeting that had just been written would not
	// show up in a query.
	q := datastore.NewQuery(investmentEntityName).Ancestor(investmentKey(c)).Order("name")
	var i []*Investment
	keys, err := q.GetAll(c, &i)
	for pos, k := range keys {
		i[pos].Key = k.Encode()
	}
	return i, err
}
