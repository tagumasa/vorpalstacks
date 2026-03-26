import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/cloudfront_provider.dart';

class CloudfrontPage extends ConsumerWidget {
  const CloudfrontPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cfState = ref.watch(cloudfrontListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'CloudFront',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Content delivery network',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: cfState.isLoading && cfState.distributions.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: cfState.distributions.length,
                      itemBuilder: (context, index) {
                        final dist = cfState.distributions[index];
                        Color statusColor = dist.status == 'Deployed'
                            ? const Color(0xFF4CAF50)
                            : const Color(0xFFFF9800);
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ListTile(
                            leading: Icon(
                              dist.enabled
                                  ? Icons.cloud_done_outlined
                                  : Icons.cloud_off_outlined,
                              color: dist.enabled
                                  ? cs.primary
                                  : cs.onSurface.withValues(alpha: 0.4),
                            ),
                            title: Row(
                              children: [
                                Text(dist.id),
                                const SizedBox(width: 8),
                                Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 6,
                                    vertical: 2,
                                  ),
                                  decoration: BoxDecoration(
                                    color: statusColor.withValues(alpha: 0.1),
                                    borderRadius: BorderRadius.circular(4),
                                  ),
                                  child: Text(
                                    dist.status,
                                    style: TextStyle(
                                      fontSize: 10,
                                      color: statusColor,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  dist.domainName,
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: cs.onSurface.withValues(alpha: 0.6),
                                  ),
                                ),
                                Text(
                                  dist.comment,
                                  style: TextStyle(
                                    fontSize: 11,
                                    color: cs.onSurface.withValues(alpha: 0.5),
                                  ),
                                ),
                              ],
                            ),
                            trailing: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Text(
                                  '${dist.origins.length} origin(s)',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: cs.onSurface.withValues(alpha: 0.6),
                                  ),
                                ),
                                Text(
                                  dist.enabled ? 'Enabled' : 'Disabled',
                                  style: TextStyle(
                                    fontSize: 11,
                                    color: dist.enabled
                                        ? const Color(0xFF4CAF50)
                                        : cs.onSurface.withValues(alpha: 0.5),
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
}
