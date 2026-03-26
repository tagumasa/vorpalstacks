class Secret {
  final String name;
  final String arn;
  final String? description;
  final String? kmsKeyId;
  final DateTime? createdDate;
  final DateTime? lastChangedDate;
  final DateTime? lastAccessedDate;
  final int versionIdsToStagesCount;
  final bool replicationEnabled;

  Secret({
    required this.name,
    required this.arn,
    this.description,
    this.kmsKeyId,
    this.createdDate,
    this.lastChangedDate,
    this.lastAccessedDate,
    this.versionIdsToStagesCount = 0,
    this.replicationEnabled = false,
  });

  factory Secret.fromJson(Map<String, dynamic> json) {
    return Secret(
      name: json['name'] ?? '',
      arn: json['arn'] ?? '',
      description: json['description'],
      kmsKeyId: json['kmsKeyId'],
      createdDate: json['createdDate'] != null
          ? DateTime.parse(json['createdDate'])
          : null,
      lastChangedDate: json['lastChangedDate'] != null
          ? DateTime.parse(json['lastChangedDate'])
          : null,
      lastAccessedDate: json['lastAccessedDate'] != null
          ? DateTime.parse(json['lastAccessedDate'])
          : null,
      versionIdsToStagesCount: json['versionIdsToStagesCount'] ?? 0,
      replicationEnabled: json['replicationEnabled'] ?? false,
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'arn': arn,
    'description': description,
    'kmsKeyId': kmsKeyId,
    'createdDate': createdDate?.toIso8601String(),
    'lastChangedDate': lastChangedDate?.toIso8601String(),
    'lastAccessedDate': lastAccessedDate?.toIso8601String(),
    'versionIdsToStagesCount': versionIdsToStagesCount,
    'replicationEnabled': replicationEnabled,
  };
}
