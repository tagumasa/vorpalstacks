import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/cloudtrail.pbgrpc.dart';
import 'api_providers.dart';

class CloudTrailTrail {
  final String trailARN;
  final String name;
  final String homeRegion;
  final DateTime created;

  CloudTrailTrail({
    required this.trailARN,
    required this.name,
    this.homeRegion = '',
    required this.created,
  });
}

class CloudTrailEventDataStore {
  final String eventDataStoreARN;
  final String name;
  final String status;
  final DateTime created;

  CloudTrailEventDataStore({
    required this.eventDataStoreARN,
    required this.name,
    required this.status,
    required this.created,
  });
}

class CloudTrailState {
  final List<CloudTrailTrail> trails;
  final List<CloudTrailEventDataStore> eventDataStores;
  final bool isLoading;
  final String? error;

  const CloudTrailState({
    this.trails = const [],
    this.eventDataStores = const [],
    this.isLoading = false,
    this.error,
  });

  CloudTrailState copyWith({
    List<CloudTrailTrail>? trails,
    List<CloudTrailEventDataStore>? eventDataStores,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return CloudTrailState(
      trails: trails ?? this.trails,
      eventDataStores: eventDataStores ?? this.eventDataStores,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class CloudTrailListNotifier extends Notifier<CloudTrailState> {
  @override
  CloudTrailState build() => const CloudTrailState();

  Future<void> loadTrails() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cloudTrailProvider);
      final request = ListTrailsRequest();
      final response = await client.listTrails(request);

      final trails = response.trails.map((t) {
        return CloudTrailTrail(
          trailARN: t.trailarn,
          name: t.name,
          homeRegion: t.homeregion,
          created: DateTime.now(),
        );
      }).toList();

      state = state.copyWith(trails: trails, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> loadEventDataStores() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cloudTrailProvider);
      final request = ListEventDataStoresRequest();
      final response = await client.listEventDataStores(request);

      final eventDataStores = response.eventdatastores.map((eds) {
        return CloudTrailEventDataStore(
          eventDataStoreARN: eds.eventdatastorearn,
          name: eds.name,
          status: eds.status.name.replaceAll('EVENT_DATA_STORE_STATUS_', ''),
          created: eds.createdtimestamp.isNotEmpty
              ? DateTime.parse(eds.createdtimestamp)
              : DateTime.now(),
        );
      }).toList();

      state = state.copyWith(
        eventDataStores: eventDataStores,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final cloudTrailListProvider =
    NotifierProvider<CloudTrailListNotifier, CloudTrailState>(() {
      return CloudTrailListNotifier();
    });
