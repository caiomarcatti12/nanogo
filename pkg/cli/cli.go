package cli

import (
	"errors"
	flaggo "flag"
	"fmt"
	"os"

	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
)

// Flag define a estrutura para uma flag
type Flag struct {
	Key          string
	DefaultValue interface{}
	Value        interface{}
	Func         func(args ...Flags)
}

// Flags é um mapa de flags
type Flags map[string]Flag

// ICLI define a interface para o CLI
type ICLI interface {
	Initialize()
	Add(argument string, defaultValue interface{}, description string, function ...func(args ...Flags))
	Execute(args []string) error
	Get(key string, defaultValue ...interface{}) (interface{}, error)
	GetAll() (map[string]interface{}, error)
}

// CLI implementa a interface ICLI
type CLI struct {
	log   log.ILog
	flags Flags
}

// NewCLI cria uma nova instância de CLI
func Factory(log log.ILog) ICLI {
	return &CLI{
		log:   log,
		flags: make(Flags),
	}
}

// Add adiciona um novo argumento ao CLI
func (c *CLI) Add(argument string, defaultValue interface{}, description string, function ...func(args ...Flags)) {
	flag := Flag{
		Key:          argument,
		DefaultValue: defaultValue,
	}

	// flag.Value = flaggo.(argument, string(defaultValue), description)

	// if len(function) > 0 {
	// 	flag.Func = function[0]
	// }

	c.flags[argument] = flag
}

// Initialize inicializa as flags
func (c *CLI) Initialize() {

	flagCommandLine := flaggo.NewFlagSet("cli", flaggo.ContinueOnError)
	args := flagCommandLine.Parse(os.Args)

	flaggo.VisitAll(func(f *flaggo.Flag) {
		fmt.Println(f.Name)
	})
	fmt.Printf("Args: %v\n", args)

	for key, f := range c.flags {
		flagValue := flaggo.Lookup(key)
		if flagValue != nil {
			c.flags[key] = Flag{
				Key:          f.Key,
				DefaultValue: f.DefaultValue,
				Value:        flagValue.Value,
				Func:         f.Func,
			}
		}
	}
}

// Get retorna o valor de uma flag
func (c *CLI) Get(key string, defaultValue ...interface{}) (interface{}, error) {
	if flag, exists := c.flags[key]; exists {
		return flag.Value, nil
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return nil, errors.New("flag não encontrada")
}

// Execute executa a função associada à flag se existir
func (c *CLI) Execute(args []string) error {
	for _, arg := range args {
		if flag, exists := c.flags[arg]; exists && flag.Func != nil {
			// Passa as flags como argumento para a função
			flag.Func(c.flags)
		}
	}
	return nil
}

// GetAll retorna todos os valores das flags
func (c *CLI) GetAll() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, flag := range c.flags {
		result[key] = flag.Value
	}
	return result, nil
}
