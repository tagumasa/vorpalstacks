package inspection

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// applyTextTransformations runs the configured transformations against
// one request-component value, in ascending Priority order. AWS WAF
// processes all transformations from lowest priority to highest before
// inspecting the result (API_TextTransformation.Priority). The pipeline
// is byte-oriented: hash and decode transformations emit raw binary,
// which later byte matching compares against the SearchString as-is.
func applyTextTransformations(value []byte, tts []*wafstore.TextTransformation) []byte {
	ordered := make([]*wafstore.TextTransformation, len(tts))
	copy(ordered, tts)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Priority < ordered[j].Priority
	})
	out := value
	for _, tt := range ordered {
		if tt == nil {
			continue
		}
		out = applyTextTransformation(out, tt.Type)
	}
	return out
}

// applyTextTransformation applies a single transformation type. The
// behaviours follow the per-type descriptions in the AWS WAF Developer
// Guide section "Using text transformations"
// (waf-rule-statement-transformation.html).
func applyTextTransformation(value []byte, t string) []byte {
	switch t {
	case "NONE":
		return value
	case "COMPRESS_WHITE_SPACE":
		return compressWhiteSpace(value)
	case "HTML_ENTITY_DECODE":
		return htmlEntityDecode(value)
	case "LOWERCASE":
		return []byte(strings.ToLower(string(value)))
	case "UPPERCASE":
		return []byte(strings.ToUpper(string(value)))
	case "CMD_LINE":
		return cmdLineTransform(value, cmdLineGeneric)
	case "CMD_LINE_UNIX":
		return cmdLineTransform(value, cmdLineUnix)
	case "CMD_LINE_WIN":
		return cmdLineTransform(value, cmdLineWindows)
	case "URL_DECODE":
		return urlDecode(value)
	case "URL_DECODE_UNI":
		return urlDecodeUni(value)
	case "BASE64_DECODE":
		return base64Decode(value, false)
	case "BASE64_DECODE_EXT":
		return base64Decode(value, true)
	case "HEX_DECODE":
		return hexDecode(value)
	case "SQL_HEX_DECODE":
		return sqlHexDecode(value)
	case "MD5":
		sum := md5.Sum(value) //nolint:gosec // MD5 is the AWS-specified transformation, not a security choice.
		return sum[:]
	case "SHA256":
		sum := sha256.Sum256(value)
		return sum[:]
	case "REPLACE_COMMENTS":
		return replaceComments(value)
	case "REMOVE_COMMENTS_CHAR":
		return removeCommentsChar(value)
	case "ESCAPE_SEQ_DECODE":
		return escapeSeqDecode(value)
	case "CSS_DECODE":
		return cssDecode(value)
	case "JS_DECODE":
		return jsDecode(value)
	case "JS_DECODE_EXT":
		return jsDecodeExt(value)
	case "NORMALIZE_PATH":
		return normalizePath(value, false)
	case "NORMALIZE_PATH_WIN":
		return normalizePath(value, true)
	case "REMOVE_NULLS":
		return bytesReplaceAll(value, 0, -1)
	case "REPLACE_NULLS":
		return bytesReplaceAll(value, 0, ' ')
	case "REMOVE_WHITESPACE":
		return removeWhiteSpace(value)
	case "TRIM":
		return trimASCIISpaceBytes(value)
	case "TRIM_LEFT":
		return trimASCIISpaceLeft(value)
	case "TRIM_RIGHT":
		return trimASCIISpaceRight(value)
	case "UTF8_TO_UNICODE":
		return utf8ToUnicode(value)
	default:
		// Unknown types reach here only if validation was bypassed;
		// leaving the value untouched is the least surprising outcome.
		return value
	}
}

// cmdLine variant selectors. The generic variant deletes the characters
// `\ " ' ^`, the Unix variant deletes `\ " '`, and the Windows variant
// deletes `" ' ^` while collapsing runs of backslashes (single
// backslashes survive so Windows paths stay intact).
const (
	cmdLineGeneric = 0
	cmdLineUnix    = 1
	cmdLineWindows = 2
)

func cmdLineTransform(value []byte, variant int) []byte {
	s := string(value)
	if variant == cmdLineWindows {
		s = strings.ReplaceAll(s, "^\n", "")
		s = strings.ReplaceAll(s, "^\r\n", "")
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\', '"', '\'', '^':
			if (variant == cmdLineUnix && c == '^') ||
				(variant == cmdLineWindows && c == '\\') {
				// Unix keeps carets; Windows keeps backslashes (runs are
				// collapsed after the scan below).
				b.WriteByte(c)
			}
			continue
		case '\t', '\n', '\r', '\v', '\f':
			if variant == cmdLineGeneric {
				b.WriteByte(c)
			} else {
				b.WriteByte(' ')
			}
			continue
		case ' ', ',', ';':
			if variant == cmdLineGeneric && (c == ',' || c == ';') {
				b.WriteByte(' ')
			} else {
				b.WriteByte(c)
			}
			continue
		default:
			b.WriteByte(c)
		}
	}
	s = b.String()
	if variant == cmdLineWindows {
		for strings.Contains(s, "\\\\") {
			s = strings.ReplaceAll(s, "\\\\", "\\")
		}
	}
	// All variants collapse multiple spaces and lowercase A-Z; the two
	// platform variants additionally trim outer spaces.
	s = collapseSpaces(s, variant != cmdLineGeneric)
	s = lowercaseASCII(s)
	if variant == cmdLineGeneric {
		// The generic variant additionally deletes spaces before the
		// characters / and (, per the Developer Guide's CMD_LINE entry.
		s = strings.ReplaceAll(s, " /", "/")
		s = strings.ReplaceAll(s, " (", "(")
	}
	return []byte(s)
}

// collapseSpaces replaces runs of spaces with one space and optionally
// trims leading and trailing spaces. It does not change character case:
// the Developer Guide defines case folding only for the command-line
// transformations, so COMPRESS_WHITE_SPACE must preserve the input's
// case (the guide's compress-whitespace entry lists no case change).
func collapseSpaces(s string, trim bool) string {
	var b strings.Builder
	inSpace := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' {
			inSpace = true
			continue
		}
		if inSpace && b.Len() > 0 {
			b.WriteByte(' ')
		}
		inSpace = false
		b.WriteByte(c)
	}
	return b.String()
}

// lowercaseASCII converts uppercase letters A-Z to lowercase a-z, the
// exact case folding the Developer Guide specifies for the command-line
// transformations.
func lowercaseASCII(s string) string {
	out := []byte(s)
	for i, c := range out {
		if 'A' <= c && c <= 'Z' {
			out[i] = c + 'a' - 'A'
		}
	}
	return string(out)
}

// compressWhiteSpace replaces formfeed, tab, newline, carriage return,
// vertical tab and non-breaking space (ASCII 160) with a space, then
// collapses multiple spaces into one.
func compressWhiteSpace(value []byte) []byte {
	out := make([]byte, 0, len(value))
	for _, c := range value {
		switch c {
		case '\f', '\t', '\n', '\r', '\v', 160:
			out = append(out, ' ')
		default:
			out = append(out, c)
		}
	}
	return []byte(collapseSpaces(string(out), false))
}

// htmlEntityDecode replaces numeric entities (&#xhhhh; / &#nnnn;) and
// the named entities with their characters. Entity handling is
// case insensitive, per the Developer Guide description.
func htmlEntityDecode(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] != '&' {
			b.WriteByte(s[i])
			i++
			continue
		}
		semi := strings.IndexByte(s[i:], ';')
		// The window must accommodate the longest documented named
		// entity (&NonBreakingSpace;, &DiacriticalGrave;).
		if semi < 2 || semi > 24 {
			b.WriteByte(s[i])
			i++
			continue
		}
		entity := s[i : i+semi+1]
		if decoded, ok := decodeHTMLEntity(entity); ok {
			b.WriteString(decoded)
			i += len(entity)
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return []byte(b.String())
}

func decodeHTMLEntity(entity string) (string, bool) {
	body := entity[1 : len(entity)-1]
	if len(body) == 0 {
		return "", false
	}
	if body[0] == '#' {
		var r rune
		var err error
		if len(body) > 1 && (body[1] == 'x' || body[1] == 'X') {
			var v int64
			v, err = parseInt64(body[2:], 16)
			r = rune(v)
		} else {
			var v int64
			v, err = parseInt64(body[1:], 10)
			r = rune(v)
		}
		if err != nil || r == 0 || r > 0x10FFFF {
			return "", false
		}
		return string(r), true
	}
	if decoded, ok := namedEntities[strings.ToLower(body)]; ok {
		return decoded, true
	}
	return "", false
}

// namedEntities maps the named entities the Developer Guide's HTML
// entity decode entry lists ("AWS WAF replaces the following
// HTML-encoded characters with unencoded characters"). Keys are the
// lower-cased entity bodies because the guide states the handling is
// case insensitive. Where the guide's table offers several spellings
// for one character (for example lcub and lbrace), every spelling is a
// key. Two entries required a judgement call the guide's rendered table
// leaves ambiguous: apos is mapped to the HTML-standard apostrophe '
// (the table cell renders a backslash, which the backslash entity bsol
// already covers, so the rendering is taken to be an escaping artefact)
// and minus is mapped to the plain hyphen the table displays.
var namedEntities = map[string]string{
	"quot": "\"", "amp": "&", "lt": "<", "gt": ">",
	"nbsp": "\u00a0", "nonbreakingspace": "\u00a0",
	"newline": "\n", "tab": "\t",
	"lcub": "{", "lbrace": "{", "rcub": "}", "rbrace": "}",
	"excl": "!", "num": "#", "dollar": "$", "percnt": "%",
	"apos": "'", "lpar": "(", "rpar": ")",
	"ast": "*", "midast": "*", "plus": "+", "comma": ",",
	"period": ".", "sol": "/", "colon": ":", "semi": ";",
	"equals": "=", "quest": "?",
	"verbar": "|", "vert": "|", "verticalline": "|",
	"lsqb": "[", "lbrack": "[", "rsqb": "]", "rbrack": "]",
	"hat": "^", "lowbar": "_", "underbar": "_",
	"grave": "`", "diacriticalgrave": "`",
	"tilde": "~", "diacriticaltilde": "~",
	"minus": "-", "bsol": "\\",
}

func urlDecode(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c == '%' && i+2 < len(s) {
			hi, okHi := hexNibble(s[i+1])
			lo, okLo := hexNibble(s[i+2])
			if okHi && okLo {
				b.WriteByte(hi<<4 | lo)
				i += 3
				continue
			}
		}
		if c == '+' {
			b.WriteByte(' ')
			i++
			continue
		}
		b.WriteByte(c)
		i++
	}
	return []byte(b.String())
}

func hexNibble(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// urlDecodeUni decodes percent-encoding plus Microsoft %uXXXX encoding.
// Full-width ASCII codes (FF01-FF5E) are adjusted to their ASCII
// counterparts; otherwise the low byte is kept and the high byte is
// zeroed.
func urlDecodeUni(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '%' && i+5 < len(s) && (s[i+1] == 'u' || s[i+1] == 'U') {
			var code int
			valid := true
			for k := 0; k < 4; k++ {
				n, ok := hexNibble(s[i+2+k])
				if !ok {
					valid = false
					break
				}
				code = code<<4 | int(n)
			}
			if valid {
				b.WriteRune(adjustFullWidthRune(rune(code)))
				i += 6
				continue
			}
		}
		if s[i] == '%' && i+2 < len(s) {
			hi, okHi := hexNibble(s[i+1])
			lo, okLo := hexNibble(s[i+2])
			if okHi && okLo {
				b.WriteByte(hi<<4 | lo)
				i += 3
				continue
			}
		}
		if s[i] == '+' {
			b.WriteByte(' ')
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return []byte(b.String())
}

// base64Decode decodes Base64. The forgiving variant ignores invalid
// characters (BASE64_DECODE_EXT); the strict variant returns the input
// unchanged when decoding fails — AWS does not document failure
// behaviour for the strict form, and leaving the value undecoded keeps
// matching deterministic instead of silently dropping the component.
func base64Decode(value []byte, forgiving bool) []byte {
	if forgiving {
		filtered := make([]byte, 0, len(value))
		for _, c := range value {
			if isBase64Char(c) {
				filtered = append(filtered, c)
			}
		}
		decoded, err := base64.StdEncoding.DecodeString(string(filtered))
		if err != nil {
			return value
		}
		return decoded
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(value)))
	if err != nil {
		return value
	}
	return decoded
}

func isBase64Char(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' ||
		c == '+' || c == '/' || c == '='
}

// hexDecode decodes a string of hexadecimal characters into binary.
func hexDecode(value []byte) []byte {
	s := strings.TrimSpace(string(value))
	if len(s)%2 != 0 {
		return value
	}
	out := make([]byte, 0, len(s)/2)
	for i := 0; i+1 < len(s); i += 2 {
		hi, okHi := hexNibble(s[i])
		lo, okLo := hexNibble(s[i+1])
		if !okHi || !okLo {
			return value
		}
		out = append(out, hi<<4|lo)
	}
	return out
}

// sqlHexDecode decodes SQL hex literals of the form 0x414243 into the
// corresponding bytes.
func sqlHexDecode(value []byte) []byte {
	s := string(value)
	if len(s) < 3 || (s[0] != '0' || (s[1] != 'x' && s[1] != 'X')) {
		return value
	}
	decoded := hexDecode([]byte(s[2:]))
	if len(decoded) == 0 && len(s) > 2 {
		return value
	}
	return decoded
}

// replaceComments replaces each complete C-style comment with one space
// (without compressing consecutive occurrences), replaces an
// unterminated comment with one space, and leaves a standalone `*/`
// untouched.
func replaceComments(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '/' && i+1 < len(s) && s[i+1] == '*' {
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				b.WriteByte(' ')
				break
			}
			b.WriteByte(' ')
			i = i + 2 + end + 2
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return []byte(b.String())
}

// removeCommentsChar removes the common comment markers /*, */, -- and
// # from the input.
func removeCommentsChar(value []byte) []byte {
	s := string(value)
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	s = strings.ReplaceAll(s, "--", "")
	s = strings.ReplaceAll(s, "#", "")
	return []byte(s)
}

// escapeSeqDecode decodes ANSI C escape sequences; sequences that are
// not valid remain in the output.
func escapeSeqDecode(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			i++
			continue
		}
		next := s[i+1]
		switch next {
		case 'a':
			b.WriteByte('\a')
			i += 2
		case 'b':
			b.WriteByte('\b')
			i += 2
		case 'f':
			b.WriteByte('\f')
			i += 2
		case 'n':
			b.WriteByte('\n')
			i += 2
		case 'r':
			b.WriteByte('\r')
			i += 2
		case 't':
			b.WriteByte('\t')
			i += 2
		case 'v':
			b.WriteByte('\v')
			i += 2
		case '\\':
			b.WriteByte('\\')
			i += 2
		case '?':
			b.WriteByte('?')
			i += 2
		case '\'':
			b.WriteByte('\'')
			i += 2
		case '"':
			b.WriteByte('"')
			i += 2
		case 'x':
			if n, adv, ok := decodeHexRunes(s[i+2:]); ok {
				b.WriteByte(n)
				i += 2 + adv
			} else {
				b.WriteByte(c)
				i++
			}
		case '0':
			if n, adv, ok := decodeOctalRunes(s[i+2:]); ok {
				b.WriteByte(n)
				i += 2 + adv
			} else {
				b.WriteByte(c)
				i++
			}
		default:
			b.WriteByte(c)
			i++
		}
	}
	return []byte(b.String())
}

func decodeHexRunes(s string) (byte, int, bool) {
	var v int
	used := 0
	for used < 2 && used < len(s) {
		n, ok := hexNibble(s[used])
		if !ok {
			break
		}
		v = v<<4 | int(n)
		used++
	}
	if used == 0 {
		return 0, 0, false
	}
	return byte(v), used, true
}

func decodeOctalRunes(s string) (byte, int, bool) {
	var v int
	used := 0
	for used < 3 && used < len(s) {
		c := s[used]
		if c < '0' || c > '7' {
			break
		}
		v = v<<3 | int(c-'0')
		used++
	}
	if used == 0 {
		return 0, 0, false
	}
	return byte(v), used, true
}

// cssDecode decodes CSS 2.x escape rules: a backslash followed by up to
// four hexadecimal digits (the two decoding bytes the Developer Guide
// allots this transformation) and an optional single whitespace
// character is replaced by the decoded character; a backslash followed
// by a non-hexadecimal character drops the backslash (uncovering
// evasion combinations such as `ja\vascript`).
func cssDecode(value []byte) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '\\' {
			if n, adv, ok := decodeCSSHex(s[i+1:]); ok {
				b.WriteRune(rune(n))
				i += 1 + adv
				// An escaped code point may absorb one whitespace.
				if i < len(s) && isASCIISpaceByte(s[i]) {
					i++
				}
				continue
			}
			if i+1 < len(s) {
				b.WriteByte(s[i+1])
				i += 2
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return []byte(b.String())
}

func decodeCSSHex(s string) (int, int, bool) {
	var v int
	used := 0
	for used < 4 && used < len(s) {
		n, ok := hexNibble(s[used])
		if !ok {
			break
		}
		v = v<<4 | int(n)
		used++
	}
	if used == 0 {
		return 0, 0, false
	}
	return v, used, true
}

// jsDecode decodes JavaScript escape sequences. A \uHHHH code in the
// full-width ASCII range FF01-FF5E is adjusted using its high byte;
// otherwise the low byte is kept and the high byte is zeroed. Classic
// single-character escapes are decoded to their characters.
func jsDecode(value []byte) []byte {
	return jsDecodeVariant(value, false)
}

// jsDecodeExt behaves like jsDecode but preserves Windows-style paths:
// recognised single-character escapes and unrecognised \C sequences are
// kept as-is, and only \\, \/, \' and \" lose their backslash.
func jsDecodeExt(value []byte) []byte {
	return jsDecodeVariant(value, true)
}

func jsDecodeVariant(value []byte, ext bool) []byte {
	s := string(value)
	var b strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			i++
			continue
		}
		next := s[i+1]
		if next == 'u' && i+5 < len(s) {
			var code int
			valid := true
			for k := 0; k < 4; k++ {
				n, ok := hexNibble(s[i+2+k])
				if !ok {
					valid = false
					break
				}
				code = code<<4 | int(n)
			}
			if valid {
				b.WriteRune(adjustFullWidthRune(rune(code)))
				i += 6
				continue
			}
		}
		switch next {
		case '\\', '/', '\'', '"':
			b.WriteByte(next)
			i += 2
		case 'a', 'b', 'f', 'n', 'r', 't', 'v':
			if ext {
				b.WriteByte(c)
				b.WriteByte(next)
			} else {
				b.WriteByte(jsEscapeChar(next))
			}
			i += 2
		default:
			if ext {
				b.WriteByte(c)
				b.WriteByte(next)
				i += 2
			} else {
				b.WriteByte(next)
				i += 2
			}
		}
	}
	return []byte(b.String())
}

// adjustFullWidthRune keeps the high-byte adjustment the Developer
// Guide describes: full-width ASCII FF01-FF5E maps to its ASCII
// counterpart (U+FF01 corresponds to U+0021, so the offset is 0xFEE0),
// and other codes keep only the low byte (the high byte is zeroed).
func adjustFullWidthRune(r rune) rune {
	if r >= 0xFF01 && r <= 0xFF5E {
		return r - 0xFEE0
	}
	return r & 0xFF
}

func jsEscapeChar(c byte) byte {
	switch c {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	default:
		return '\v'
	}
}

// normalizePath removes multiple slashes, directory self-references
// and non-leading parent references. The Windows variant first
// converts backslashes to forward slashes.
func normalizePath(value []byte, windows bool) []byte {
	s := string(value)
	if windows {
		s = strings.ReplaceAll(s, "\\", "/")
	}
	parts := strings.Split(s, "/")
	out := make([]string, 0, len(parts))
	for idx, part := range parts {
		switch part {
		case "":
			continue
		case ".":
			continue
		case "..":
			// Parent references at the very beginning are preserved.
			if idx == 0 || len(out) == 0 {
				out = append(out, part)
			} else if out[len(out)-1] != ".." {
				out = out[:len(out)-1]
			} else {
				out = append(out, part)
			}
		default:
			out = append(out, part)
		}
	}
	// Preserve leading and trailing separators when the input had one.
	leading := strings.HasPrefix(s, "/")
	trailing := strings.HasSuffix(s, "/") && !strings.HasSuffix(s, "//")
	joined := strings.Join(out, "/")
	if leading {
		joined = "/" + joined
	}
	if trailing && joined != "" {
		joined += "/"
	}
	return []byte(joined)
}

func removeWhiteSpace(value []byte) []byte {
	out := make([]byte, 0, len(value))
	for _, c := range value {
		if !isASCIISpaceByte(c) && c != 160 {
			out = append(out, c)
		}
	}
	return out
}

func isASCIISpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

func trimASCIISpaceBytes(value []byte) []byte {
	return []byte(trimASCIISpace(string(value)))
}

func trimASCIISpaceLeft(value []byte) []byte {
	s := string(value)
	for len(s) > 0 && isASCIISpaceByte(s[0]) {
		s = s[1:]
	}
	return []byte(s)
}

func trimASCIISpaceRight(value []byte) []byte {
	s := string(value)
	for len(s) > 0 && isASCIISpaceByte(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return []byte(s)
}

// utf8ToUnicode converts all UTF-8 character sequences to \uXXXX escape
// form, normalising multilingual input so that byte-level matching sees
// one canonical representation.
func utf8ToUnicode(value []byte) []byte {
	var b strings.Builder
	for _, r := range string(value) {
		if r < 0x80 {
			b.WriteRune(r)
			continue
		}
		b.WriteString("\\u")
		b.WriteString(strings.ToLower(hex.EncodeToString([]byte(string(r)))))
	}
	return []byte(b.String())
}

// bytesReplaceAll removes (replacement < 0) or replaces every
// occurrence of old.
func bytesReplaceAll(value []byte, old byte, replacement int) []byte {
	out := make([]byte, 0, len(value))
	for _, c := range value {
		if c == old {
			if replacement >= 0 {
				out = append(out, byte(replacement))
			}
			continue
		}
		out = append(out, c)
	}
	return out
}

func parseInt64(s string, base int) (int64, error) {
	var v int64
	if len(s) == 0 {
		return 0, errShortNumber
	}
	for i := 0; i < len(s); i++ {
		n, ok := hexNibble(s[i])
		if !ok || int64(n) >= int64(base) {
			return 0, errBadDigit
		}
		v = v*int64(base) + int64(n)
	}
	return v, nil
}

type parseError string

func (e parseError) Error() string { return string(e) }

const (
	errShortNumber parseError = "empty numeric entity"
	errBadDigit    parseError = "invalid numeric entity digit"
)
