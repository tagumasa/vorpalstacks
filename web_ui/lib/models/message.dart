class SqsMessage {
  final String messageId;
  final String receiptHandle;
  final String body;
  final DateTime timestamp;
  final int approximateReceiveCount;
  final String? deduplicationId;
  final String? messageGroupId;

  SqsMessage({
    required this.messageId,
    required this.receiptHandle,
    required this.body,
    required this.timestamp,
    this.approximateReceiveCount = 0,
    this.deduplicationId,
    this.messageGroupId,
  });

  factory SqsMessage.fromJson(Map<String, dynamic> json) {
    return SqsMessage(
      messageId: json['messageId'] ?? '',
      receiptHandle: json['receiptHandle'] ?? '',
      body: json['body'] ?? '',
      timestamp: json['timestamp'] != null
          ? DateTime.parse(json['timestamp'])
          : DateTime.now(),
      approximateReceiveCount: json['approximateReceiveCount'] ?? 0,
      deduplicationId: json['deduplicationId'],
      messageGroupId: json['messageGroupId'],
    );
  }

  Map<String, dynamic> toJson() => {
    'messageId': messageId,
    'receiptHandle': receiptHandle,
    'body': body,
    'timestamp': timestamp.toIso8601String(),
    'approximateReceiveCount': approximateReceiveCount,
    'deduplicationId': deduplicationId,
    'messageGroupId': messageGroupId,
  };
}
