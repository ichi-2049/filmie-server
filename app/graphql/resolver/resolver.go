package resolver

import "github.com/ichi-2049/filmie-server/graphql/resolver/container"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	container *container.Container
}

func NewResolver(container *container.Container) *Resolver {
	return &Resolver{
		container: container,
	}
}
