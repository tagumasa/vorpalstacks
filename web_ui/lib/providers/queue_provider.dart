import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/sqs.pbgrpc.dart';
import 'api_providers.dart' show sqsProvider;

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
    this.url = '',
    this.isFifo = false,
    this.availableMessages = 0,
    this.notVisibleMessages = 0,
    this.approximateNumberOfMessages = 0,
    required this.created,
  });
}

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
}

class SQSState {
  final List<Queue> queues;
  final bool isLoading;
  final String? error;
  final Queue? selectedQueue;
  final List<SqsMessage> items;

  const SQSState({
    this.queues = const [],
    this.isLoading = false,
    this.error,
    this.selectedQueue,
    this.items = const [],
  });

  SQSState copyWith({
    List<Queue>? queues,
    bool? isLoading,
    String? error,
    bool clearError = false,
    Queue? selectedQueue,
    List<SqsMessage>? items,
    bool clearSelectedQueue = false,
  }) {
    return SQSState(
      queues: queues ?? this.queues,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
      selectedQueue: clearSelectedQueue
          ? null
          : (selectedQueue ?? this.selectedQueue),
      items: items ?? this.items,
    );
  }
}

final queueListProvider = NotifierProvider<SQSNotifier, SQSState>(() {
  return SQSNotifier();
});

class SQSNotifier extends Notifier<SQSState> {
  @override
  SQSState build() => const SQSState();

  SQSServiceClient get _client => ref.read(sqsProvider);

  Future<void> loadQueues() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.listQueues(ListQueuesRequest());
      final queues = <Queue>[];
      for (final queueUrl in response.queueurls) {
        final name = queueUrl.split('/').last;
        final isFifo = name.endsWith('.fifo');
        final created = DateTime.now();
        queues.add(
          Queue(name: name, url: queueUrl, isFifo: isFifo, created: created),
        );
      }
      state = state.copyWith(queues: queues, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  void selectQueue(Queue queue) {
    state = state.copyWith(selectedQueue: queue, items: []);
  }

  Future<void> createQueue(String name, {bool isFifo = false}) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final request = CreateQueueRequest(queuename: name);
      if (isFifo) {
        request.attributes['FifoQueue'] = 'true';
      }
      final response = await _client.createQueue(request);
      final created = DateTime.now();
      final queue = Queue(
        name: name,
        url: response.queueurl.isNotEmpty
            ? response.queueurl
            : 'http://localhost:9090/queue/$name',
        isFifo: isFifo,
        created: created,
      );
      state = state.copyWith(
        queues: [...state.queues, queue],
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteQueue(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final existing = state.queues.where((q) => q.name == name).firstOrNull;
      final queueUrl = existing?.url ?? 'http://localhost:9090/queue/$name';
      await _client.deleteQueue(DeleteQueueRequest(queueurl: queueUrl));
      state = state.copyWith(
        queues: state.queues.where((q) => q.name != name).toList(),
        clearSelectedQueue: state.selectedQueue?.name == name,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> getMessages(String queueUrl) async {
    try {
      final response = await _client.receiveMessage(
        ReceiveMessageRequest(queueurl: queueUrl),
      );
      final messages = response.messages.map((msg) {
        final timestampStr = msg.attributes['SentTimestamp'] ?? '';
        return SqsMessage(
          messageId: msg.messageid,
          receiptHandle: msg.receipthandle,
          body: msg.body,
          timestamp: DateTime.tryParse(timestampStr) ?? DateTime.now(),
          approximateReceiveCount:
              int.tryParse(msg.attributes['ApproximateReceiveCount'] ?? '') ??
              0,
          messageGroupId: msg.attributes['MessageGroupId'],
          deduplicationId: msg.attributes['MessageDeduplicationId'],
        );
      }).toList();
      state = state.copyWith(items: messages);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}
