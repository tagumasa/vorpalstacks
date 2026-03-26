import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/queue_provider.dart';

class SqsPage extends ConsumerWidget {
  const SqsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final sqsState = ref.watch(queueListProvider);
    final sqsNotifier = ref.read(queueListProvider.notifier);
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
                  'SQS',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: sqsState.isLoading
                      ? null
                      : () async {
                          final nameController = TextEditingController();
                          final isFifo = ValueNotifier(false);
                          final result = await showDialog<Map<String, dynamic>>(
                            context: context,
                            builder: (context) => AlertDialog(
                              title: const Text('Create Queue'),
                              content: Column(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  TextField(
                                    controller: nameController,
                                    decoration: const InputDecoration(
                                      labelText: 'Queue Name',
                                    ),
                                    autofocus: true,
                                  ),
                                  const SizedBox(height: 16),
                                  ValueListenableBuilder<bool>(
                                    valueListenable: isFifo,
                                    builder: (context, value, _) =>
                                        CheckboxListTile(
                                          title: const Text('FIFO Queue'),
                                          value: value,
                                          onChanged: (v) =>
                                              isFifo.value = v ?? false,
                                          contentPadding: EdgeInsets.zero,
                                        ),
                                  ),
                                ],
                              ),
                              actions: [
                                TextButton(
                                  onPressed: () => Navigator.pop(context),
                                  child: const Text('Cancel'),
                                ),
                                ElevatedButton(
                                  onPressed: () => Navigator.pop(context, {
                                    'name': nameController.text,
                                    'isFifo': isFifo.value,
                                  }),
                                  child: const Text('Create'),
                                ),
                              ],
                            ),
                          );
                          if (result != null &&
                              result['name'].toString().isNotEmpty) {
                            await sqsNotifier.createQueue(
                              result['name'],
                              isFifo: result['isFifo'] ?? false,
                            );
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Queue'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Message queue service',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: sqsState.isLoading && sqsState.queues.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : GridView.builder(
                      gridDelegate:
                          const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2,
                            crossAxisSpacing: 16,
                            mainAxisSpacing: 16,
                            childAspectRatio: 2,
                          ),
                      itemCount: sqsState.queues.length,
                      itemBuilder: (context, index) {
                        final queue = sqsState.queues[index];
                        return Card(
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  children: [
                                    Icon(
                                      Icons.message_outlined,
                                      color: cs.primary,
                                    ),
                                    const SizedBox(width: 8),
                                    Expanded(
                                      child: Text(
                                        queue.name,
                                        style: const TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w600,
                                        ),
                                        overflow: TextOverflow.ellipsis,
                                      ),
                                    ),
                                    Container(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 6,
                                        vertical: 2,
                                      ),
                                      decoration: BoxDecoration(
                                        color: queue.isFifo
                                            ? const Color(
                                                0xFFFF9800,
                                              ).withValues(alpha: 0.1)
                                            : const Color(
                                                0xFF2196F3,
                                              ).withValues(alpha: 0.1),
                                        borderRadius: BorderRadius.circular(4),
                                      ),
                                      child: Text(
                                        queue.isFifo ? 'FIFO' : 'Standard',
                                        style: TextStyle(
                                          fontSize: 10,
                                          color: queue.isFifo
                                              ? const Color(0xFFFF9800)
                                              : const Color(0xFF2196F3),
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                const Spacer(),
                                Row(
                                  children: [
                                    _MetricChip(
                                      label: 'Available',
                                      value: '${queue.availableMessages}',
                                    ),
                                    const SizedBox(width: 16),
                                    _MetricChip(
                                      label: 'In Flight',
                                      value: '${queue.notVisibleMessages}',
                                    ),
                                  ],
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

class _MetricChip extends StatelessWidget {
  final String label;
  final String value;
  const _MetricChip({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
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
