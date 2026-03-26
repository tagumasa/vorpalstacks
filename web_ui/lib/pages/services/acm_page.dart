import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/acm_provider.dart';

class AcmPage extends ConsumerWidget {
  const AcmPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final acmState = ref.watch(acmListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Certificate Manager',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'SSL/TLS certificate management',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: acmState.isLoading && acmState.certificates.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: acmState.certificates.length,
                      itemBuilder: (context, index) {
                        final cert = acmState.certificates[index];
                        Color statusColor;
                        switch (cert.status) {
                          case 'ISSUED':
                            statusColor = const Color(0xFF4CAF50);
                            break;
                          case 'PENDING_VALIDATION':
                            statusColor = const Color(0xFFFF9800);
                            break;
                          case 'EXPIRED':
                            statusColor = const Color(0xFFF44336);
                            break;
                          default:
                            statusColor = const Color(0xFF9E9E9E);
                        }
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ListTile(
                            leading: Icon(
                              cert.inUse
                                  ? Icons.verified
                                  : Icons.verified_outlined,
                              color: statusColor,
                            ),
                            title: Text(cert.domainName),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                if (cert.subjectAlternativeNames.isNotEmpty)
                                  Text(
                                    'SANs: ${cert.subjectAlternativeNames.take(3).join(", ")}${cert.subjectAlternativeNames.length > 3 ? "..." : ""}',
                                    style: TextStyle(
                                      fontSize: 12,
                                      color: cs.onSurface.withValues(
                                        alpha: 0.6,
                                      ),
                                    ),
                                  ),
                                if (cert.notAfter != null)
                                  Text(
                                    'Expires: ${_formatDate(cert.notAfter!)}',
                                    style: TextStyle(
                                      fontSize: 11,
                                      color: cert.status == 'EXPIRED'
                                          ? statusColor
                                          : cs.onSurface.withValues(alpha: 0.5),
                                    ),
                                  ),
                              ],
                            ),
                            trailing: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 8,
                                    vertical: 4,
                                  ),
                                  decoration: BoxDecoration(
                                    color: statusColor.withValues(alpha: 0.1),
                                    borderRadius: BorderRadius.circular(4),
                                  ),
                                  child: Text(
                                    cert.status,
                                    style: TextStyle(
                                      color: statusColor,
                                      fontWeight: FontWeight.w500,
                                      fontSize: 12,
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  cert.inUse ? 'In use' : 'Not in use',
                                  style: TextStyle(
                                    fontSize: 11,
                                    color: cs.onSurface.withValues(alpha: 0.5),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = date.difference(now);
    if (diff.inDays < 0) return 'Expired';
    if (diff.inDays == 0) return 'Today';
    if (diff.inDays < 30) return 'In ${diff.inDays} days';
    return '${date.day}/${date.month}/${date.year}';
  }
}
