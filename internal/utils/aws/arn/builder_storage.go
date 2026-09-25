// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// S3Builder provides methods for constructing S3 (Simple Storage Service) ARNs.
type S3Builder struct{ *ARNBuilder }

// S3 returns an S3Builder for constructing S3 ARNs.
func (b *ARNBuilder) S3() *S3Builder { return &S3Builder{b} }

// Bucket constructs an ARN for an S3 bucket.
func (b *S3Builder) Bucket(name string) string { return b.BuildGlobal("s3", name) }

// Object constructs an ARN for an S3 object.
func (b *S3Builder) Object(bucket, key string) string { return b.BuildGlobal("s3", bucket+"/"+key) }

// ObjectVersion constructs an ARN for a specific version of an S3 object.
func (b *S3Builder) ObjectVersion(bucket, key, versionId string) string {
	return b.BuildGlobal("s3", bucket+"/"+key+"?versionId="+versionId)
}

// DynamoDBBuilder provides methods for constructing DynamoDB ARNs.
type DynamoDBBuilder struct{ *ARNBuilder }

// DynamoDB returns a DynamoDBBuilder for constructing DynamoDB ARNs.
func (b *ARNBuilder) DynamoDB() *DynamoDBBuilder { return &DynamoDBBuilder{b} }

// Table constructs an ARN for a DynamoDB table.
func (b *DynamoDBBuilder) Table(name string) string { return b.Build("dynamodb", "table/"+name) }

// TableStream constructs an ARN for a DynamoDB table stream.
func (b *DynamoDBBuilder) TableStream(name, label string) string {
	return b.Build("dynamodb", "table/"+name+"/stream/"+label)
}

// Index constructs an ARN for a DynamoDB global secondary index.
func (b *DynamoDBBuilder) Index(table, index string) string {
	return b.Build("dynamodb", "table/"+table+"/index/"+index)
}

// Backup constructs an ARN for a DynamoDB table backup. The identity
// segment after /backup/ is the backup's generated id, never its name —
// the documented example shows table/Music/backup/01489602797149-73d8d5bc
// for the backup named MusicBackup — so same-named backups of one table
// carry distinct ARNs.
func (b *DynamoDBBuilder) Backup(table, id string) string {
	return b.Build("dynamodb", "table/"+table+"/backup/"+id)
}

// GlobalTable constructs an ARN for a DynamoDB global table.
func (b *DynamoDBBuilder) GlobalTable(name string) string {
	return b.Build("dynamodb", "globaltable/"+name)
}

// Export constructs an ARN for a DynamoDB table export. The resource
// path embeds the exporting table — the documented format is
// arn:aws:dynamodb:<region>:<account>:table/<TableName>/export/<id> —
// because the ID segment alone distinguishes concurrent exports of
// the same table, it carries the creation's unique identity.
func (b *DynamoDBBuilder) Export(tableArn, exportId string) string {
	_, _, _, _, resource := SplitARN(tableArn)
	return b.Build("dynamodb", resource+"/export/"+exportId)
}

// Import constructs an ARN for a DynamoDB table import, embedding the
// importing table's resource path the same way Export does for the
// export family.
func (b *DynamoDBBuilder) Import(tableArn, importId string) string {
	_, _, _, _, resource := SplitARN(tableArn)
	return b.Build("dynamodb", resource+"/import/"+importId)
}

// ParseTableName extracts the table name from a DynamoDB table ARN.
func (b *DynamoDBBuilder) ParseTableName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if len(resource) > 6 && resource[:6] == "table/" {
		rest := resource[6:]
		if idx := strings.Index(rest, "/"); idx != -1 {
			return rest[:idx]
		}
		return rest
	}
	return ""
}
