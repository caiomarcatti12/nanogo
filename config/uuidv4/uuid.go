package uuid

import (
	"encoding/json"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID // Embutindo diretamente, sem criar um campo aninhado explícito.
}

// Nil representa um UUID nulo na nossa estrutura customizada.
var Nil UUID = UUID{uuid.Nil}

func New() UUID {
	return UUID{uuid.New()}
}

func NewRandom() (UUID, error) {
	id, err := uuid.NewRandom()
	return UUID{id}, err
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	u.UUID = id
	return nil
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID{id}, err
}
