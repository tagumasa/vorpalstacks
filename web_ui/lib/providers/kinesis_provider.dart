import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/kinesis.pbgrpc.dart';
import 'api_providers.dart';

class KinesisStream {
  final String streamName;
  final String streamARN;
  final String status;
  final DateTime created;

  KinesisStream({
    required this.streamName,
    required this.streamARN,
    required this.status,
    required this.created,
  });
}

class KinesisState {
  final List<KinesisStream> streams;
  final bool isLoading;
  final String? error;

  const KinesisState({
    this.streams = const [],
    this.isLoading = false,
    this.error,
  });

  KinesisState copyWith({
    List<KinesisStream>? streams,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return KinesisState(
      streams: streams ?? this.streams,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class KinesisListNotifier extends Notifier<KinesisState> {
  @override
  KinesisState build() => const KinesisState();

  Future<void> loadStreams() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(kinesisProvider);
      final request = ListStreamsInput();
      final response = await client.listStreams(request);

      final streams = response.streamsummaries.map((s) {
        return KinesisStream(
          streamName: s.streamname,
          streamARN: s.streamarn,
          status: s.streamstatus.name.replaceAll('STREAM_STATUS_', ''),
          created: s.streamcreationtimestamp.isNotEmpty
              ? DateTime.parse(s.streamcreationtimestamp)
              : DateTime.now(),
        );
      }).toList();

      state = state.copyWith(streams: streams, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final kinesisListProvider = NotifierProvider<KinesisListNotifier, KinesisState>(
  () {
    return KinesisListNotifier();
  },
);
