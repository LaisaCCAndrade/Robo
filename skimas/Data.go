package skimas

type Data struct {
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Country      string `json:"country"`
	Organization string `json:"organization"`
	CNPJ         string `json:"cnpj"`
}
