package test

import (
	"replite_web/internal/app/dao"
	"testing"
)

func TestSchemaInit(t *testing.T) {
	dao.InitMogoSchema()
}
