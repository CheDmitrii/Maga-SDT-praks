package graph

import "Prak_12/internal/store"

// Resolver is the root GraphQL resolver.
type Resolver struct {
	Store *store.Store
}
