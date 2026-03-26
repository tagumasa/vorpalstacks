class UserPool {
  final String id;
  final String name;
  final String status;
  final int usersCount;
  final DateTime created;

  UserPool({
    required this.id,
    required this.name,
    this.status = 'Active',
    this.usersCount = 0,
    required this.created,
  });

  factory UserPool.fromJson(Map<String, dynamic> json) {
    return UserPool(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      status: json['status'] ?? 'Active',
      usersCount: json['usersCount'] ?? 0,
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'status': status,
    'usersCount': usersCount,
    'created': created.toIso8601String(),
  };
}

class CognitoUser {
  final String username;
  final String? email;
  final String status;
  final DateTime created;
  final DateTime? lastModified;
  final List<String> groups;

  CognitoUser({
    required this.username,
    this.email,
    this.status = 'UNCONFIRMED',
    required this.created,
    this.lastModified,
    this.groups = const [],
  });

  factory CognitoUser.fromJson(Map<String, dynamic> json) {
    return CognitoUser(
      username: json['username'] ?? '',
      email: json['email'],
      status: json['status'] ?? 'UNCONFIRMED',
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
      lastModified: json['lastModified'] != null
          ? DateTime.parse(json['lastModified'])
          : null,
      groups: List<String>.from(json['groups'] ?? []),
    );
  }

  Map<String, dynamic> toJson() => {
    'username': username,
    'email': email,
    'status': status,
    'created': created.toIso8601String(),
    'lastModified': lastModified?.toIso8601String(),
    'groups': groups,
  };
}
