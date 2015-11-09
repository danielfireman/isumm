package isumm

import (
	"sort"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	now := time.Now()
	o := Operations{{Date: now}, {Date: now.Add(incMin)}}
	sort.Sort(o)
	if o[0].Date.Before(o[1].Date) {
		t.Errorf("want %s before than %s", o[0], o[1])
	}
}
