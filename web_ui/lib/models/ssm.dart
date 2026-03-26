class SSMParameter {
  final String name;
  final String arn;
  final String type;
  final String value;
  final int version;
  final DateTime? lastModifiedDate;
  final String? description;
  final String? dataType;

  SSMParameter({
    required this.name,
    required this.arn,
    required this.type,
    required this.value,
    this.version = 1,
    this.lastModifiedDate,
    this.description,
    this.dataType,
  });

  factory SSMParameter.fromJson(Map<String, dynamic> json) {
    return SSMParameter(
      name: json['name'] ?? '',
      arn: json['arn'] ?? '',
      type: json['type'] ?? 'String',
      value: json['value'] ?? '',
      version: json['version'] ?? 1,
      lastModifiedDate: json['lastModifiedDate'] != null
          ? DateTime.parse(json['lastModifiedDate'])
          : null,
      description: json['description'],
      dataType: json['dataType'],
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'arn': arn,
    'type': type,
    'value': value,
    'version': version,
    'lastModifiedDate': lastModifiedDate?.toIso8601String(),
    'description': description,
    'dataType': dataType,
  };
}
