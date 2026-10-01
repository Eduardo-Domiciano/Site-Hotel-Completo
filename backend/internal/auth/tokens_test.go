package auth_test

import (
	"testing"
	"time"

	"github.com/Eduardo-Domiciano/sitehotel-api/internal/auth"
	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := auth.HashPassword("senha-secreta")
	if err != nil {
		t.Fatal(err)
	}
	if !auth.CheckPassword(hash, "senha-secreta") {
		t.Fatal("expected password to match")
	}
	if auth.CheckPassword(hash, "outra") {
		t.Fatal("expected password mismatch")
	}
}

func TestIssueAndParseAccess(t *testing.T) {
	tokens := auth.NewTokens("test-secret-key", 15*time.Minute, 24*time.Hour)
	id := uuid.New()
	pair, hash, exp, err := tokens.Issue(id, "a@b.com")
	if err != nil {
		t.Fatal(err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" || hash == "" || exp.IsZero() {
		t.Fatal("incomplete token pair")
	}
	claims, err := tokens.ParseAccess(pair.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != id || claims.Email != "a@b.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestHashTokenStable(t *testing.T) {
	a := auth.HashToken("abc")
	b := auth.HashToken("abc")
	if a != b || a == "" {
		t.Fatal("hash should be stable and non-empty")
	}
}
