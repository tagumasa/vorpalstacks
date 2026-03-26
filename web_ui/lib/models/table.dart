class DynamoTable {
  final String name;
  final String? arn;
  final int itemCount;
  final int sizeBytes;
  final String status;
  final String? partitionKey;
  final String? sortKey;
  final DateTime created;

  DynamoTable({
    required this.name,
    this.arn,
    this.itemCount = 0,
    this.sizeBytes = 0,
    this.status = 'ACTIVE',
    this.partitionKey,
    this.sortKey,
    required this.created,
  });

  factory DynamoTable.fromJson(Map<String, dynamic> json) {
    return DynamoTable(
      name: json['name'] ?? '',
      arn: json['arn'],
      itemCount: json['itemCount'] ?? 0,
      sizeBytes: json['sizeBytes'] ?? 0,
      status: json['status'] ?? 'ACTIVE',
      partitionKey: json['partitionKey'],
      sortKey: json['sortKey'],
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'arn': arn,
    'itemCount': itemCount,
    'sizeBytes': sizeBytes,
    'status': status,
    'partitionKey': partitionKey,
    'sortKey': sortKey,
    'created': created.toIso8601String(),
  };

  String get formattedSize {
    if (sizeBytes < 1024) return '$sizeBytes B';
    if (sizeBytes < 1024 * 1024) {
      return '${(sizeBytes / 1024).toStringAsFixed(1)} KB';
    }
    if (sizeBytes < 1024 * 1024 * 1024) {
      return '${(sizeBytes / (1024 * 1024)).toStringAsFixed(1)} MB';
    }
    return '${(sizeBytes / (1024 * 1024 * 1024)).toStringAsFixed(1)} GB';
  }
}

class DynamoItem {
  final Map<String, dynamic> data;
  final String? id;

  DynamoItem({required this.data, this.id});

  factory DynamoItem.fromJson(Map<String, dynamic> json) {
    return DynamoItem(data: json);
  }
}
