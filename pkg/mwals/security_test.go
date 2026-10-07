package mwals_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwals"
)

func TestSecureRegistration(t *testing.T) {
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "store.json")
	svc, err := mwals.NewService(key, path)
	if err != nil {
		t.Fatal(err)
	}
	record := &mwals.AliasRecord{
		Alias: "@Alice", Status: mwals.AliasStatusActive,
		VerificationProof: "secret-challenge",
		Endpoints:         []mwals.Endpoint{{Provider: "TNM", Destination: "265881234567"}},
	}
	if err := svc.Register(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	record.Endpoints[0].Destination = "tampered"
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("permissions: %v", info.Mode())
	}
	stored, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(stored), "secret-challenge") {
		t.Fatal("verification proof persisted")
	}
	resp, err := svc.Resolve(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Endpoints[0].Destination != "265881234567" {
		t.Fatal("caller mutated stored record")
	}
	sig, err := hex.DecodeString(resp.SecuritySig)
	if err != nil {
		t.Fatal(err)
	}
	resp.SecuritySig = ""
	data, _ := json.Marshal(resp)
	if !ed25519.Verify(pub, data, sig) {
		t.Fatal("invalid response signature")
	}
	resp.Endpoints[0].Destination = "tampered"
	data, _ = json.Marshal(resp)
	if ed25519.Verify(pub, data, sig) {
		t.Fatal("endpoint excluded from signature")
	}
	if err := svc.Register(context.Background(), &mwals.AliasRecord{Alias: "../oops"}); err == nil {
		t.Fatal("invalid alias accepted")
	}
	private := &mwals.AliasRecord{
		Alias: "bob", IsPrivate: true,
		Endpoints: []mwals.Endpoint{{Provider: "TNM", Destination: "123"}},
	}
	if err := svc.Register(context.Background(), private); err == nil {
		t.Fatal("unredeemable token accepted")
	}
}

func TestRegisterAuthAndRollback(t *testing.T) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	svc, err := mwals.NewService(key, "")
	if err != nil {
		t.Fatal(err)
	}
	body := `{"alias":"alice","endpoints":[{"provider":"TNM","destination":"123"}]}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	res := httptest.NewRecorder()
	mwals.NewHandler(svc).Register(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated registration: %d", res.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer test-token")
	res = httptest.NewRecorder()
	mwals.NewAuthenticatedHandler(svc, "test-token").Register(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("authorized registration: %d: %s", res.Code, res.Body.String())
	}
	bad, err := mwals.NewService(key, filepath.Join(t.TempDir(), "missing", "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	record := &mwals.AliasRecord{Alias: "alice", Endpoints: []mwals.Endpoint{{Provider: "TNM", Destination: "123"}}}
	if err := bad.Register(context.Background(), record); err == nil {
		t.Fatal("save error ignored")
	}
	if _, err := bad.Resolve(context.Background(), "alice"); err == nil {
		t.Fatal("failed registration visible")
	}
}

func TestConcurrentResolveAndAttestation(t *testing.T) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	svc, err := mwals.NewService(key, "")
	if err != nil {
		t.Fatal(err)
	}
	svc.Seed(&mwals.AliasRecord{Alias: "alice", Status: mwals.AliasStatusActive, VerificationProof: "123456"})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = svc.Resolve(context.Background(), "alice")
			}
		}()
	}
	if err := svc.AttestAlias("alice", mwals.AttestationVerified, "123456"); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}
