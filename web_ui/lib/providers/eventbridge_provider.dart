import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/events.pbgrpc.dart';
import 'api_providers.dart';

class EventRule {
  final String name;
  final String arn;
  final String state;
  final String description;
  final String eventPattern;
  final String scheduleExpression;

  EventRule({
    required this.name,
    this.arn = '',
    this.state = 'ENABLED',
    this.description = '',
    this.eventPattern = '',
    this.scheduleExpression = '',
  });
}

class EventBridgeState {
  final List<EventRule> rules;
  final bool isLoading;
  final String? error;

  const EventBridgeState({
    this.rules = const [],
    this.isLoading = false,
    this.error,
  });

  EventBridgeState copyWith({
    List<EventRule>? rules,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return EventBridgeState(
      rules: rules ?? this.rules,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class EventBridgeListNotifier extends Notifier<EventBridgeState> {
  @override
  EventBridgeState build() => const EventBridgeState();

  Future<void> loadRules() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cloudWatchEventsProvider);
      final request = ListRulesRequest();
      final response = await client.listRules(request);

      final rules = response.rules.map((r) {
        return EventRule(
          name: r.name,
          arn: r.arn,
          state: r.state.name.replaceAll('RULE_STATE_', ''),
          description: r.description,
          eventPattern: r.eventpattern,
          scheduleExpression: r.scheduleexpression,
        );
      });

      state = state.copyWith(rules: rules.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createRule({
    required String name,
    String? scheduleExpression,
    String? eventPattern,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(cloudWatchEventsProvider);
      final request = PutRuleRequest(
        name: name,
        scheduleexpression: scheduleExpression,
        eventpattern: eventPattern,
      );
      await client.putRule(request);
      await loadRules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteRule(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(cloudWatchEventsProvider);
      final request = DeleteRuleRequest(name: name);
      await client.deleteRule(request);
      await loadRules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> enableRule(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(cloudWatchEventsProvider);
      final request = EnableRuleRequest(name: name);
      await client.enableRule(request);
      await loadRules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> disableRule(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(cloudWatchEventsProvider);
      final request = DisableRuleRequest(name: name);
      await client.disableRule(request);
      await loadRules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final eventBridgeListProvider =
    NotifierProvider<EventBridgeListNotifier, EventBridgeState>(() {
      return EventBridgeListNotifier();
    });
