package umqr

import "fmt"

// GenerateMerchantQR creates a standard UMQR string for a merchant.
func GenerateMerchantQR(merchantName, city, alias, provider string, amount float64, reference string) string {
	enc := NewEncoder()
	enc.Set(TagPayloadFormatIndicator, "01")
	enc.Set(TagPointOfInitiationMethod, "11") // Static Sticker

	// Malawi Specific Interop Data
	tag26Value := EncodeTag26(provider, alias)
	enc.Set(TagMalawiMerchantAccount, tag26Value)

	enc.Set(TagMerchantCategoryCode, "5411") // Default to Grocery Stores
	enc.Set(TagTransactionCurrency, "454")   // MWK

	if amount > 0 {
		enc.Set(TagTransactionAmount, fmt.Sprintf("%.2f", amount))
	}

	enc.Set(TagCountryCode, "MW")
	enc.Set(TagMerchantName, merchantName)
	enc.Set(TagMerchantCity, city)

	if reference != "" {
		enc.Set(TagAdditionalData, TLV{Tag: "01", Value: reference}.String())
	}

	return enc.Encode()
}
