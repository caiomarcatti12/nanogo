package webserver_types

// Route define a estrutura para rotas no servidor.
type Route struct {
	Path    string
	Method  string
	IHandler interface{}
	HandlerFunc string
}
