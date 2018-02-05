package metadata

import "manageiq-exchange/models/utils"

type Metadata struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

func (m *Metadata) Init(data map[string]interface{}) {
	utils.CreateFromMap(data, &m)
}
