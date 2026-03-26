class HostedZone {
  final String id;
  final String name;
  final String arn;
  final int recordSetCount;
  final String? comment;
  final DateTime? creationTime;

  HostedZone({
    required this.id,
    required this.name,
    required this.arn,
    this.recordSetCount = 0,
    this.comment,
    this.creationTime,
  });

  factory HostedZone.fromJson(Map<String, dynamic> json) {
    return HostedZone(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      arn: json['arn'] ?? '',
      recordSetCount: json['recordSetCount'] ?? 0,
      comment: json['comment'],
      creationTime: json['creationTime'] != null
          ? DateTime.parse(json['creationTime'])
          : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'arn': arn,
    'recordSetCount': recordSetCount,
    'comment': comment,
    'creationTime': creationTime?.toIso8601String(),
  };
}

class ResourceRecord {
  final String name;
  final String type;
  final String value;
  final int ttl;
  final String? routingPolicy;

  ResourceRecord({
    required this.name,
    required this.type,
    required this.value,
    this.ttl = 300,
    this.routingPolicy,
  });

  factory ResourceRecord.fromJson(Map<String, dynamic> json) {
    return ResourceRecord(
      name: json['name'] ?? '',
      type: json['type'] ?? 'A',
      value: json['value'] ?? '',
      ttl: json['ttl'] ?? 300,
      routingPolicy: json['routingPolicy'],
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'type': type,
    'value': value,
    'ttl': ttl,
    'routingPolicy': routingPolicy,
  };
}
