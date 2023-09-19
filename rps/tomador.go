package rps

type Tomador struct {
	CNPJ               string    `bson:"cnpj,omitempty" json:"cnpj,omitempty"`
	CPF                string    `bson:"cpf,omitempty" json:"cpf,omitempty"`
	InscricaoMunicipal string    `bson:"inscricaoMunicipal,omitempty" json:"inscricaoMunicipal,omitempty"`
	RazaoSocial        string    `bson:"razaoSocial" json:"razaoSocial"`
	Endereco           *Endereco `bson:"endereco,omitempty" json:"endereco,omitempty"`
}
