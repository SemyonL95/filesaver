package handlers

type API struct {
	FileStorage FileStorage
	Cache       Cache
}

func NewAPI(storage FileStorage, cache Cache) *API {
	return &API{
		FileStorage: storage,
		Cache:       cache,
	}
}
