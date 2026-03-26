// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

// ParsedARN represents a parsed Amazon Resource Name with its constituent components.
type ParsedARN struct {
	Partition string // The partition (e.g., "aws" or "aws-cn")
	Service   string // The AWS service name (e.g., "s3", "ec2")
	Region    string // The region (e.g., "us-east-1")
	AccountID string // The account ID
	Resource  string // The resource identifier
	RoleName  string // The role name (if applicable)
}

// IsCrossAccount returns true if the ARN belongs to a different account than the local account ID.
func (a *ParsedARN) IsCrossAccount(localAccountID string) bool {
	return a.AccountID != "" && a.AccountID != localAccountID
}

// PartitionAWS is the standard AWS partition.
// PartitionAWSCN is the China AWS partition.
const (
	PartitionAWS   = "aws"    // The standard AWS partition
	PartitionAWSCN = "aws-cn" // The China AWS partition
)
