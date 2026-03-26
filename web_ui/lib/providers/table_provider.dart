import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/dynamodb.pbgrpc.dart';
import 'api_providers.dart' show dynamoDbProvider;

class DynamoTable {
  final String name;
  final String arn;
  final int itemCount;
  final int sizeBytes;
  final String status;
  final String partitionKey;
  final String? sortKey;
  final DateTime created;

  DynamoTable({
    required this.name,
    this.arn = '',
    this.itemCount = 0,
    this.sizeBytes = 0,
    this.status = 'ACTIVE',
    this.partitionKey = '',
    this.sortKey,
    required this.created,
  });
}

class DynamoItem {
  final Map<String, dynamic> data;

  DynamoItem({required this.data});
}

class DynamoDBState {
  final List<DynamoTable> tables;
  final DynamoTable? selectedTable;
  final List<DynamoItem> items;
  final bool isLoading;
  final String? error;

  const DynamoDBState({
    this.tables = const [],
    this.selectedTable,
    this.items = const [],
    this.isLoading = false,
    this.error,
  });

  DynamoDBState copyWith({
    List<DynamoTable>? tables,
    DynamoTable? selectedTable,
    bool clearSelectedTable = false,
    List<DynamoItem>? items,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return DynamoDBState(
      tables: tables ?? this.tables,
      selectedTable: clearSelectedTable
          ? null
          : (selectedTable ?? this.selectedTable),
      items: items ?? this.items,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class DynamoDBNotifier extends Notifier<DynamoDBState> {
  @override
  DynamoDBState build() => const DynamoDBState();

  DynamoDBServiceClient get _client => ref.read(dynamoDbProvider);

  Future<void> loadTables() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.listTables(ListTablesInput());
      final tables = <DynamoTable>[];
      for (final tableName in response.tablenames) {
        try {
          final desc = await _client.describeTable(
            DescribeTableInput(tablename: tableName),
          );
          final table = desc.table;
          DateTime created;
          try {
            created = table.creationdatetime.isNotEmpty
                ? DateTime.parse(table.creationdatetime)
                : DateTime.now();
          } catch (_) {
            created = DateTime.now();
          }
          String partitionKey = '';
          String? sortKey;
          for (final key in table.keyschema) {
            if (key.keytype == KeyType.KEY_TYPE_HASH) {
              partitionKey = key.attributename;
            } else if (key.keytype == KeyType.KEY_TYPE_RANGE) {
              sortKey = key.attributename;
            }
          }
          tables.add(
            DynamoTable(
              name: table.tablename,
              arn: table.tablearn,
              itemCount: table.itemcount.toInt(),
              sizeBytes: table.tablesizebytes.toInt(),
              status: table.tablestatus.name.replaceAll('TABLE_STATUS_', ''),
              partitionKey: partitionKey,
              sortKey: sortKey,
              created: created,
            ),
          );
        } catch (_) {}
      }
      state = state.copyWith(tables: tables, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  void selectTable(DynamoTable table) {
    state = state.copyWith(selectedTable: table, items: []);
  }

  Future<void> createTable(
    String name,
    String partitionKey, {
    String partitionKeyType = 'S',
    String? sortKey,
    String? sortKeyType,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final pkType = _parseAttributeType(partitionKeyType);
      final keySchema = <KeySchemaElement>[
        KeySchemaElement(
          attributename: partitionKey,
          keytype: KeyType.KEY_TYPE_HASH,
        ),
      ];
      if (sortKey != null) {
        keySchema.add(
          KeySchemaElement(
            attributename: sortKey,
            keytype: KeyType.KEY_TYPE_RANGE,
          ),
        );
      }
      final attributeDefs = <AttributeDefinition>[
        AttributeDefinition(attributename: partitionKey, attributetype: pkType),
      ];
      if (sortKey != null) {
        final skType = _parseAttributeType(sortKeyType ?? 'S');
        attributeDefs.add(
          AttributeDefinition(attributename: sortKey, attributetype: skType),
        );
      }
      await _client.createTable(
        CreateTableInput(
          tablename: name,
          keyschema: keySchema,
          attributedefinitions: attributeDefs,
          billingmode: BillingMode.BILLING_MODE_PAY_PER_REQUEST,
        ),
      );
      await loadTables();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteTable(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.deleteTable(DeleteTableInput(tablename: name));
      state = state.copyWith(
        tables: state.tables.where((t) => t.name != name).toList(),
        clearSelectedTable: state.selectedTable?.name == name,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> scanTable(String tableName) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.scan(
        ScanInput(tablename: tableName, limit: 100),
      );
      final items = <DynamoItem>[];
      response.items.forEach((key, value) {
        final data = <String, dynamic>{};
        data[key] = _parseAttributeValue(value);
        items.add(DynamoItem(data: data));
      });
      state = state.copyWith(items: items, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  dynamic _parseAttributeValue(AttributeValue value) {
    if (value.s.isNotEmpty) return value.s;
    if (value.n.isNotEmpty) return num.tryParse(value.n) ?? value.n;
    if (value.b.isNotEmpty) return value.b;
    if (value.hasBool_307144766()) return value.bool_307144766;
    if (value.hasNull_426761765()) return null;
    return null;
  }

  ScalarAttributeType _parseAttributeType(String type) {
    switch (type.toUpperCase()) {
      case 'N':
        return ScalarAttributeType.SCALAR_ATTRIBUTE_TYPE_N;
      case 'B':
        return ScalarAttributeType.SCALAR_ATTRIBUTE_TYPE_B;
      default:
        return ScalarAttributeType.SCALAR_ATTRIBUTE_TYPE_S;
    }
  }
}

final tableListProvider = NotifierProvider<DynamoDBNotifier, DynamoDBState>(() {
  return DynamoDBNotifier();
});
