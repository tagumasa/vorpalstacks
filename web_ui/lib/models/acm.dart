class Certificate {
  final String certificateArn;
  final String domainName;
  final String status;
  final String type;
  final List<String> subjectAlternativeNames;
  final DateTime? notBefore;
  final DateTime? notAfter;
  final DateTime? createdAt;
  final bool inUse;
  final String? keyAlgorithm;

  Certificate({
    required this.certificateArn,
    required this.domainName,
    required this.status,
    this.type = 'AMAZON_ISSUED',
    this.subjectAlternativeNames = const [],
    this.notBefore,
    this.notAfter,
    this.createdAt,
    this.inUse = false,
    this.keyAlgorithm,
  });

  factory Certificate.fromJson(Map<String, dynamic> json) {
    return Certificate(
      certificateArn: json['certificateArn'] ?? '',
      domainName: json['domainName'] ?? '',
      status: json['status'] ?? 'PENDING_VALIDATION',
      type: json['type'] ?? 'AMAZON_ISSUED',
      subjectAlternativeNames: json['subjectAlternativeNames'] != null
          ? List<String>.from(json['subjectAlternativeNames'])
          : [],
      notBefore: json['notBefore'] != null
          ? DateTime.parse(json['notBefore'])
          : null,
      notAfter: json['notAfter'] != null
          ? DateTime.parse(json['notAfter'])
          : null,
      createdAt: json['createdAt'] != null
          ? DateTime.parse(json['createdAt'])
          : null,
      inUse: json['inUse'] ?? false,
      keyAlgorithm: json['keyAlgorithm'],
    );
  }

  Map<String, dynamic> toJson() => {
    'certificateArn': certificateArn,
    'domainName': domainName,
    'status': status,
    'type': type,
    'subjectAlternativeNames': subjectAlternativeNames,
    'notBefore': notBefore?.toIso8601String(),
    'notAfter': notAfter?.toIso8601String(),
    'createdAt': createdAt?.toIso8601String(),
    'inUse': inUse,
    'keyAlgorithm': keyAlgorithm,
  };
}
