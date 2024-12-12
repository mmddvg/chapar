package postgres_test

import (
	"mmddvg/chapar/pkg/repositories/postgres"
	"testing"
)

func TestIdGen(t *testing.T) {
	postgres.Initialize(1)
	t.Log("id : ", postgres.GenerateId())
}
