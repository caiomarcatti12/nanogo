package rps

type IdentificacaoRPS struct {
	Numero int         `bson:"numero" json:"numero"`
	Serie  string      `bson:"serie" json:"serie"`
	Tipo   TipoRPSType `bson:"tipo" json:"tipo"`
}
