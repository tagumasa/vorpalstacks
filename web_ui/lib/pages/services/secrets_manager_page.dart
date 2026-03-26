import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/secrets_manager_provider.dart';

class SecretsManagerPage extends ConsumerWidget {
  const SecretsManagerPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final secretsState = ref.watch(secretsManagerListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Secrets Manager',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Secrets and credentials management',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: secretsState.isLoading && secretsState.secrets.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: secretsState.secrets.length,
                      itemBuilder: (context, index) {
                        final secret = secretsState.secrets[index];
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ListTile(
                            leading: Icon(
                              Icons.key_outlined,
                              color: cs.primary,
                            ),
                            title: Text(
                              secret.name,
                              style: const TextStyle(
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                if (secret.description != null)
                                  Text(
                                    secret.description!,
                                    style: TextStyle(
                                      fontSize: 12,
                                      color: cs.onSurface.withValues(
                                        alpha: 0.6,
                                      ),
                                    ),
                                  ),
                                if (secret.lastAccessedDate != null)
                                  Text(
                                    'Last accessed: ${_formatDate(secret.lastAccessedDate!)}',
                                    style: TextStyle(
                                      fontSize: 11,
                                      color: cs.onSurface.withValues(
                                        alpha: 0.5,
                                      ),
                                    ),
                                  ),
                              ],
                            ),
                            trailing: _MetricChip(
                              label: 'Versions',
                              value: '${secret.versionIdsToStagesCount}',
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
    final diff = now.difference(date);
    if (diff.inDays == 0) return 'Today';
    if (diff.inDays == 1) return 'Yesterday';
    if (diff.inDays < 7) return '${diff.inDays} days ago';
    return '${date.day}/${date.month}/${date.year}';
  }
}

class _MetricChip extends StatelessWidget {
  final String label;
  final String value;
  const _MetricChip({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        Text(
          label,
          style: const TextStyle(fontSize: 10, color: Color(0xFF879596)),
        ),
        Text(
          value,
          style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
        ),
      ],
    );
  }
}
