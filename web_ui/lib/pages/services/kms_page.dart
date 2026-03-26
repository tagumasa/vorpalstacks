import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/kms_provider.dart';
import '../../widgets/common/status_badge.dart';

class KMSPage extends ConsumerWidget {
  const KMSPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final kmsState = ref.watch(kmsListProvider);
    final kmsNotifier = ref.read(kmsListProvider.notifier);
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
                  'KMS',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: kmsState.isLoading
                      ? null
                      : () async {
                          final result = await showDialog<String>(
                            context: context,
                            builder: (context) => _CreateKeyDialog(),
                          );
                          if (result != null) {
                            await kmsNotifier.createKey(description: result);
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Key'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage your encryption keys',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: kmsState.isLoading && kmsState.keys.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: kmsState.keys.length,
                      itemBuilder: (context, index) {
                        final key = kmsState.keys[index];
                        return _buildKeyCard(context, key, kmsNotifier);
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildKeyCard(
    BuildContext context,
    KMSKey key,
    KMSListNotifier notifier,
  ) {
    final cs = Theme.of(context).colorScheme;
    final isEnabled = key.state == 'Enabled';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFFF9900).withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.vpn_key,
                color: Color(0xFFFF9900),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    key.keyId,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    key.keyArn,
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            StatusBadge(
              label: key.state,
              type: isEnabled ? StatusType.success : StatusType.pending,
            ),
            const SizedBox(width: 12),
            PopupMenuButton(
              icon: Icon(
                Icons.more_vert,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
              itemBuilder: (context) => [
                PopupMenuItem(
                  value: isEnabled ? 'disable' : 'enable',
                  child: Row(
                    children: [
                      Icon(
                        isEnabled ? Icons.block : Icons.check_circle,
                        size: 18,
                      ),
                      const SizedBox(width: 12),
                      Text(isEnabled ? 'Disable' : 'Enable'),
                    ],
                  ),
                ),
              ],
              onSelected: (value) async {
                if (value == 'enable') {
                  await notifier.enableKey(key.keyId);
                } else if (value == 'disable') {
                  await notifier.disableKey(key.keyId);
                }
              },
            ),
          ],
        ),
      ),
    );
  }
}

class _CreateKeyDialog extends StatefulWidget {
  @override
  State<_CreateKeyDialog> createState() => _CreateKeyDialogState();
}

class _CreateKeyDialogState extends State<_CreateKeyDialog> {
  final _descriptionController = TextEditingController();

  @override
  void dispose() {
    _descriptionController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Create Key'),
      content: SizedBox(
        width: 400,
        child: TextField(
          controller: _descriptionController,
          decoration: const InputDecoration(
            labelText: 'Description (optional)',
            hintText: 'My encryption key',
          ),
          autofocus: true,
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () => Navigator.pop(context, _descriptionController.text),
          child: const Text('Create'),
        ),
      ],
    );
  }
}
