import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/route53.pbgrpc.dart';
import '../src/generated/route53.pb.dart';
import 'api_providers.dart';

class HostedZone {
  final String id;
  final String name;
  final String arn;
  final int recordSetCount;
  final String? comment;
  final DateTime? creationTime;

  HostedZone({
    required this.id,
    required this.name,
    this.arn = '',
    this.recordSetCount = 0,
    this.comment,
    this.creationTime,
  });
}

class ResourceRecord {
  final String name;
  final String type;
  final String value;
  final int ttl;

  ResourceRecord({
    required this.name,
    required this.type,
    required this.value,
    this.ttl = 300,
  });
}

class Route53State {
  final List<HostedZone> hostedZones;
  final List<ResourceRecord> records;
  final bool isLoading;
  final String? error;

  const Route53State({
    this.hostedZones = const [],
    this.records = const [],
    this.isLoading = false,
    this.error,
  });

  Route53State copyWith({
    List<HostedZone>? hostedZones,
    List<ResourceRecord>? records,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return Route53State(
      hostedZones: hostedZones ?? this.hostedZones,
      records: records ?? this.records,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class Route53ListNotifier extends Notifier<Route53State> {
  @override
  Route53State build() => const Route53State();

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(route53Provider);
      final request = ListHostedZonesRequest();
      final response = await client.listHostedZones(request);

      final hostedZones = response.hostedzones.map((hz) {
        return HostedZone(
          id: hz.id,
          name: hz.name,
          recordSetCount: hz.resourcerecordsetcount.toInt(),
        );
      });

      state = state.copyWith(
        hostedZones: hostedZones.toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final route53ListProvider = NotifierProvider<Route53ListNotifier, Route53State>(
  () {
    return Route53ListNotifier();
  },
);
