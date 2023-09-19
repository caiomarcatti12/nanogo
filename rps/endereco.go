package rps

type Endereco struct {
	Endereco        string `bson:"endereco" json:"endereco"`
	Numero          string `bson:"numero" json:"numero"`
	Complemento     string `bson:"complemento" json:"complemento"`
	Bairro          string `bson:"bairro" json:"bairro"`
	CodigoMunicipio int    `bson:"codigoMunicipio" json:"codigoMunicipio"`
	UF              string `bson:"uf" json:"uf"`
	CEP             int    `bson:"cep" json:"cep"`
}
