package main

import (
	"fmt"
	"github.com/caiomarcatti12/nanogo/pkg/mapper"
)

type Policy struct {
	id   string
	Name string
	Rule Rule
}

type Rule struct {
	id    string
	Name  string
	Value int
}

func main() {
	policy := Policy{
		id:   "abc-123",
		Name: "Teste",
		Rule: Rule{
			id:    "rule-456",
			Name:  "Regra Teste",
			Value: 42,
		},
	}

	// Serialize simples
	serialized := mapper.Serialize(policy)
	fmt.Println("Serialized struct:", serialized)

	// Serialize com slice
	policies := []Policy{
		{
			id:   "id-1",
			Name: "Policy 1",
			Rule: Rule{id: "rule-1", Name: "Regra 1", Value: 1},
		},
		{
			id:   "id-2",
			Name: "Policy 2",
			Rule: Rule{id: "rule-2", Name: "Regra 2", Value: 2},
		},
	}
	serializedSlice := mapper.Serialize(policies)
	fmt.Println("Serialized slice:", serializedSlice)

	// Deserialize
	data := map[string]interface{}{
		"id":   "novo-789",
		"Name": "Nova Policy",
		"Rule": map[string]interface{}{
			"id":    "rule-789",
			"Name":  "Nova Regra",
			"Value": 99,
		},
	}
	var newPolicy Policy
	mapper.Deserialize(data, &newPolicy)
	fmt.Printf("Deserialized struct: %+v\n", newPolicy)
}
