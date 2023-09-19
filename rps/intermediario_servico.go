package rps

type IntermediarioServico struct {
	RazaoSocial        string `bson:"razaoSocial" json:"razaoSocial"`
	CNPJ               string `bson:"cnpj,omitempty" json:"cnpj,omitempty"`
	CPF                string `bson:"cpf,omitempty" json:"cpf,omitempty"`
	InscricaoMunicipal string `bson:"inscricaoMunicipal,omitempty" json:"inscricaoMunicipal,omitempty"`
}
