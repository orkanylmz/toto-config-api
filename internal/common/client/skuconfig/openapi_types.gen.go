// Package skuconfig provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package skuconfig

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// SKUResponse defines model for SKUResponse.
type SKUResponse struct {
	MainSku string `json:"main_sku"`
}

// GetSKUParams defines parameters for GetSKU.
type GetSKUParams struct {
	Package string `form:"package" json:"package"`
}
