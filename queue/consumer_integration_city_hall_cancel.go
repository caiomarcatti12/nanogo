package queue

import "log"

func ConsumerIntegrationCityHallCancel(body map[string]interface{}, headers map[string]interface{}) {
	log.Printf("Headers: %v", headers)
	log.Printf("Body: %s", body)
}
