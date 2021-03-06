package finance

import (
	"github.com/piquette/finance-go/form"
)

// Query is the function used to get a page listing.
type Query func(*form.Values) ([]interface{}, error)

// Iter provides a convenient interface
// for iterating over the elements
// returned from paginated list API calls.
// Successive calls to the Next method
// will step through each item in the list,
// fetching pages of items as needed.
// Iterators are not thread-safe, so they should not be consumed
// across multiple goroutines.
type Iter struct {
	cur    interface{}
	err    error
	qs     *form.Values
	query  Query
	values []interface{}
}

// GetIter returns a new Iter for a given query and its options.
func GetIter(qs *form.Values, query Query) *Iter {
	iter := &Iter{}
	iter.query = query

	q := qs
	if q == nil {
		q = &form.Values{}
	}
	iter.qs = q

	iter.getPage()
	return iter
}

// GetErrIter returns a iter wrapping an error.
func GetErrIter(e error) *Iter {
	iter := &Iter{}
	iter.err = e
	return iter
}

func (it *Iter) getPage() {
	it.values, it.err = it.query(it.qs)
}

// Next advances the Iter to the next item in the list,
// which will then be available
// through the Current method.
// It returns false when the iterator stops
// at the end of the list.
func (it *Iter) Next() bool {

	if it.values == nil {
		return false
	}
	if len(it.values) == 0 {
		it.getPage()
	}
	if len(it.values) == 0 {
		return false
	}
	it.cur = it.values[0]
	it.values = it.values[1:]
	return true
}

// Current returns the most recent item
// visited by a call to Next.
func (it *Iter) Current() interface{} {
	return it.cur
}

// Err returns the error, if any,
// that caused the Iter to stop.
// It must be inspected
// after Next returns false.
func (it *Iter) Err() error {
	return it.err
}

// Count returns the list count.
func (it *Iter) Count() int {
	return len(it.values)
}
