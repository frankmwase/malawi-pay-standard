package mwussd

import (
	"testing"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
)

func TestGenerateAirtelSession(t *testing.T) {
	router := NewRouter()
	txn := &mwjson.Transaction{
		Payload: mwjson.Payload{
			Amount: 5000,
			Sender: mwjson.Participant{
				Provider: mwjson.ProviderAirtelMoney,
			},
			Receiver: mwjson.Participant{
				ID: "0999123456",
			},
		},
	}

	steps, err := router.GenerateSession(txn, "1234")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(steps) != 5 {
		t.Errorf("Expected 5 steps for Airtel Money flow, got %d", len(steps))
	}

	if steps[0].Action != ActionDial || steps[0].Content != "*121#" {
		t.Errorf("Step 1 failed: Expected DIAL *121#, got %v %v", steps[0].Action, steps[0].Content)
	}

	if steps[4].Action != ActionReply || steps[4].Content != "1234" {
		t.Errorf("Step 5 failed: Expected PIN entry, got %v %v", steps[4].Action, steps[4].Content)
	}
}

func TestGenerateTNMSession(t *testing.T) {
	router := NewRouter()
	txn := &mwjson.Transaction{
		Payload: mwjson.Payload{
			Amount: 2500,
			Sender: mwjson.Participant{
				Provider: mwjson.ProviderTNMPamba,
			},
			Receiver: mwjson.Participant{
				ID: "0888123456",
			},
		},
	}

	steps, err := router.GenerateSession(txn, "9999")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(steps) != 5 {
		t.Errorf("Expected 5 steps for TNM Mpamba flow, got %d", len(steps))
	}

	if steps[0].Action != ActionDial || steps[0].Content != "*444#" {
		t.Errorf("Step 1 failed: Expected DIAL *444#, got %v %v", steps[0].Action, steps[0].Content)
	}
}
