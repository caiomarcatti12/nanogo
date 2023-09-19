package rps

type OrgaoGerador struct {
	CodigoMunicipio int    `bson:"codigoMunicipio" json:"codigoMunicipio"`
	Uf              string `bson:"uf" json:"uf"`
}
