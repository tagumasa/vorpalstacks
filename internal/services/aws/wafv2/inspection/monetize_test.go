package inspection

import (
	"encoding/json"
	"testing"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

func monetizeACL(config interface{}, multiplier string) *wafstore.WebACL {
	action := &wafstore.Action{Monetize: &wafstore.MonetizeAction{}}
	if multiplier != "" {
		action.Monetize.PriceMultiplier = multiplier
	}
	return &wafstore.WebACL{
		Name:               "monetize-acl",
		Scope:              "CLOUDFRONT",
		DefaultAction:      allowDefaultAction(),
		MonetizationConfig: config,
		Rules: []*wafstore.Rule{
			{
				Name:             "monetize-rule",
				Priority:         1,
				Action:           action,
				Statement:        uriMatchStatement("/premium"),
				VisibilityConfig: &wafstore.VisibilityConfig{MetricName: "monetize-rule"},
			},
		},
	}
}

func sampleMonetizationConfig() map[string]interface{} {
	return map[string]interface{}{
		"CurrencyMode": "TEST",
		"CryptoConfig": map[string]interface{}{
			"PaymentNetworks": []interface{}{
				map[string]interface{}{
					"Chain":         "BASE_SEPOLIA",
					"WalletAddress": "0x1234567890123456789012345678901234567890",
					"Prices": []interface{}{
						map[string]interface{}{"Amount": "0.010", "Currency": "USDC"},
					},
				},
			},
		},
	}
}

func TestMonetizeActionServes402PriceManifest(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	result := evaluator.Evaluate(monetizeACL(sampleMonetizationConfig(), "3"), baseRequest("/premium"))

	if result.Action != ActionMonetize {
		t.Fatalf("Action = %q, want %q", result.Action, ActionMonetize)
	}
	if result.CustomResponse == nil || result.CustomResponse.StatusCode != 402 {
		t.Fatalf("monetize response = %+v, want status 402", result.CustomResponse)
	}
	var manifest struct {
		PriceMultiplier string `json:"priceMultiplier"`
		CurrencyMode    string `json:"currencyMode"`
		PaymentNetworks []struct {
			Chain         string `json:"chain"`
			WalletAddress string `json:"walletAddress"`
			Prices        []struct {
				Amount          string `json:"amount"`
				EffectiveAmount string `json:"effectiveAmount"`
				Currency        string `json:"currency"`
			} `json:"prices"`
		} `json:"paymentNetworks"`
	}
	if err := json.Unmarshal([]byte(result.CustomResponse.Body), &manifest); err != nil {
		t.Fatalf("manifest is not JSON: %v (%q)", err, result.CustomResponse.Body)
	}
	if manifest.PriceMultiplier != "3" || manifest.CurrencyMode != "TEST" {
		t.Fatalf("manifest header fields = %+v", manifest)
	}
	if len(manifest.PaymentNetworks) != 1 {
		t.Fatalf("manifest networks = %+v", manifest.PaymentNetworks)
	}
	network := manifest.PaymentNetworks[0]
	if network.Chain != "BASE_SEPOLIA" || network.WalletAddress != "0x1234567890123456789012345678901234567890" {
		t.Fatalf("manifest network = %+v", network)
	}
	if len(network.Prices) != 1 {
		t.Fatalf("manifest prices = %+v", network.Prices)
	}
	price := network.Prices[0]
	if price.Amount != "0.010" || price.EffectiveAmount != "0.03" || price.Currency != "USDC" {
		t.Fatalf("manifest price = %+v, want amount 0.010 effective 0.03 USDC", price)
	}
}

func TestMonetizeActionDefaultsToUnitMultiplier(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	result := evaluator.Evaluate(monetizeACL(sampleMonetizationConfig(), ""), baseRequest("/premium"))
	var manifest struct {
		PaymentNetworks []struct {
			Prices []struct {
				EffectiveAmount string `json:"effectiveAmount"`
			} `json:"prices"`
		} `json:"paymentNetworks"`
	}
	if err := json.Unmarshal([]byte(result.CustomResponse.Body), &manifest); err != nil {
		t.Fatalf("manifest is not JSON: %v", err)
	}
	if got := manifest.PaymentNetworks[0].Prices[0].EffectiveAmount; got != "0.01" {
		t.Fatalf("effective amount = %q, want 0.01 with the default unit multiplier", got)
	}
}

func TestMonetizeActionWithoutConfigurationStillInterrupts(t *testing.T) {
	evaluator := NewEvaluator(Resolvers{})
	result := evaluator.Evaluate(monetizeACL(nil, ""), baseRequest("/premium"))
	if result.Action != ActionMonetize || result.CustomResponse.StatusCode != 402 {
		t.Fatalf("result = %+v, want a 402 interruption even without a parsable configuration", result)
	}
}

func TestEffectivePrice(t *testing.T) {
	cases := []struct {
		amount     string
		multiplier int64
		want       string
	}{
		{"0.010", 3, "0.03"},
		{"0.001", 100, "0.1"},
		{"0.5", 2, "1"},
		{"1", 7, "7"},
		{"2.500", 4, "10"},
		{"not-a-price", 5, "not-a-price"},
	}
	for _, tc := range cases {
		if got := effectivePrice(tc.amount, tc.multiplier); got != tc.want {
			t.Errorf("effectivePrice(%q, %d) = %q, want %q", tc.amount, tc.multiplier, got, tc.want)
		}
	}
}

func TestParsePriceMillis(t *testing.T) {
	valid := []struct {
		amount string
		millis int64
	}{
		{"1", 1000},
		{"0.001", 1},
		{"0.5", 500},
		{"0.05", 50},
		{"0.005", 5},
		{"0.500", 500},
		{"10.25", 10250},
		{"999999999.999", 999999999999},
	}
	for _, tc := range valid {
		millis, err := ParsePriceMillis(tc.amount)
		if err != nil || millis != tc.millis {
			t.Errorf("ParsePriceMillis(%q) = (%d, %v), want (%d, nil)", tc.amount, millis, err, tc.millis)
		}
	}
	invalid := []string{
		"", "01.5", "1.", "+1", "-1", ".5", "0.000", "1.0000", " 1", "1 000",
	}
	for _, amount := range invalid {
		if millis, err := ParsePriceMillis(amount); err == nil {
			t.Errorf("ParsePriceMillis(%q) = %d, want an error", amount, millis)
		}
	}
}
