package types

type Course struct {
	Name            string  `json:"name,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Picture         string  `json:"picture,omitempty"`
	PreviewUrlVideo string  `json:"preview-url-video,omitempty"`
}
