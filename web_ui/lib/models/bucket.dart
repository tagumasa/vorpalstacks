class S3Bucket {
  final String name;
  final String? region;
  final DateTime creationDate;
  final int objectCount;
  final int sizeBytes;

  S3Bucket({
    required this.name,
    this.region,
    required this.creationDate,
    this.objectCount = 0,
    this.sizeBytes = 0,
  });

  factory S3Bucket.fromJson(Map<String, dynamic> json) {
    return S3Bucket(
      name: json['name'] ?? '',
      region: json['region'],
      creationDate: json['creationDate'] != null
          ? DateTime.parse(json['creationDate'])
          : DateTime.now(),
      objectCount: json['objectCount'] ?? 0,
      sizeBytes: json['sizeBytes'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'region': region,
    'creationDate': creationDate.toIso8601String(),
    'objectCount': objectCount,
    'sizeBytes': sizeBytes,
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

class S3Object {
  final String key;
  final int size;
  final DateTime lastModified;
  final String? eTag;
  final String storageClass;

  S3Object({
    required this.key,
    this.size = 0,
    required this.lastModified,
    this.eTag,
    this.storageClass = 'STANDARD',
  });

  factory S3Object.fromJson(Map<String, dynamic> json) {
    return S3Object(
      key: json['key'] ?? '',
      size: json['size'] ?? 0,
      lastModified: json['lastModified'] != null
          ? DateTime.parse(json['lastModified'])
          : DateTime.now(),
      eTag: json['eTag'],
      storageClass: json['storageClass'] ?? 'STANDARD',
    );
  }

  String get formattedSize {
    if (size < 1024) return '$size B';
    if (size < 1024 * 1024) return '${(size / 1024).toStringAsFixed(1)} KB';
    if (size < 1024 * 1024 * 1024) {
      return '${(size / (1024 * 1024)).toStringAsFixed(1)} MB';
    }
    return '${(size / (1024 * 1024 * 1024)).toStringAsFixed(1)} GB';
  }

  bool get isFolder => key.endsWith('/');
}
