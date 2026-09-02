package inspection

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// monetizeConfig mirrors the API's MonetizationConfig, the payment
// configuration a Monetize action reads its pricing from.
type monetizeConfig struct {
	CryptoConfig *monetizeCryptoConfig `json:"CryptoConfig,omitempty"`
	CurrencyMode string                `json:"CurrencyMode,omitempty"`
}

type monetizeCryptoConfig struct {
	PaymentNetworks []monetizePaymentNetwork `json:"PaymentNetworks,omitempty"`
}

type monetizePaymentNetwork struct {
	Chain         string          `json:"Chain,omitempty"`
	WalletAddress string          `json:"WalletAddress,omitempty"`
	Prices        []monetizePrice `json:"Prices,omitempty"`
}

type monetizePrice struct {
	Amount   string `json:"Amount,omitempty"`
	Currency string `json:"Currency,omitempty"`
}

// monetizeResponse builds the HTTP 402 response of a matching Monetize
// rule. The action documentation specifies the status and that the
// client uses the returned pricing information to complete payment; it
// does not specify a manifest format, so this platform serves the
// manifest as JSON derived from the web ACL's monetization
// configuration and the action's price multiplier — the effective price
// per request is each network's price multiplied by the multiplier.
func (ctx *evalCtx) monetizeResponse(action *wafstore.Action) *ResolvedResponse {
	multiplier := int64(1)
	if action != nil && action.Monetize != nil {
		if parsed, err := strconv.ParseInt(action.Monetize.PriceMultiplier, 10, 64); err == nil && parsed >= 1 && parsed <= 100 {
			multiplier = parsed
		}
	}
	return &ResolvedResponse{
		StatusCode: 402,
		Headers:    []wafstore.CustomHTTPHeader{{Name: "Content-Type", Value: "application/json"}},
		Body:       ctx.priceManifest(multiplier),
	}
}

// priceManifest serialises the 402 price manifest. A web ACL without a
// parsable monetization configuration still serves 402 with the
// multiplier and an empty network list, so the terminating action is
// never silently dropped.
func (ctx *evalCtx) priceManifest(multiplier int64) string {
	manifest := map[string]interface{}{
		"priceMultiplier": strconv.FormatInt(multiplier, 10),
		"paymentNetworks": []interface{}{},
	}
	config, ok := normaliseThroughJSON[*monetizeConfig](ctx.acl.MonetizationConfig)
	if !ok || config == nil {
		encoded, _ := json.Marshal(manifest)
		return string(encoded)
	}
	if config.CurrencyMode != "" {
		manifest["currencyMode"] = config.CurrencyMode
	}
	networks := []interface{}{}
	if config.CryptoConfig != nil {
		for _, network := range config.CryptoConfig.PaymentNetworks {
			entry := map[string]interface{}{
				"chain":         network.Chain,
				"walletAddress": network.WalletAddress,
			}
			prices := make([]interface{}, 0, len(network.Prices))
			for _, price := range network.Prices {
				prices = append(prices, map[string]interface{}{
					"amount":          price.Amount,
					"effectiveAmount": effectivePrice(price.Amount, multiplier),
					"currency":        price.Currency,
				})
			}
			entry["prices"] = prices
			networks = append(networks, entry)
		}
	}
	manifest["paymentNetworks"] = networks
	encoded, _ := json.Marshal(manifest)
	return string(encoded)
}

// effectivePrice multiplies a decimal price string by an integer
// multiplier exactly. Prices carry at most three decimal places, so the
// product is computed in milli-units and formatted back with trailing
// zeros trimmed; an unparsable amount passes through unchanged.
func effectivePrice(amount string, multiplier int64) string {
	millis, err := ParsePriceMillis(amount)
	if err != nil {
		return amount
	}
	return millisToPrice(millis * multiplier)
}

// priceAmountPattern mirrors the Smithy pattern trait on PriceAmount: a
// decimal with an optional one-to-three digit fraction, no leading zeros,
// no trailing dot, no sign, and an all-zero fraction only through a
// non-zero tail.
var priceAmountPattern = regexp.MustCompile(`^([1-9][0-9]*(\.[0-9]{1,3})?|0\.([1-9][0-9]{0,2}|0[1-9][0-9]?|00[1-9]))$`)

// ParsePriceMillis converts a decimal price string with at most three
// fractional digits to its milli-unit integer form. It is the single
// price parser shared by the manifest builder and the configuration
// validator, so it enforces the PriceAmount grammar as well as the
// numeric range checked by its callers.
func ParsePriceMillis(amount string) (int64, error) {
	if !priceAmountPattern.MatchString(amount) {
		return 0, fmt.Errorf("price %q is not in the PriceAmount format", amount)
	}
	whole, frac, _ := strings.Cut(amount, ".")
	for len(frac) < 3 {
		frac += "0"
	}
	return strconv.ParseInt(whole+frac, 10, 64)
}

func millisToPrice(millis int64) string {
	whole, frac := millis/1000, millis%1000
	if frac == 0 {
		return strconv.FormatInt(whole, 10)
	}
	fracText := strings.TrimRight(fmt.Sprintf("%03d", frac), "0")
	return strconv.FormatInt(whole, 10) + "." + fracText
}
