import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/sns.pbgrpc.dart';
import '../src/generated/sns.pb.dart';
import 'api_providers.dart';

class SNSTopic {
  final String topicArn;
  final String name;
  final String displayName;
  final int subscriptionsConfirmed;
  final int subscriptionsPending;
  final DateTime created;

  SNSTopic({
    required this.topicArn,
    required this.name,
    this.displayName = '',
    this.subscriptionsConfirmed = 0,
    this.subscriptionsPending = 0,
    required this.created,
  });
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
    this.status = '',
    this.created,
  });
}

class SNSState {
  final List<SNSTopic> topics;
  final List<SNSSubscription> subscriptions;
  final bool isLoading;
  final String? error;

  const SNSState({
    this.topics = const [],
    this.subscriptions = const [],
    this.isLoading = false,
    this.error,
  });

  SNSState copyWith({
    List<SNSTopic>? topics,
    List<SNSSubscription>? subscriptions,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return SNSState(
      topics: topics ?? this.topics,
      subscriptions: subscriptions ?? this.subscriptions,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class SNSListNotifier extends Notifier<SNSState> {
  @override
  SNSState build() => const SNSState();

  Future<void> loadTopics() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(snsProvider);
      final request = ListTopicsInput();
      final response = await client.listTopics(request);

      final topics = response.topics.map((t) {
        final arn = t.topicarn;
        final name = arn.split(':').last;
        return SNSTopic(topicArn: arn, name: name, created: DateTime.now());
      });

      state = state.copyWith(topics: topics.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final snsListProvider = NotifierProvider<SNSListNotifier, SNSState>(() {
  return SNSListNotifier();
});
