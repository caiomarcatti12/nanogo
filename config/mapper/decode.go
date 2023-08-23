package mapper

import (
	"github.com/mitchellh/mapstructure"
)

func Transform(input interface{}, output any) error {
	return mapstructure.Decode(input, &output)
}
