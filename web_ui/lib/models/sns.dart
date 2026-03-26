class SNSTopic {
  final String topicArn;
  final String name;
  final int subscriptionsConfirmed;
  final int subscriptionsPending;
  final String? displayName;
  final DateTime created;

  SNSTopic({
    required this.topicArn,
    required this.name,
    this.subscriptionsConfirmed = 0,
    this.subscriptionsPending = 0,
    this.displayName,
    required this.created,
  });

  factory SNSTopic.fromJson(Map<String, dynamic> json) {
    return SNSTopic(
      topicArn: json['topicArn'] ?? '',
      name: json['name'] ?? '',
      subscriptionsConfirmed: json['subscriptionsConfirmed'] ?? 0,
      subscriptionsPending: json['subscriptionsPending'] ?? 0,
      displayName: json['displayName'],
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'topicArn': topicArn,
    'name': name,
    'subscriptionsConfirmed': subscriptionsConfirmed,
    'subscriptionsPending': subscriptionsPending,
    'displayName': displayName,
    'created': created.toIso8601String(),
  };
}

class SNSSubscription {
  final String subscriptionArn;
  final String topicArn;
  final String protocol;
  final String endpoint;
  final String status;
  final DateTime? created;

  SNSSubscription({
    required this.subscriptionArn,
    required this.topicArn,
    required this.protocol,
    required this.endpoint,
    required this.status,
    this.created,
  });

  factory SNSSubscription.fromJson(Map<String, dynamic> json) {
    return SNSSubscription(
      subscriptionArn: json['subscriptionArn'] ?? '',
      topicArn: json['topicArn'] ?? '',
      protocol: json['protocol'] ?? '',
      endpoint: json['endpoint'] ?? '',
      status: json['status'] ?? 'pending',
      created: json['created'] != null ? DateTime.parse(json['created']) : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'subscriptionArn': subscriptionArn,
    'topicArn': topicArn,
    'protocol': protocol,
    'endpoint': endpoint,
    'status': status,
    'created': created?.toIso8601String(),
  };
}
