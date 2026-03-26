import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/sts_provider.dart';

class STSPage extends ConsumerWidget {
  const STSPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(stsListProvider);
    final notifier = ref.read(stsListProvider.notifier);
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
                  'STS Caller Identity',
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
                      : () => notifier.loadCallerIdentity(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'View current caller identity information',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            if (state.isLoading)
              const Center(child: CircularProgressIndicator())
            else if (state.error != null)
              Center(
                child: Text(
                  'Error: ${state.error}',
                  style: TextStyle(color: cs.error),
                ),
              )
            else if (state.identity != null)
              _buildIdentityCard(context, state.identity!),
          ],
        ),
      ),
    );
  }

  Widget _buildIdentityCard(BuildContext context, STSIdentity identity) {
    final cs = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: const Color(0xFFF3E5F5).withValues(alpha: 0.3),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Icon(
                    Icons.account_circle,
                    color: Color(0xFF7B1FA2),
                    size: 32,
                  ),
                ),
                const SizedBox(width: 16),
                Text(
                  'Caller Identity',
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),
            _buildInfoRow(context, 'User ID', identity.userId),
            const SizedBox(height: 16),
            _buildInfoRow(context, 'Account', identity.account),
            const SizedBox(height: 16),
            _buildInfoRow(context, 'ARN', identity.arn),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(BuildContext context, String label, String value) {
    final cs = Theme.of(context).colorScheme;

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 100,
          child: Text(
            label,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: cs.onSurface.withValues(alpha: 0.6),
            ),
          ),
        ),
        Expanded(
          child: SelectableText(
            value,
            style: TextStyle(fontSize: 14, color: cs.onSurface),
          ),
        ),
      ],
    );
  }
}
