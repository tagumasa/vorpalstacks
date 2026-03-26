class ConfigEntry {
  final String key;
  final String value;
  final String type;
  final String source;
  final String description;
  final bool readOnly;
  final int updatedAt;
  final String envVar;
  final String category;

  const ConfigEntry({
    required this.key,
    required this.value,
    required this.type,
    required this.source,
    required this.description,
    required this.readOnly,
    required this.updatedAt,
    required this.envVar,
    required this.category,
  });

  factory ConfigEntry.fromJson(Map<String, dynamic> json) {
    return ConfigEntry(
      key: json['key'] as String? ?? '',
      value: json['value'] as String? ?? '',
      type: json['type'] as String? ?? '',
      source: json['source'] as String? ?? '',
      description: json['description'] as String? ?? '',
      readOnly: json['readOnly'] as bool? ?? false,
      updatedAt: json['updatedAt'] as int? ?? 0,
      envVar: json['envVar'] as String? ?? '',
      category: json['category'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'key': key,
    'value': value,
    'type': type,
    'source': source,
    'description': description,
    'readOnly': readOnly,
    'updatedAt': updatedAt,
    'envVar': envVar,
    'category': category,
  };
}

class ServiceInfo {
  final String name;
  final bool enabled;
  final int resourceCount;

  const ServiceInfo({
    required this.name,
    required this.enabled,
    required this.resourceCount,
  });

  factory ServiceInfo.fromJson(Map<String, dynamic> json) {
    return ServiceInfo(
      name: json['name'] as String? ?? '',
      enabled: json['enabled'] as bool? ?? false,
      resourceCount: json['resourceCount'] as int? ?? 0,
    );
  }
}

class ServiceStatus {
  final String name;
  final bool enabled;
  final String mode;
  final int resourceCount;
  final int lastRequestAt;

  const ServiceStatus({
    required this.name,
    required this.enabled,
    required this.mode,
    required this.resourceCount,
    required this.lastRequestAt,
  });

  factory ServiceStatus.fromJson(Map<String, dynamic> json) {
    return ServiceStatus(
      name: json['name'] as String? ?? '',
      enabled: json['enabled'] as bool? ?? false,
      mode: json['mode'] as String? ?? '',
      resourceCount: json['resourceCount'] as int? ?? 0,
      lastRequestAt: json['lastRequestAt'] as int? ?? 0,
    );
  }
}
