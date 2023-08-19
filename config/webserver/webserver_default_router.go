package webserver


func WebserverDefaultRouter() {
	AddRouter("GET", "/healthcheck", HealthcheckHandler)
}
