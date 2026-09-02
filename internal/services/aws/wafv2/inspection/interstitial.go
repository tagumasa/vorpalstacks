package inspection

import (
	"fmt"
	"strconv"
	"strings"
)

// ChallengeInterstitialHTML renders the JavaScript interstitial of an
// interrupting Captcha or Challenge response. The page is
// self-contained: an embedded parameter block carries the challenge
// identifier, the solution difficulty and the token exchange endpoint;
// the inline script computes a proof of work over the challenge
// identifier, exchanges it for an aws-waf-token cookie at the endpoint
// and then re-requests the original URL. The parameter block is a JSON
// script element so non-browser clients can parse the same challenge
// data the script consumes.
func ChallengeInterstitialHTML(kind, challengeID, tokenEndpoint string, difficulty int) string {
	var builder strings.Builder
	builder.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	builder.WriteString("<meta charset=\"utf-8\">\n")
	if kind == ActionCaptcha {
		builder.WriteString("<title>Verification required</title>\n")
	} else {
		builder.WriteString("<title>Checking your request</title>\n")
	}
	builder.WriteString("</head>\n<body>\n")
	builder.WriteString("<noscript><p>JavaScript is required to proceed.</p></noscript>\n")
	fmt.Fprintf(&builder,
		"<script type=\"application/json\" id=\"awswaf-challenge\">{\"challengeId\":%q,\"kind\":%q,\"tokenEndpoint\":%q,\"difficulty\":%d}</script>\n",
		challengeID, strings.ToLower(kind), tokenEndpoint, difficulty)
	builder.WriteString(interstitialScript)
	builder.WriteString("</body>\n</html>\n")
	return builder.String()
}

// interstitialScript solves the embedded challenge and retries the
// original request once the exchange endpoint has set the token cookie.
// The retry is a navigation to the original URL, which serves the
// browser flows the interstitial exists for; API clients integrate the
// token exchange themselves.
const interstitialScript = `<script>
(function () {
  var params = JSON.parse(document.getElementById('awswaf-challenge').textContent);
  var prefix = new Array(params.difficulty + 1).join('0');
  var encoder = new TextEncoder();
  function hex(buffer) {
    return Array.prototype.map.call(new Uint8Array(buffer), function (byte) {
      return byte.toString(16).padStart(2, '0');
    }).join('');
  }
  function solve(counter) {
    return crypto.subtle.digest('SHA-256', encoder.encode(params.challengeId + '.' + counter)).then(function (digest) {
      if (hex(digest).substring(0, params.difficulty) === prefix) {
        return counter;
      }
      return solve(counter + 1);
    });
  }
  solve(0).then(function (counter) {
    return fetch(params.tokenEndpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ challengeId: params.challengeId, counter: String(counter) })
    });
  }).then(function (response) {
    if (response.ok) {
      window.location.reload();
    }
  });
})();
</script>
`

// InterstitialChallengeID extracts the challenge identifier from an
// interstitial page, for tests and diagnostics.
func InterstitialChallengeID(html string) string {
	marker := "\"challengeId\":"
	start := strings.Index(html, marker)
	if start < 0 {
		return ""
	}
	rest := html[start+len(marker):]
	end := strings.IndexByte(rest, '"')
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// InterstitialDifficulty extracts the difficulty value from an
// interstitial page.
func InterstitialDifficulty(html string) int {
	marker := "\"difficulty\":"
	start := strings.Index(html, marker)
	if start < 0 {
		return 0
	}
	rest := html[start+len(marker):]
	end := strings.IndexAny(rest, ",}")
	if end < 0 {
		return 0
	}
	difficulty, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil {
		return 0
	}
	return difficulty
}
