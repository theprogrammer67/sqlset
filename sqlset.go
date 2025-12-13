// Package sqlset is a way to store SQL queries separated from the go code.
// Query sets are stored in the .sql files, every filename without extension is an SQL set ID.
// Every file contains queries, marked with query IDs using special syntax,
// see `testdata/valid/*.sql` files for examples.
// Also file may contain JSON-encoded query set metadata with name and description.
package sqlset

import "fmt"

// SQLQueriesProvider is the interface for getting SQL queries.
type SQLQueriesProvider interface {
	// Get returns a query by set ID and query ID.
	// If the set or query is not found, it returns an error.
	Get(setID string, queryID string) (string, error)
	// MustGet returns a query by set ID and query ID.
	// It panics if the set or query is not found.
	MustGet(setID string, queryID string) string
}

// SQLSetsProvider is the interface for getting information about query sets.
type SQLSetsProvider interface {
	// GetAllMetas returns metadata for all registered query sets.
	GetAllMetas() []QuerySetMeta
}

// SQLSet is a container for multiple query sets, organized by set ID.
// It provides methods to access SQL queries and metadata.
// Use New to create a new instance.
type SQLSet struct {
	sets map[string]QuerySet
}

// Get retrieves a specific SQL query by its set ID and query ID.
// It returns an error if the query set or the query itself cannot be found.
func (s *SQLSet) Get(setID string, queryID string) (string, error) {
	return s.findQuery(setID, queryID)
}

// MustGet is like Get but panics if the query set or query is not found.
// This is useful for cases where the query is expected to exist and its absence is a critical error.
func (s *SQLSet) MustGet(setID string, queryID string) string {
	q, err := s.findQuery(setID, queryID)
	if err != nil {
		panic(err)
	}

	return q
}

// GetAllMetas returns a slice of metadata for all the query sets loaded.
// The order of the returned slice is not guaranteed.
func (s *SQLSet) GetAllMetas() []QuerySetMeta {
	metas := make([]QuerySetMeta, 0, len(s.sets))

	for _, qs := range s.sets {
		metas = append(metas, qs.GetMeta())
	}

	return metas
}

func (s *SQLSet) findQuery(setID string, queryID string) (string, error) {
	if s.sets == nil {
		return "", fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
	}

	qs, ok := s.sets[setID]
	if !ok {
		return "", fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
	}

	q, err := qs.findQuery(queryID)
	if err != nil {
		return "", err
	}

	return q, nil
}

func (s *SQLSet) registerQuerySet(setID string, qs QuerySet) {
	if s.sets == nil {
		s.sets = make(map[string]QuerySet)
	}

	s.sets[setID] = qs
}

// QuerySet represents a single set of queries, usually from a single .sql file.
type QuerySet struct {
	meta    QuerySetMeta
	queries map[string]string
}

// GetMeta returns the metadata associated with the query set.
func (qs *QuerySet) GetMeta() QuerySetMeta {
	return qs.meta
}

func (qs *QuerySet) registerQuery(id string, query string) {
	if qs.queries == nil {
		qs.queries = make(map[string]string)
	}

	qs.queries[id] = query
}

func (qs *QuerySet) findQuery(id string) (string, error) {
	if qs.queries == nil {
		return "", fmt.Errorf("%s: %w", id, ErrQueryNotFound)
	}

	q, ok := qs.queries[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", id, ErrQueryNotFound)
	}

	return q, nil
}

// QuerySetMeta holds the metadata for a query set.
type QuerySetMeta struct {
	// ID is the unique identifier for the set, derived from the filename.
	ID string `json:"id"`
	// Name is a human-readable name for the query set, from the metadata block.
	Name string `json:"name"`
	// Description provides more details about the query set, from the metadata block.
	Description string `json:"description,omitempty"`
}
