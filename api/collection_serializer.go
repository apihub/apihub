package api

type CollectionSerializer struct {
	Items interface{} `json:"items"`
	Count int         `json:"item_count"`
}

func Collection(data interface{}, count int) *CollectionSerializer {
	cs := &CollectionSerializer{
		Items: data,
		Count: count,
	}

	return cs
}
