package test_util

import (
	"net/http"
	"testing"

	db "github.com/ot07/next-bazaar/db/sqlc"
)

type SeedData = map[string]interface{}

func NoopSetupAuth(request *http.Request, sessionToken string) {}

func NoopCreateSeed(t *testing.T, store db.Store) {}

func NoopCreateAndReturnSeed(t *testing.T, store db.Store) SeedData {
	return nil
}
