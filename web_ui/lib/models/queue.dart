class Queue {
  final String name;
  final String url;
  final bool isFifo;
  final int availableMessages;
  final int notVisibleMessages;
  final int approximateNumberOfMessages;
  final DateTime created;

  Queue({
    required this.name,
    required this.url,
    this.isFifo = false,
    this.availableMessages = 0,
    this.notVisibleMessages = 0,
    this.approximateNumberOfMessages = 0,
    required this.created,
  });

  factory Queue.fromJson(Map<String, dynamic> json) {
    return Queue(
      name: json['name'] ?? '',
      url: json['url'] ?? '',
      isFifo: json['isFifo'] ?? false,
      availableMessages: json['availableMessages'] ?? 0,
      notVisibleMessages: json['notVisibleMessages'] ?? 0,
      approximateNumberOfMessages: json['approximateNumberOfMessages'] ?? 0,
      created: json['created'] != null
          ? DateTime.parse(json['created'])
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'url': url,
    'isFifo': isFifo,
    'availableMessages': availableMessages,
    'notVisibleMessages': notVisibleMessages,
    'approximateNumberOfMessages': approximateNumberOfMessages,
    'created': created.toIso8601String(),
  };

  Queue copyWith({
    String? name,
    String? url,
    bool? isFifo,
    int? availableMessages,
    int? notVisibleMessages,
    int? approximateNumberOfMessages,
    DateTime? created,
  }) {
    return Queue(
      name: name ?? this.name,
      url: url ?? this.url,
      isFifo: isFifo ?? this.isFifo,
      availableMessages: availableMessages ?? this.availableMessages,
      notVisibleMessages: notVisibleMessages ?? this.notVisibleMessages,
      approximateNumberOfMessages:
          approximateNumberOfMessages ?? this.approximateNumberOfMessages,
      created: created ?? this.created,
    );
  }
}
