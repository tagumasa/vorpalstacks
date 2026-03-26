class IAMUser {
  final String arn;
  final String userName;
  final String userId;
  final DateTime created;
  final String status;
  final List<String> groups;

  IAMUser({
    required this.arn,
    required this.userName,
    required this.userId,
    required this.created,
    this.status = 'Active',
    this.groups = const [],
  });

  factory IAMUser.fromJson(Map<String, dynamic> json) {
    return IAMUser(
      arn: json['arn'] ?? '',
      userName: json['userName'] ?? '',
      userId: json['userId'] ?? '',
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
      status: json['status'] ?? 'Active',
      groups: List<String>.from(json['groups'] ?? []),
    );
  }

  Map<String, dynamic> toJson() => {
    'arn': arn,
    'userName': userName,
    'userId': userId,
    'created': created.toIso8601String(),
    'status': status,
    'groups': groups,
  };
}

class IAMRole {
  final String arn;
  final String roleName;
  final String roleId;
  final DateTime created;
  final String? path;

  IAMRole({
    required this.arn,
    required this.roleName,
    required this.roleId,
    required this.created,
    this.path,
  });

  factory IAMRole.fromJson(Map<String, dynamic> json) {
    return IAMRole(
      arn: json['arn'] ?? '',
      roleName: json['roleName'] ?? '',
      roleId: json['roleId'] ?? '',
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
      path: json['path'],
    );
  }

  Map<String, dynamic> toJson() => {
    'arn': arn,
    'roleName': roleName,
    'roleId': roleId,
    'created': created.toIso8601String(),
    'path': path,
  };
}

class IAMPolicy {
  final String arn;
  final String policyName;
  final String type;
  final int attachedCount;
  final String? description;

  IAMPolicy({
    required this.arn,
    required this.policyName,
    this.type = 'Managed',
    this.attachedCount = 0,
    this.description,
  });

  factory IAMPolicy.fromJson(Map<String, dynamic> json) {
    return IAMPolicy(
      arn: json['arn'] ?? '',
      policyName: json['policyName'] ?? '',
      type: json['type'] ?? 'Managed',
      attachedCount: json['attachedCount'] ?? 0,
      description: json['description'],
    );
  }

  Map<String, dynamic> toJson() => {
    'arn': arn,
    'policyName': policyName,
    'type': type,
    'attachedCount': attachedCount,
    'description': description,
  };
}
