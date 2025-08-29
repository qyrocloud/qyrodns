package health

type CheckResponse struct {
	Status string `json:"status"`
}

type ServiceInfoResponse struct {
	Service string `json:"service"`
}
