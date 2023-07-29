package test

import (
	"log"
	"replite_web/internal/app/service"
	"testing"
)

func TestQueryRuleAuthorization(t *testing.T) {
	role := "admin"
	response := service.QueryRuleAuthorization(role)
	log.Println(response)
}
