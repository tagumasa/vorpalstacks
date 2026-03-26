import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/cloudfront.pbgrpc.dart';
import '../src/generated/cloudfront.pb.dart';
import 'api_providers.dart';

class Distribution {
  final String id;
  final String arn;
  final String status;
  final bool enabled;
  final String domainName;
  final List<String> origins;
  final String defaultCacheBehavior;
  final String priceClass;
  final DateTime? lastModifiedTime;
  final bool isIpv6Enabled;
  final String comment;

  Distribution({
    required this.id,
    this.arn = '',
    this.status = '',
    this.enabled = false,
    this.domainName = '',
    this.origins = const [],
    this.defaultCacheBehavior = '',
    this.priceClass = '',
    this.lastModifiedTime,
    this.isIpv6Enabled = true,
    this.comment = '',
  });
}

class CloudFrontState {
  final List<Distribution> distributions;
  final bool isLoading;
  final String? error;

  const CloudFrontState({
    this.distributions = const [],
    this.isLoading = false,
    this.error,
  });

  CloudFrontState copyWith({
    List<Distribution>? distributions,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return CloudFrontState(
      distributions: distributions ?? this.distributions,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class CloudFrontListNotifier extends Notifier<CloudFrontState> {
  @override
  CloudFrontState build() => const CloudFrontState();

  Future<void> loadDistributions() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cloudFrontProvider);
      final request = ListDistributionsRequest();
      final response = await client.listDistributions(request);

      final distributions = response.distributionlist.items.map((d) {
        return Distribution(
          id: d.id,
          arn: d.arn,
          status: d.status,
          enabled: d.enabled,
          domainName: d.domainname,
        );
      });

      state = state.copyWith(
        distributions: distributions.toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final cloudfrontListProvider =
    NotifierProvider<CloudFrontListNotifier, CloudFrontState>(() {
      return CloudFrontListNotifier();
    });
