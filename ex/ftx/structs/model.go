package structs

type Response struct {
	Success bool `json:"success"`
	Result  any  `json:"result"`
}
