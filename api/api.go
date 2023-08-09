package api

import "vega-cli-mm/store"

type Api struct {
	store *store.Store
}

func NewApi(
	store *store.Store,
) *Api {
	return &Api{
		store: store,
	}
}

func (a *Api) Start() {
	go func() {
		// TODO - start the API
	}()
}
