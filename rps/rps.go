package rps

import "time"

type RPS struct {
	IdentificacaoRPS         IdentificacaoRPS             `bson:"identificacaoRPS" json:"identificacaoRPS,omitempty"`
	DataEmissao              time.Time                    `bson:"dataEmissao" json:"dataEmissao,omitempty"`
	NaturezaOperacao         NaturezaOperacaoType         `bson:"naturezaOperacao" json:"naturezaOperacao,omitempty"`
	RegimeEspecialTributacao RegimeEspecialTributacaoType `bson:"regimeEspecialTributacao" json:"regimeEspecialTributacao,omitempty"`
	OptanteSimplesNacional   OptanteSimplesNacionalType   `bson:"optanteSimplesNacional" json:"optanteSimplesNacional,omitempty"`
	IncentivadorCultural     IncentivadorCulturalType     `bson:"incentivadorCultural" json:"incentivadorCultural,omitempty"`
	Status                   StatusType                   `bson:"status" json:"status,omitempty"`
	Tomador                  Tomador                      `bson:"tomador" json:"tomadorServico,omitempty"`
	Servico                  Servico                      `bson:"servico" json:"servico,omitempty"`
	IntermediarioServico     *IntermediarioServico        `bson:"intermediarioServico,omitempty" json:"intermediarioServico,omitempty"`
	ConstrucaoCivil          *ConstrucaoCivil             `bson:"construcaoCivil,omitempty" json:"construcaoCivil,omitempty"`
	OrgaoGerador             *OrgaoGerador                `bson:"orgaoGerador,omitempty" json:"orgaoGerador,omitempty"`
	RpsSubstituido           *IdentificacaoRPS            `bson:"rpsSubstituido,omitempty" json:"rpsSubstituido,omitempty"`
}
