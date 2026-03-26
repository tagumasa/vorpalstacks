import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../providers/api_providers.dart';
import '../generated/admin_config.pb.dart';
export '../generated/admin_config.pb.dart' show ConfigEntry;

class ConfigState {
  final List<ConfigEntry> entries;
  final bool isLoading;
  final String? error;
  final String? selectedCategory;

  const ConfigState({
    this.entries = const [],
    this.isLoading = false,
    this.error,
    this.selectedCategory,
  });

  ConfigState copyWith({
    List<ConfigEntry>? entries,
    bool? isLoading,
    String? error,
    String? selectedCategory,
  }) {
    return ConfigState(
      entries: entries ?? this.entries,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      selectedCategory: selectedCategory ?? this.selectedCategory,
    );
  }

  List<ConfigEntry> get entriesByCategory {
    if (selectedCategory == null || selectedCategory!.isEmpty) {
      return entries;
    }
    return entries.where((e) => e.category == selectedCategory).toList();
  }

  List<String> get categories {
    final cats = entries.map((e) => e.category).toSet().toList();
    cats.sort();
    return cats;
  }
}

class ConfigNotifier extends Notifier<ConfigState> {
  @override
  ConfigState build() {
    Future.microtask(() => loadConfig());
    return const ConfigState();
  }

  Future<void> loadConfig({String? category}) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final client = ref.read(adminConfigProvider);
      final response = await client.listConfig(
        ListConfigRequest(category: category ?? ''),
      );
      state = state.copyWith(
        entries: response.entries,
        isLoading: false,
        selectedCategory: category,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> updateConfig(String key, String value) async {
    try {
      final client = ref.read(adminConfigProvider);
      final updated = await client.updateConfig(
        UpdateConfigRequest(key: key, value: value),
      );
      final newEntries = state.entries.map((e) {
        return e.key == key ? updated : e;
      }).toList();
      state = state.copyWith(entries: newEntries);
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  Future<bool> resetConfig(String key) async {
    try {
      final client = ref.read(adminConfigProvider);
      final reset = await client.resetConfig(ResetConfigRequest(key: key));
      final newEntries = state.entries.map((e) {
        return e.key == key ? reset : e;
      }).toList();
      state = state.copyWith(entries: newEntries);
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  void filterByCategory(String? category) {
    state = state.copyWith(selectedCategory: category);
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

final configProvider = NotifierProvider<ConfigNotifier, ConfigState>(() {
  return ConfigNotifier();
});

final categoriesProvider = Provider<List<String>>((ref) {
  return ref.watch(configProvider).categories;
});
