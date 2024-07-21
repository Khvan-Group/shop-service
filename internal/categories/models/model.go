// package categoryModels
package models

// swagger:model
type Category struct {
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
}
