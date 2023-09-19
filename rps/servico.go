package rps

type Servico struct {
	Valores                   Valores `bson:"valores" json:"valores"`
	ItemListaServico          int     `bson:"itemListaServico" json:"itemListaServico"`
	CodigoTributacaoMunicipio int     `bson:"codigoTributacaoMunicipio" json:"codigoTributacaoMunicipio"`
	Discriminacao             string  `bson:"discriminacao" json:"discriminacao"`
	CodigoMunicipio           int     `bson:"codigoMunicipio" json:"codigoMunicipio"`
}
