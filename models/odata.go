package models

type ODataResponse struct {
	Context string      `json:"@odata.context,omitempty"`
	Count   int         `json:"@odata.count,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}
