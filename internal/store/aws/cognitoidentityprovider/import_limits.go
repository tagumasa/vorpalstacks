package cognitoidentityprovider

// AWS-documented limits for CSV user import. Sources: Amazon Cognito
// developer guide, "Importing users into user pools from a CSV file"
// (Formatting the CSV file / Supported password hashing algorithms) and the
// UserImportJobType Status documentation. These are the single definitions;
// every other site must reference them.

const (
	// MaxImportCSVFileSizeBytes is the maximum CSV file size (100 MB).
	MaxImportCSVFileSizeBytes int64 = 100 << 20
	// MaxImportCSVRowLengthChars is the maximum length of one CSV row in
	// characters (16,000).
	MaxImportCSVRowLengthChars = 16000
	// MaxImportCSVRows is the maximum number of data rows in one file,
	// excluding the header row (500,000).
	MaxImportCSVRows = 500000
	// ImportJobExpiryHours is the age after which a Created job that was
	// never started becomes Expired. AWS documents the window as
	// "24-48 hours"; this implementation fixes the deterministic lower
	// bound of that window.
	ImportJobExpiryHours = 24
)

// Parameter ceilings for imported password hashes. A hash whose embedded
// parameters exceed these ceilings fails the import for that user.
const (
	// MinImportHashParamValue is the floor for every numeric hash
	// parameter: a zero-valued cost, block size, parallelism, memory, or
	// iteration count is not a legitimate hash (RFC 9106 requires at
	// least one round and one thread for Argon2id).
	MinImportHashParamValue = 1
	// MaxImportHashBcryptCost is the maximum bcrypt cost factor.
	MaxImportHashBcryptCost = 12
	// MaxImportHashScryptN is the maximum scrypt CPU/memory cost.
	MaxImportHashScryptN = 65536
	// MaxImportHashScryptR is the maximum scrypt block size.
	MaxImportHashScryptR = 8
	// MaxImportHashScryptP is the maximum scrypt parallelism.
	MaxImportHashScryptP = 1
	// MaxImportHashArgon2MemKiB is the maximum Argon2id memory in KiB.
	MaxImportHashArgon2MemKiB = 19456
	// MaxImportHashArgon2Time is the maximum Argon2id iteration count.
	MaxImportHashArgon2Time = 2
	// MaxImportHashArgon2Parallelism is the maximum Argon2id parallelism.
	MaxImportHashArgon2Parallelism = 1
	// MaxImportHashPbkdf2Iterations is the maximum PBKDF2-SHA256 iteration
	// count.
	MaxImportHashPbkdf2Iterations = 600000
)
