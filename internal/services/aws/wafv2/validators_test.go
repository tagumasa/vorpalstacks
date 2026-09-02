package wafv2

import (
	"strings"
	"testing"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

func TestValidateVisibilityConfig(t *testing.T) {
	valid := &wafstore.VisibilityConfig{MetricName: "test-metric"}
	if err := validateVisibilityConfig(valid); err != nil {
		t.Fatalf("valid config err = %v", err)
	}
	if err := validateVisibilityConfig(nil); err == nil {
		t.Fatal("nil config: expected error, got nil")
	}

	cases := []struct {
		name string
		vc   *wafstore.VisibilityConfig
	}{
		{"empty MetricName", &wafstore.VisibilityConfig{MetricName: ""}},
		{"too long MetricName", &wafstore.VisibilityConfig{MetricName: string(make([]byte, 256))}},
		{"invalid characters", &wafstore.VisibilityConfig{MetricName: "bad metric!"}},
		{"space rejected", &wafstore.VisibilityConfig{MetricName: "has space"}},
	}
	for _, tc := range cases {
		if err := validateVisibilityConfig(tc.vc); err == nil {
			t.Fatalf("%s: expected error, got nil", tc.name)
		}
	}

	// Pattern-allowed punctuation must pass.
	allowed := &wafstore.VisibilityConfig{MetricName: "My#Metric:1.0-2/x"}
	if err := validateVisibilityConfig(allowed); err != nil {
		t.Fatalf("allowed punctuation err = %v", err)
	}
}

func TestValidResourceTypes(t *testing.T) {
	for _, rt := range []string{
		"APPLICATION_LOAD_BALANCER", "API_GATEWAY", "APPSYNC",
		"COGNITO_USER_POOL", "APP_RUNNER_SERVICE", "VERIFIED_ACCESS_INSTANCE",
		"AMPLIFY", "AGENTCORE_GATEWAY",
	} {
		if !validResourceTypes[rt] {
			t.Errorf("expected %s to be valid", rt)
		}
	}
	for _, rt := range []string{"", "EC2_INSTANCE", "application_load_balancer", "S3"} {
		if validResourceTypes[rt] {
			t.Errorf("expected %s to be invalid", rt)
		}
	}
}

func TestCalculateRulesCapacity(t *testing.T) {
	// An anchored string match costs 2 WCUs; AndStatement costs 1 plus
	// the sum of nested statements.
	rules := []*wafstore.Rule{
		{
			Name: "byte-match",
			Statement: &wafstore.Statement{
				ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("x")},
			},
		},
		{
			Name: "and",
			Statement: &wafstore.Statement{
				AndStatement: &wafstore.AndStatement{
					Statements: []*wafstore.Statement{
						{ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("y")}},
						{ByteMatchStatement: &wafstore.ByteMatchStatement{SearchString: []byte("z")}},
					},
				},
			},
		},
	}
	// 2 + (1 + 2 + 2) = 7
	if got := calculateRulesCapacity(rules); got != 7 {
		t.Fatalf("capacity = %d, want 7", got)
	}
	if got := calculateRulesCapacity(nil); got != 0 {
		t.Fatalf("nil rules capacity = %d, want 0", got)
	}
}

// TestCalculateStatementCapacityModifiers pins the documented cost
// modifiers: the substring positional constraints, per-transformation
// cost (NONE free), the JSON body doubling, the AllQueryArguments
// surcharge and the rate-based custom-key cost.
func TestCalculateStatementCapacityModifiers(t *testing.T) {
	contains := &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
		PositionalConstraint: "CONTAINS",
		SearchString:         []byte("x"),
	}}
	if got := calculateStatementCapacity(contains); got != wcuByteMatchSubstring {
		t.Fatalf("CONTAINS base = %d, want %d", got, wcuByteMatchSubstring)
	}
	withLowercase := &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
		FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
		PositionalConstraint: "CONTAINS",
		SearchString:         []byte("x"),
		TextTransformations: []*wafstore.TextTransformation{
			{Priority: 0, Type: "NONE"},
			{Priority: 1, Type: "LOWERCASE"},
		},
	}}
	if got := calculateStatementCapacity(withLowercase); got != wcuByteMatchSubstring+wcuTextTransformation {
		t.Fatalf("CONTAINS with one charged transformation = %d, want %d", got, wcuByteMatchSubstring+wcuTextTransformation)
	}
	jsonBody := &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
		FieldToMatch:         &wafstore.FieldToMatch{JsonBody: &wafstore.JsonBody{}},
		PositionalConstraint: "EXACTLY",
		SearchString:         []byte("x"),
	}}
	if got := calculateStatementCapacity(jsonBody); got != wcuByteMatchAnchored*2 {
		t.Fatalf("EXACTLY on the JSON body = %d, want %d", got, wcuByteMatchAnchored*2)
	}
	allArgs := &wafstore.Statement{ByteMatchStatement: &wafstore.ByteMatchStatement{
		FieldToMatch:         &wafstore.FieldToMatch{AllQueryArguments: &wafstore.All{}},
		PositionalConstraint: "EXACTLY",
		SearchString:         []byte("x"),
	}}
	if got := calculateStatementCapacity(allArgs); got != wcuByteMatchAnchored+wcuAllQueryArguments {
		t.Fatalf("EXACTLY on all query arguments = %d, want %d", got, wcuByteMatchAnchored+wcuAllQueryArguments)
	}
	rate := &wafstore.Statement{RateBasedStatement: &wafstore.RateBasedStatement{
		Limit:            100,
		AggregateKeyType: "CUSTOM_KEYS",
		CustomKeys: []*wafstore.RateBasedStatementCustomKey{
			{Header: &wafstore.RateLimitHeaderKey{Name: "X-Tenant"}},
			{HTTPMethod: &wafstore.RateLimitEmptyKey{}},
		},
	}}
	if got := calculateStatementCapacity(rate); got != wcuRateBase+2*wcuRateCustomKey {
		t.Fatalf("rate base with two custom keys = %d, want %d", got, wcuRateBase+2*wcuRateCustomKey)
	}
}

func TestEnsureRuleGroupNotReferenced(t *testing.T) {
	// ARN-based scan helper: nil/empty ARN is a no-op.
	if err := ensureRuleGroupNotReferenced(nil, ""); err != nil {
		t.Fatalf("empty ARN err = %v, want nil", err)
	}
}

func TestEnsureNotAssociated(t *testing.T) {
	// No association stores and an empty ARN are both no-ops.
	if err := ensureNotAssociated(nil, ""); err != nil {
		t.Fatalf("empty ARN err = %v, want nil", err)
	}
	if err := ensureNotAssociated(nil, "arn:aws:wafv2::regional/webacl/x/y"); err != nil {
		t.Fatalf("nil stores err = %v, want nil", err)
	}
}

func TestValidateDefaultAction(t *testing.T) {
	if err := validateDefaultAction(nil); err == nil {
		t.Fatal("nil action: expected error, got nil")
	}
	if err := validateDefaultAction(&wafstore.Action{Allow: &wafstore.AllowAction{}}); err != nil {
		t.Fatalf("allow action err = %v", err)
	}
	if err := validateDefaultAction(&wafstore.Action{Block: &wafstore.BlockAction{}}); err != nil {
		t.Fatalf("block action err = %v", err)
	}
	if err := validateDefaultAction(&wafstore.Action{Count: &wafstore.CountAction{}}); err == nil {
		t.Fatal("count action: expected error, got nil")
	}
}

// TestValidateEntityNamePattern pins the Smithy EntityName @pattern
// ^[\w\-]+$ alongside the length gate: spaces, punctuation outside the
// class, and multibyte characters are rejected.
func TestValidateEntityNamePattern(t *testing.T) {
	for _, name := range []string{"abc-DEF_1", "ipset", "123"} {
		if err := validateEntityName(name); err != nil {
			t.Errorf("valid name %q rejected: %v", name, err)
		}
	}
	for _, name := range []string{"bad name", "name!", "bad/name", "\u65e5\u672c\u8a9e", "a.b"} {
		if err := validateEntityName(name); err == nil {
			t.Errorf("invalid name %q accepted", name)
		}
	}
}

// TestValidateEntityDescriptionPattern pins the Smithy EntityDescription
// @pattern: the value must start and end with a class character and carry
// at least one middle character (the pattern's minimum is 3 characters).
// An empty description stays accepted: the protocol layer cannot
// distinguish an omitted optional member from an explicitly empty one.
func TestValidateEntityDescriptionPattern(t *testing.T) {
	for _, desc := range []string{"valid desc", "a:b", "x-y", "WebACL for tests 2026"} {
		if err := validateEntityDescription(desc); err != nil {
			t.Errorf("valid description %q rejected: %v", desc, err)
		}
	}
	for _, desc := range []string{"a", "ab", " leading", "trailing ", "desc!", "\u65e5\u672c"} {
		if err := validateEntityDescription(desc); err == nil {
			t.Errorf("invalid description %q accepted", desc)
		}
	}
	if err := validateEntityDescription(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
	// A multibyte description over 256 bytes but under 256 characters
	// passes the length gate and is rejected by the pattern instead, so
	// the reported character count must reflect runes, not bytes.
	cjkLong := strings.Repeat("\u65e5", 100)
	err := validateEntityDescription(cjkLong)
	if err == nil {
		t.Error("100-character CJK description accepted")
	} else if !strings.Contains(err.Error(), "must start and end with an allowed character") {
		t.Errorf("100-character CJK description rejected on length, not pattern: %v", err)
	}
}

// TestValidateTokenDomainPattern pins the Smithy TokenDomain @pattern
// ^[\w./-]+$ on WebACL API-key token domains.
func TestValidateTokenDomainPattern(t *testing.T) {
	valid := []interface{}{"abc.com", "store.abc.com", "a-b.example"}
	if err := validateTokenDomains(valid); err != nil {
		t.Errorf("valid token domains rejected: %v", err)
	}
	for _, dom := range []string{"a b", "abc!", "abc\\com"} {
		invalid := []interface{}{dom}
		if err := validateTokenDomains(invalid); err == nil {
			t.Errorf("invalid token domain %q accepted", dom)
		}
	}
}

// TestValidateCustomResponseBodiesKeyPattern pins that CustomResponseBodies
// map keys follow the EntityName shape (the map key targets the same
// ^[\w\-]+$ pattern).
func TestValidateCustomResponseBodiesKeyPattern(t *testing.T) {
	if err := validateCustomResponseBodies(map[string]interface{}{"ok-key_1": nil}); err != nil {
		t.Errorf("valid key rejected: %v", err)
	}
	for _, key := range []string{"bad key", "key!", "key/name"} {
		if err := validateCustomResponseBodies(map[string]interface{}{key: nil}); err == nil {
			t.Errorf("invalid key %q accepted", key)
		}
	}
}

// TestValidateStatementRegexStringUnicodeLengths pins that the
// RegexMatchStatement RegexString bound follows the Smithy
// RegexPatternString @length(1, 512) trait counted in Unicode characters:
// the shape's pattern is ".*", so multibyte regex patterns are valid input
// and must not be rejected on byte length.
func TestValidateStatementRegexStringUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	inRange := &wafstore.Statement{RegexMatchStatement: &wafstore.RegexMatchStatement{
		RegexString: strings.Repeat(cjk, 200),
	}}
	if err := validateStatement(inRange, false); err != nil {
		t.Errorf("200-character CJK RegexString rejected: %v", err)
	}

	overRange := &wafstore.Statement{RegexMatchStatement: &wafstore.RegexMatchStatement{
		RegexString: strings.Repeat(cjk, 513),
	}}
	if err := validateStatement(overRange, false); err == nil {
		t.Error("513-character CJK RegexString accepted")
	}
}

// TestValidateRateCustomKeys pins the model requirements for CUSTOM_KEYS
// aggregation: at least one aggregation key, exactly one union member per
// entry, and no address-only key set (aggregating on only the IP or
// forwarded IP address belongs to the IP and FORWARDED_IP aggregate key
// types).
func TestValidateRateCustomKeys(t *testing.T) {
	base := func(keys []*wafstore.RateBasedStatementCustomKey) *wafstore.Statement {
		return &wafstore.Statement{RateBasedStatement: &wafstore.RateBasedStatement{
			Limit:            100,
			AggregateKeyType: "CUSTOM_KEYS",
			CustomKeys:       keys,
		}}
	}
	if err := validateStatement(base(nil), false); err == nil {
		t.Error("CUSTOM_KEYS without CustomKeys accepted")
	}
	if err := validateStatement(base([]*wafstore.RateBasedStatementCustomKey{
		{IP: &wafstore.RateLimitEmptyKey{}},
	}), false); err == nil {
		t.Error("address-only CustomKeys accepted")
	}
	if err := validateStatement(base([]*wafstore.RateBasedStatementCustomKey{
		{IP: &wafstore.RateLimitEmptyKey{}, HTTPMethod: &wafstore.RateLimitEmptyKey{}},
	}), false); err == nil {
		t.Error("CustomKeys entry with two union members accepted")
	}
	if err := validateStatement(base([]*wafstore.RateBasedStatementCustomKey{
		{IP: &wafstore.RateLimitEmptyKey{}},
		{HTTPMethod: &wafstore.RateLimitEmptyKey{}},
	}), false); err != nil {
		t.Errorf("IP combined with HTTPMethod rejected: %v", err)
	}
	if err := validateStatement(base([]*wafstore.RateBasedStatementCustomKey{
		{Header: &wafstore.RateLimitHeaderKey{Name: "X-Tenant"}},
	}), false); err != nil {
		t.Errorf("single header key rejected: %v", err)
	}
}

func TestValidateMonetizationConfig(t *testing.T) {
	// Each case builds its configuration from scratch: mutation helpers
	// sharing nested maps poison one another's fixtures.
	build := func(networks ...map[string]interface{}) map[string]interface{} {
		list := make([]interface{}, 0, len(networks))
		for _, network := range networks {
			list = append(list, network)
		}
		return map[string]interface{}{
			"CurrencyMode": "TEST",
			"CryptoConfig": map[string]interface{}{"PaymentNetworks": list},
		}
	}
	network := func(chain, wallet string, prices ...map[string]interface{}) map[string]interface{} {
		list := make([]interface{}, 0, len(prices))
		for _, price := range prices {
			list = append(list, price)
		}
		return map[string]interface{}{
			"Chain":         chain,
			"WalletAddress": wallet,
			"Prices":        list,
		}
	}
	price := func(amount, currency string) map[string]interface{} {
		return map[string]interface{}{"Amount": amount, "Currency": currency}
	}
	const (
		evmWallet  = "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"
		solWallet  = "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM"
		prodWallet = "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"
	)

	if err := validateMonetizationConfig(build(network("BASE_SEPOLIA", evmWallet, price("0.010", "USDC")))); err != nil {
		t.Fatalf("valid config err = %v", err)
	}
	if err := validateMonetizationConfig(nil); err != nil {
		t.Fatalf("nil config err = %v", err)
	}

	invalid := []struct {
		name   string
		config map[string]interface{}
	}{
		{"bad chain", build(network("ETHEREUM", evmWallet, price("0.010", "USDC")))},
		{"zero networks", build()},
		{
			"three networks",
			build(
				network("BASE_SEPOLIA", evmWallet, price("0.010", "USDC")),
				network("SOLANA_DEVNET", solWallet, price("1", "USDC")),
				network("BASE", prodWallet, price("1", "USDC")),
			),
		},
		{"no prices", build(network("BASE_SEPOLIA", evmWallet))},
		{"non-USDC currency", build(network("BASE_SEPOLIA", evmWallet, price("0.010", "ETH")))},
		{"amount below minimum", build(network("BASE_SEPOLIA", evmWallet, price("0.0001", "USDC")))},
		{"amount above maximum", build(network("BASE_SEPOLIA", evmWallet, price("1000000000", "USDC")))},
		{"amount with four decimals", build(network("BASE_SEPOLIA", evmWallet, price("0.1234", "USDC")))},
		{
			"mixed production and test networks",
			build(
				network("BASE_SEPOLIA", evmWallet, price("0.010", "USDC")),
				network("SOLANA", solWallet, price("1", "USDC")),
			),
		},
	}
	for _, tc := range invalid {
		if err := validateMonetizationConfig(tc.config); err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
		}
	}

	mixedMode := build(network("BASE_SEPOLIA", evmWallet, price("0.010", "USDC")))
	mixedMode["CurrencyMode"] = "PLAY"
	if err := validateMonetizationConfig(mixedMode); err == nil {
		t.Error("bad currency mode accepted")
	}
}

func TestValidWalletAddress(t *testing.T) {
	// Checksum vectors from the EIP-55 specification.
	for _, address := range []string{
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
		"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
		"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
		"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
	} {
		if !validWalletAddress("BASE", address) {
			t.Errorf("EIP-55 vector %s rejected", address)
		}
	}
	// One-case addresses bypass the checksum, per the WalletAddress
	// documentation.
	if !validWalletAddress("BASE", strings.ToLower("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed")) {
		t.Error("all-lowercase address must bypass the checksum")
	}
	// A mixed-case address with a corrupted checksum digit fails.
	if validWalletAddress("BASE", "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1bEaed") &&
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1bEaed" != "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed" {
		t.Error("corrupted checksum accepted")
	}
	if validWalletAddress("BASE", "0x123") {
		t.Error("short EVM address accepted")
	}
	if validWalletAddress("BASE", "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaez") {
		t.Error("non-hex EVM address accepted")
	}
	// A Solana public key (Base58, 32-44 characters).
	if !validWalletAddress("SOLANA", "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM") {
		t.Error("valid Solana address rejected")
	}
	if validWalletAddress("SOLANA", "0I1O0l") {
		t.Error("Base58-excluded characters accepted")
	}
	if validWalletAddress("SOLANA", "2short") {
		t.Error("short Solana address accepted")
	}
}

func TestValidateMonetizeRules(t *testing.T) {
	monetizeRule := func(statement *wafstore.Statement) *wafstore.Rule {
		return &wafstore.Rule{
			Name:     "monetize-rule",
			Priority: 1,
			Action: &wafstore.Action{Monetize: &wafstore.MonetizeAction{
				PriceMultiplier: "2",
			}},
			Statement: statement,
		}
	}
	uriStatement := &wafstore.Statement{
		ByteMatchStatement: &wafstore.ByteMatchStatement{
			FieldToMatch:         &wafstore.FieldToMatch{UriPath: &wafstore.All{}},
			SearchString:         []byte("/premium"),
			PositionalConstraint: "CONTAINS",
			TextTransformations:  []*wafstore.TextTransformation{{Priority: 0, Type: "NONE"}},
		},
	}
	config := map[string]interface{}{"CurrencyMode": "TEST"}

	if err := validateMonetizeRules([]*wafstore.Rule{monetizeRule(uriStatement)}, "CLOUDFRONT", config); err != nil {
		t.Fatalf("valid monetize rule err = %v", err)
	}
	if err := validateMonetizeRules([]*wafstore.Rule{}, "REGIONAL", nil); err != nil {
		t.Fatalf("rule list without monetize rules err = %v", err)
	}
	if err := validateMonetizeRules([]*wafstore.Rule{monetizeRule(uriStatement)}, "CLOUDFRONT", nil); err == nil {
		t.Error("monetize rule without MonetizationConfig accepted")
	}
	if err := validateMonetizeRules([]*wafstore.Rule{monetizeRule(uriStatement)}, "REGIONAL", config); err == nil {
		t.Error("monetize rule on a regional web ACL accepted")
	}
	rateStatement := &wafstore.Statement{
		RateBasedStatement: &wafstore.RateBasedStatement{Limit: 100, AggregateKeyType: "IP"},
	}
	if err := validateMonetizeRules([]*wafstore.Rule{monetizeRule(rateStatement)}, "CLOUDFRONT", config); err == nil {
		t.Error("monetize rule on a rate-based statement accepted")
	}
	networkConfig := func(amount string) map[string]interface{} {
		return map[string]interface{}{
			"CurrencyMode": "TEST",
			"CryptoConfig": map[string]interface{}{
				"PaymentNetworks": []interface{}{map[string]interface{}{
					"Chain":         "BASE_SEPOLIA",
					"WalletAddress": "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
					"Prices":        []interface{}{map[string]interface{}{"Amount": amount, "Currency": "USDC"}},
				}},
			},
		}
	}
	if err := validateMonetizationConfig(networkConfig("0.010")); err != nil {
		t.Fatalf("network with a well-formed price err = %v", err)
	}
	for _, amount := range []string{"01.500", "1.", "+1", ".5"} {
		if err := validateMonetizationConfig(networkConfig(amount)); err == nil {
			t.Errorf("price amount %q accepted", amount)
		}
	}
}

func TestValidateImmunityConfig(t *testing.T) {
	config := func(seconds int64) map[string]interface{} {
		return map[string]interface{}{
			"ImmunityTimeProperty": map[string]interface{}{"ImmunityTime": seconds},
		}
	}
	if err := validateImmunityConfig(nil, "Captcha"); err != nil {
		t.Fatalf("nil config err = %v", err)
	}
	if err := validateImmunityConfig(map[string]interface{}{}, "Captcha"); err != nil {
		t.Fatalf("config without ImmunityTimeProperty err = %v", err)
	}
	if err := validateImmunityConfig(config(60), "Captcha"); err != nil {
		t.Fatalf("captcha immunity 60 err = %v", err)
	}
	if err := validateImmunityConfig(config(300), "Challenge"); err != nil {
		t.Fatalf("challenge immunity 300 err = %v", err)
	}
	for _, tc := range []struct {
		action   string
		immunity int64
	}{
		{"Captcha", 0},
		{"Captcha", 59},
		{"Captcha", 259201},
		{"Challenge", 299},
	} {
		if err := validateImmunityConfig(config(tc.immunity), tc.action); err == nil {
			t.Errorf("%s immunity %d accepted", tc.action, tc.immunity)
		}
	}
}

func TestValidateRulesManagedRuleGroupPlacement(t *testing.T) {
	managedRule := &wafstore.Rule{
		Name:     "managed-rule",
		Priority: 1,
		Action:   &wafstore.Action{Count: &wafstore.CountAction{}},
		Statement: &wafstore.Statement{
			ManagedRuleGroupStatement: &wafstore.ManagedRuleGroupStatement{
				VendorName: "AWS",
				Name:       "AWSManagedRulesCommonRuleSet",
			},
		},
	}
	if err := validateRules([]*wafstore.Rule{managedRule}, true); err != nil {
		t.Fatalf("web ACL rule with a top-level managed rule group err = %v", err)
	}
	if err := validateRules([]*wafstore.Rule{managedRule}, false); err == nil {
		t.Error("managed rule group inside a rule group accepted")
	}
	nested := &wafstore.Rule{
		Name:     "nested-managed-rule",
		Priority: 1,
		Action:   &wafstore.Action{Count: &wafstore.CountAction{}},
		Statement: &wafstore.Statement{
			NotStatement: &wafstore.NotStatement{Statement: managedRule.Statement},
		},
	}
	if err := validateRules([]*wafstore.Rule{nested}, true); err == nil {
		t.Error("managed rule group nested inside a NotStatement accepted")
	}
}
