package umqr_test

import (
	"strings"
	"testing"

	"github.com/frankmwase/malawi-pay-standard/pkg/umqr"
)

func TestCRC16(t *testing.T) {
	// Standard EMVCo test vector?
	// Let's use a known string: "123456789"
	// CRC16-CCITT (False) for "123456789" is 0x29B1
	data := []byte("123456789")
	got := umqr.CalculateCRC16CCITT(data)
	want := uint16(0x29B1)
	if got != want {
		t.Errorf("CRC16 failed: got %04X, want %04X", got, want)
	}
}

func TestUMQREncoding(t *testing.T) {
	qr := umqr.GenerateMerchantQR("MUBAS Cafeteria", "Blantyre", "@mubas_cafe", "AIRTEL_MONEY", 1500.00, "LUNCH-001")

	// Basic checks
	if !strings.HasPrefix(qr, "000201") {
		t.Errorf("Missing Payload Format Indicator: %s", qr)
	}

	if !strings.Contains(qr, "5303454") {
		t.Errorf("Missing Currency MWK: %s", qr)
	}

	if !strings.Contains(qr, "MW.GOV.NATSWITCH") {
		t.Errorf("Missing Malawi Interop ID: %s", qr)
	}

	// Check if CRC tag 63 is at the end
	if !strings.Contains(qr, "6304") {
		t.Errorf("Missing CRC Tag: %s", qr)
	}

	t.Logf("Generated QR: %s", qr)
}
