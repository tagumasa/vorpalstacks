package cloudtrail

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// Delivery of Lake query results to S3 (StartQuery DeliveryS3Uri). After a
// query FINISHES, the rows are delivered as numbered gzip CSV files plus a
// JSON sign file under
// [prefix/]AWSLogs/<account>/CloudTrail-Lake/Query/<Y>/<M>/<D>/<queryID>/.
// The sign file carries the SHA-256 of each compressed file and an
// RSA-SHA256 signature over the space-joined hash list, made with a
// per-delivery signing key whose public half is stored so ListPublicKeys
// can serve it to validators matching the fingerprint.

// deliveryS3URIPattern is the model's DeliveryS3Uri trait: an S3 URI whose
// bucket part follows the S3 bucket-name rules, with an optional key
// prefix.
const deliveryS3URIPattern = `^s3://[a-z0-9][.\-a-z0-9]{1,61}[a-z0-9](/.*)?$`

var deliveryS3URIRe = regexp.MustCompile(deliveryS3URIPattern)

// splitDeliveryS3URI validates a DeliveryS3Uri against the model pattern
// and splits it into bucket and key prefix (no trailing slash).
func splitDeliveryS3URI(uri string) (bucket, prefix string, err error) {
	if !deliveryS3URIRe.MatchString(uri) {
		return "", "", newInvalidS3BucketNameException(
			fmt.Sprintf("The DeliveryS3Uri is not valid: %s", uri))
	}
	rest := strings.TrimPrefix(uri, "s3://")
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		return rest[:i], strings.Trim(rest[i+1:], "/"), nil
	}
	return rest, "", nil
}

// queryResultSignFile is the structure of result_sign.json: the delivered
// file names with the SHA-256 of their compressed content, the algorithms
// used, the delivery time, the signature over the hash list, and the
// fingerprint identifying the signing public key.
type queryResultSignFile struct {
	Version              string
	Region               string
	Files                []queryResultSignFileEntry
	HashAlgorithm        string
	SignatureAlgorithm   string
	QueryCompleteTime    string
	HashSignature        string
	PublicKeyFingerprint string
}

// queryResultSignFileEntry names one delivered query result file and the
// hexadecimal SHA-256 of its compressed content.
type queryResultSignFileEntry struct {
	FileHashValue string
	FileName      string
}

// queryResultFileSink receives each result file as it is finalised: the
// compressed content is handed over for delivery and released when the
// sink returns, so the delivery path holds one file's content at a time.
type queryResultFileSink func(name string, content []byte, hashHex string) error

// buildQueryResultFiles renders the query rows as gzip CSV files named
// result_N.csv.gz, handing each completed file to sink. The header row
// carries the statement's projected column spellings (the full Lake
// vocabulary for SELECT *), so each file is a self-describing table; the
// row order matches GetQueryResults. A new file starts when the rendered
// CSV reaches maxFileBytes, the platform operational bound that keeps one
// file's footprint small — the AWS-documented per-file ceiling is larger
// (MaxQueryResultFileBytes) and smaller numbered files stay inside the
// documented form. A query with no rows delivers a single header-only
// file. The returned entries carry each delivered file's name and
// content hash for the sign file.
func buildQueryResultFiles(pq *parsedQuery, rows [][]map[string]string, maxFileBytes int, sink queryResultFileSink) ([]queryResultSignFileEntry, error) {
	columns := pq.columns
	if len(columns) == 1 && columns[0] == "*" {
		columns = lakeColumnList
	}

	entries := []queryResultSignFileEntry{}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	rowsInChunk := 0

	flushFile := func() error {
		w.Flush()
		if err := w.Error(); err != nil {
			return err
		}
		compressed, sum, err := gzipAndHash(buf.Bytes())
		if err != nil {
			return err
		}
		entry := queryResultSignFileEntry{
			FileHashValue: hex.EncodeToString(sum[:]),
			FileName:      fmt.Sprintf("result_%d.csv.gz", len(entries)+1),
		}
		if err := sink(entry.FileName, compressed, entry.FileHashValue); err != nil {
			return err
		}
		entries = append(entries, entry)
		buf.Reset()
		w = csv.NewWriter(&buf)
		rowsInChunk = 0
		return nil
	}

	writeHeader := func() error {
		return w.Write(columns)
	}

	if err := writeHeader(); err != nil {
		return nil, err
	}
	for _, row := range rows {
		values := make([]string, len(columns))
		for i, col := range columns {
			canonical := lakeColumnCanonical(col)
			for _, entry := range row {
				if v, ok := entry[col]; ok {
					values[i] = v
					break
				}
				if v, ok := entry[canonical]; ok {
					values[i] = v
					break
				}
			}
		}
		if err := w.Write(values); err != nil {
			return nil, err
		}
		rowsInChunk++
		// Flush per row so the buffer length reflects the rendered CSV and
		// the size bound splits at real file sizes.
		w.Flush()
		if err := w.Error(); err != nil {
			return nil, err
		}
		if buf.Len() >= maxFileBytes {
			if err := flushFile(); err != nil {
				return nil, err
			}
			if err := writeHeader(); err != nil {
				return nil, err
			}
		}
	}
	if rowsInChunk > 0 || len(entries) == 0 {
		if err := flushFile(); err != nil {
			return nil, err
		}
	}
	return entries, nil
}

// deliverQueryResults writes the finished query's rows and sign file to
// the query's DeliveryS3Uri target and records the delivery status on the
// query record. A delivery failure sets DeliveryStatus and never disturbs
// the query's own FINISHED status. The recorded status starts at FAILED
// and only the fully delivered path upgrades it to SUCCESS, so a delivery
// interrupted before its end — including one interrupted by a panic —
// never records success for files that were not written.
func (s *CloudTrailService) deliverQueryResults(store cloudtrailstore.CloudTrailStoreInterface, qr *cloudtrailstore.QueryRecord) {
	status := "FAILED"
	defer func() {
		if _, err := store.MutateQuery(qr.QueryID, func(rec *cloudtrailstore.QueryRecord) error {
			rec.DeliveryStatus = status
			return nil
		}); err != nil {
			slog.Error("Failed to record query delivery status", "queryId", qr.QueryID, "error", err)
		}
	}()

	invoker := s.s3Invoker()
	if invoker == nil {
		status = "FAILED"
		return
	}
	bucket, prefix, err := splitDeliveryS3URI(qr.DeliveryS3URI)
	if err != nil {
		status = "FAILED"
		return
	}
	// A bucket that existed at admission may be gone by delivery time;
	// that outcome is the delivery-status RESOURCE_NOT_FOUND, not a
	// generic failure.
	exists, err := invoker.BucketExists(s.ctx, store.GetRegion(), bucket)
	if err != nil {
		status = "FAILED"
		return
	}
	if !exists {
		status = "RESOURCE_NOT_FOUND"
		return
	}

	pq, err := parseQueryStatement(qr.QueryStatement)
	if err != nil {
		status = "FAILED"
		return
	}

	// The per-delivery signing key: the private half signs the hash list,
	// the stored public half is what validators retrieve through
	// ListPublicKeys by fingerprint.
	pub, priv, err := store.CreateAndStoreSigningKey()
	if err != nil {
		status = "FAILED_SIGNING_FILE"
		return
	}

	// The path date is the query's completion day: validators locate the
	// delivery under the day the query finished.
	completed := time.Now().UTC()
	if qr.EndTime != nil {
		completed = qr.EndTime.UTC()
	}
	base := fmt.Sprintf("AWSLogs/%s/CloudTrail-Lake/Query/%04d/%02d/%02d/%s",
		store.GetAccountID(), completed.Year(), completed.Month(), completed.Day(), qr.QueryID)
	if prefix != "" {
		base = prefix + "/" + base
	}

	// Each completed file uploads and releases immediately — the
	// delivery path holds one file's content at a time; the sign file
	// carries every file's name and hash.
	entries, err := buildQueryResultFiles(pq, qr.QueryResultRows, cloudtrailstore.MaxQueryResultFileChunkBytes, func(name string, content []byte, _ string) error {
		return invoker.PutObject(s.ctx, store.GetRegion(), bucket, base+"/"+name, content, "application/gzip")
	})
	if err != nil {
		status = "FAILED"
		return
	}

	sign, err := buildQueryResultSignFile(store.GetRegion(), entries, pub, priv)
	if err != nil {
		status = "FAILED_SIGNING_FILE"
		return
	}
	signJSON, err := json.Marshal(sign)
	if err != nil {
		status = "FAILED_SIGNING_FILE"
		return
	}
	if err := invoker.PutObject(s.ctx, store.GetRegion(), bucket, base+"/result_sign.json", signJSON, "application/json"); err != nil {
		status = "FAILED_SIGNING_FILE"
		return
	}
	status = "SUCCESS"
}

// buildQueryResultSignFile assembles and signs result_sign.json for the
// delivered files: the signature is RSA-SHA256 over the space-joined
// fileHashValue list, hexadecimal encoded.
func buildQueryResultSignFile(region string, entries []queryResultSignFileEntry, pub *cloudtrailstore.PublicKey, priv *rsa.PrivateKey) (*queryResultSignFile, error) {
	hashes := make([]string, len(entries))
	for i, e := range entries {
		hashes[i] = e.FileHashValue
	}

	dataSigningString := strings.Join(hashes, " ")
	digest := sha256.Sum256([]byte(dataSigningString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
	if err != nil {
		return nil, err
	}

	return &queryResultSignFile{
		Version:              "1.0",
		Region:               region,
		Files:                entries,
		HashAlgorithm:        "SHA-256",
		SignatureAlgorithm:   "SHA256withRSA",
		QueryCompleteTime:    time.Now().UTC().Format(time.RFC3339),
		HashSignature:        hex.EncodeToString(signature),
		PublicKeyFingerprint: pub.Fingerprint(),
	}, nil
}
