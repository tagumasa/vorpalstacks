import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/acm.pbgrpc.dart';
import '../src/generated/acm.pb.dart';
import 'api_providers.dart';

class Certificate {
  final String certificateArn;
  final String domainName;
  final String status;
  final String type;
  final List<String> subjectAlternativeNames;
  final DateTime? notBefore;
  final DateTime? notAfter;
  final DateTime? createdAt;
  final bool inUse;
  final String? keyAlgorithm;

  Certificate({
    required this.certificateArn,
    required this.domainName,
    this.status = '',
    this.type = '',
    this.subjectAlternativeNames = const [],
    this.notBefore,
    this.notAfter,
    this.createdAt,
    this.inUse = false,
    this.keyAlgorithm,
  });
}

class ACMState {
  final List<Certificate> certificates;
  final bool isLoading;
  final String? error;

  const ACMState({
    this.certificates = const [],
    this.isLoading = false,
    this.error,
  });

  ACMState copyWith({
    List<Certificate>? certificates,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return ACMState(
      certificates: certificates ?? this.certificates,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class ACMListNotifier extends Notifier<ACMState> {
  @override
  ACMState build() => const ACMState();

  Future<void> loadCertificates() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(acmProvider);
      final request = ListCertificatesRequest();
      final response = await client.listCertificates(request);

      final certificates = response.certificatesummarylist.map((c) {
        return Certificate(
          certificateArn: c.certificatearn,
          domainName: c.domainname,
          status: c.status.name.replaceAll('CERTIFICATE_STATUS_', ''),
        );
      });

      state = state.copyWith(
        certificates: certificates.toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final acmListProvider = NotifierProvider<ACMListNotifier, ACMState>(() {
  return ACMListNotifier();
});
