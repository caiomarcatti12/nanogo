package rps

type Valores struct {
	ValorServicos          float64       `bson:"valorServicos" json:"valorServicos"`
	ValorDeducoes          float64       `bson:"valorDeducoes,omitempty" json:"valorDeducoes,omitempty"`
	ValorPIS               float64       `bson:"valorPIS,omitempty" json:"valorPIS,omitempty"`
	ValorCOFINS            float64       `bson:"valorCOFINS,omitempty" json:"valorCOFINS,omitempty"`
	ValorINSS              float64       `bson:"valorINSS,omitempty" json:"valorINSS,omitempty"`
	ValorIR                float64       `bson:"valorIR,omitempty" json:"valorIR,omitempty"`
	ValorCSLL              float64       `bson:"valorCSLL,omitempty" json:"valorCSLL,omitempty"`
	ISSRetido              ISSRetidoType `bson:"issRetido,omitempty" json:"issRetido,omitempty"`
	ValorISS               float64       `bson:"valorISS,omitempty" json:"valorISS,omitempty"`
	OutrasRetencoes        float64       `bson:"outrasRetencoes,omitempty" json:"outrasRetencoes,omitempty"`
	BaseCalculo            float64       `bson:"baseCalculo,omitempty" json:"baseCalculo,omitempty"`
	Aliquota               float64       `bson:"aliquota,omitempty" json:"aliquota,omitempty"`
	ValorLiquidoNfse       float64       `bson:"valorLiquidoNfse,omitempty" json:"valorLiquidoNfse,omitempty"`
	DescontoIncondicionado float64       `bson:"descontoIncondicionado,omitempty" json:"descontoIncondicionado,omitempty"`
	DescontoCondicionado   float64       `bson:"descontoCondicionado,omitempty" json:"descontoCondicionado,omitempty"`
}
