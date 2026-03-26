import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/wafv2.pbgrpc.dart';
import 'api_providers.dart';

class WAFv2WebACL {
  final String id;
  final String name;
  final String arn;
  final String description;

  WAFv2WebACL({
    required this.id,
    required this.name,
    required this.arn,
    required this.description,
  });
}

class WAFv2State {
  final List<WAFv2WebACL> webAcls;
  final bool isLoading;
  final String? error;

  const WAFv2State({
    this.webAcls = const [],
    this.isLoading = false,
    this.error,
  });

  WAFv2State copyWith({
    List<WAFv2WebACL>? webAcls,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return WAFv2State(
      webAcls: webAcls ?? this.webAcls,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class WAFv2ListNotifier extends Notifier<WAFv2State> {
  @override
  WAFv2State build() => const WAFv2State();

  Future<void> loadWebAcls() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(wafv2Provider);
      final request = ListWebACLsRequest();
      final response = await client.listWebACLs(request);

      final webAcls = response.webacls.map((w) {
        return WAFv2WebACL(
          id: w.id,
          name: w.name,
          arn: w.arn,
          description: w.description,
        );
      }).toList();

      state = state.copyWith(webAcls: webAcls, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final wafv2ListProvider = NotifierProvider<WAFv2ListNotifier, WAFv2State>(() {
  return WAFv2ListNotifier();
});
