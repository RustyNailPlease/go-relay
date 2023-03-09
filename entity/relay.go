package entity

import "encoding/json"

//	{
//		"name": <string identifying relay>,
//		"description": <string with detailed information>,
//		"pubkey": <administrative contact pubkey>,
//		"contact": <administrative alternate contact>,
//		"supported_nips": <a list of NIP numbers supported by the relay>,
//		"software": <string identifying relay software URL>,
//		"version": <string version identifier>
//	 }

type Nips []int

type RelayMeta struct {
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Pubkey        string          `json:"pubkey"`
	Contact       string          `json:"contact"`
	SupportedNips json.RawMessage `json:"supported_nips" gorm:"type:jsonb"`
	Software      string          `json:"software"`
	Version       string          `json:"version"`
}
