package types

//Course is a representation object of course
type Course struct {
	Name            string  `json:"name,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Picture         string  `json:"picture,omitempty"`
	PreviewURLVideo string  `json:"preview-url-video,omitempty"`
}
