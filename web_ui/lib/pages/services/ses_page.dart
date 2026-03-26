import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/ses_provider.dart';
import '../../widgets/common/status_badge.dart';

class SESPage extends ConsumerWidget {
  const SESPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(sesListProvider);
    final notifier = ref.read(sesListProvider.notifier);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  'SES Email Identities',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                IconButton(
                  onPressed: state.isLoading
                      ? null
                      : () => notifier.loadIdentities(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage email sending identities',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: state.isLoading && state.identities.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : state.identities.isEmpty
                  ? Center(
                      child: Text(
                        'No email identities found',
                        style: TextStyle(
                          color: cs.onSurface.withValues(alpha: 0.6),
                        ),
                      ),
                    )
                  : ListView.builder(
                      itemCount: state.identities.length,
                      itemBuilder: (context, index) {
                        final identity = state.identities[index];
                        return _buildIdentityCard(context, identity);
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildIdentityCard(BuildContext context, SESIdentity identity) {
    final cs = Theme.of(context).colorScheme;
    final isVerified = identity.verificationStatus == 'SUCCESS';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFE3F2FD).withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.email,
                color: Color(0xFF1976D2),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    identity.identityName,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${identity.identityType} • Sending: ${identity.sendingEnabled ? "Enabled" : "Disabled"}',
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
            ),
            StatusBadge(
              label: isVerified ? 'Verified' : identity.verificationStatus,
              type: isVerified ? StatusType.success : StatusType.pending,
            ),
          ],
        ),
      ),
    );
  }
}
