import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/waf.pbgrpc.dart';
import 'api_providers.dart';

class WAFWebACL {
  final String webAclId;
  final String name;

  WAFWebACL({required this.webAclId, required this.name});
}

class WAFState {
  final List<WAFWebACL> webAcls;
  final bool isLoading;
  final String? error;

  const WAFState({this.webAcls = const [], this.isLoading = false, this.error});

  WAFState copyWith({
    List<WAFWebACL>? webAcls,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return WAFState(
      webAcls: webAcls ?? this.webAcls,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class WAFListNotifier extends Notifier<WAFState> {
  @override
  WAFState build() => const WAFState();

  Future<void> loadWebAcls() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(wafProvider);
      final request = ListWebACLsRequest();
      final response = await client.listWebACLs(request);

      final webAcls = response.webacls.map((w) {
        return WAFWebACL(webAclId: w.webaclid, name: w.name);
      }).toList();

      state = state.copyWith(webAcls: webAcls, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final wafListProvider = NotifierProvider<WAFListNotifier, WAFState>(() {
  return WAFListNotifier();
});
