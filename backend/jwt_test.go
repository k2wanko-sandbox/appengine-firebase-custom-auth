package backend

import "testing"

func TestJWT(t *testing.T) {
	u := &User{ID: 1, Email: "test@example.com"}
	toks, err := createToken(u).SignedString(privateKey)
	if err != nil {
		t.Fatalf("createToken %v\n", err)
	}

	t.Logf("TokenString = %s", toks)

	tok, err := parseTokenString(toks)
	t.Logf("parseTokenString: Token = %#v\n", tok)
	if err != nil {
		t.Fatalf("parseTokenString: %v", err)
	}

	if !tok.Valid {
		t.Errorf("parseTokenString valid")
	}
}
