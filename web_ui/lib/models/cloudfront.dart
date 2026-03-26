class Distribution {
  final String id;
  final String arn;
  final String status;
  final bool enabled;
  final String? domainName;
  final List<String> origins;
  final String? defaultCacheBehavior;
  final String priceClass;
  final DateTime? lastModifiedTime;
  final bool isIpv6Enabled;
  final String? comment;

  Distribution({
    required this.id,
    required this.arn,
    required this.status,
    this.enabled = false,
    this.domainName,
    this.origins = const [],
    this.defaultCacheBehavior,
    this.priceClass = 'PriceClass_100',
    this.lastModifiedTime,
    this.isIpv6Enabled = true,
    this.comment,
  });

  factory Distribution.fromJson(Map<String, dynamic> json) {
    return Distribution(
      id: json['id'] ?? '',
      arn: json['arn'] ?? '',
      status: json['status'] ?? 'InProgress',
      enabled: json['enabled'] ?? false,
      domainName: json['domainName'],
      origins: json['origins'] != null
          ? List<String>.from(json['origins'])
          : [],
      defaultCacheBehavior: json['defaultCacheBehavior'],
      priceClass: json['priceClass'] ?? 'PriceClass_100',
      lastModifiedTime: json['lastModifiedTime'] != null
          ? DateTime.parse(json['lastModifiedTime'])
          : null,
      isIpv6Enabled: json['isIpv6Enabled'] ?? true,
      comment: json['comment'],
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'arn': arn,
    'status': status,
    'enabled': enabled,
    'domainName': domainName,
    'origins': origins,
    'defaultCacheBehavior': defaultCacheBehavior,
    'priceClass': priceClass,
    'lastModifiedTime': lastModifiedTime?.toIso8601String(),
    'isIpv6Enabled': isIpv6Enabled,
    'comment': comment,
  };
}
