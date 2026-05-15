package vmysql

// pebble_index.go: Secondary index integration with go-mysql-server.
// Not yet implemented — tables use primary-key-only access for Phase 0c.
// go-mysql-server handles WHERE clause filtering in-memory via its analyzer.
// rdbengine already supports index operations via index.go; the adapter
// layer (sql.IndexableTable → rdbengine index) will be added when needed.
