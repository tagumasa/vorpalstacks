import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/eventbridge_provider.dart';
import '../../widgets/common/status_badge.dart';

class EventBridgePage extends ConsumerWidget {
  const EventBridgePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final eventBridgeState = ref.watch(eventBridgeListProvider);
    final eventBridgeNotifier = ref.read(eventBridgeListProvider.notifier);
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
                  'EventBridge',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: eventBridgeState.isLoading
                      ? null
                      : () async {
                          final result = await showDialog<Map<String, dynamic>>(
                            context: context,
                            builder: (context) => _CreateRuleDialog(),
                          );
                          if (result != null) {
                            await eventBridgeNotifier.createRule(
                              name: result['name'],
                              scheduleExpression: result['scheduleExpression'],
                              eventPattern: result['eventPattern'],
                            );
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Rule'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage your event rules',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child:
                  eventBridgeState.isLoading && eventBridgeState.rules.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: eventBridgeState.rules.length,
                      itemBuilder: (context, index) {
                        final rule = eventBridgeState.rules[index];
                        return _buildRuleCard(
                          context,
                          rule,
                          eventBridgeNotifier,
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRuleCard(
    BuildContext context,
    EventRule rule,
    EventBridgeListNotifier notifier,
  ) {
    final cs = Theme.of(context).colorScheme;
    final isEnabled = rule.state == 'ENABLED';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFE7157B).withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.bolt, color: Color(0xFFE7157B), size: 24),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    rule.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  if (rule.description.isNotEmpty)
                    Text(
                      rule.description,
                      style: TextStyle(
                        fontSize: 12,
                        color: cs.onSurface.withValues(alpha: 0.6),
                      ),
                    ),
                  if (rule.scheduleExpression.isNotEmpty)
                    Text(
                      'Schedule: ${rule.scheduleExpression}',
                      style: TextStyle(
                        fontSize: 12,
                        color: cs.onSurface.withValues(alpha: 0.5),
                        fontFamily: 'monospace',
                      ),
                    ),
                ],
              ),
            ),
            StatusBadge(
              label: rule.state,
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
                const PopupMenuItem(
                  value: 'delete',
                  child: Row(
                    children: [
                      Icon(Icons.delete_outline, size: 18),
                      SizedBox(width: 12),
                      Text('Delete'),
                    ],
                  ),
                ),
              ],
              onSelected: (value) async {
                if (value == 'enable') {
                  await notifier.enableRule(rule.name);
                } else if (value == 'disable') {
                  await notifier.disableRule(rule.name);
                } else if (value == 'delete') {
                  final confirmed = await showDialog<bool>(
                    context: context,
                    builder: (context) => AlertDialog(
                      title: const Text('Delete Rule'),
                      content: Text(
                        'Are you sure you want to delete "${rule.name}"?',
                      ),
                      actions: [
                        TextButton(
                          onPressed: () => Navigator.pop(context, false),
                          child: const Text('Cancel'),
                        ),
                        ElevatedButton(
                          onPressed: () => Navigator.pop(context, true),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.red,
                          ),
                          child: const Text('Delete'),
                        ),
                      ],
                    ),
                  );
                  if (confirmed == true) {
                    await notifier.deleteRule(rule.name);
                  }
                }
              },
            ),
          ],
        ),
      ),
    );
  }
}

class _CreateRuleDialog extends StatefulWidget {
  @override
  State<_CreateRuleDialog> createState() => _CreateRuleDialogState();
}

class _CreateRuleDialogState extends State<_CreateRuleDialog> {
  final _nameController = TextEditingController();
  final _scheduleController = TextEditingController();
  final _patternController = TextEditingController();
  bool _isSchedule = true;

  @override
  void dispose() {
    _nameController.dispose();
    _scheduleController.dispose();
    _patternController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Create Rule'),
      content: SizedBox(
        width: 400,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Rule Name',
                hintText: 'my-rule',
              ),
              autofocus: true,
            ),
            const SizedBox(height: 16),
            SegmentedButton<bool>(
              segments: const [
                ButtonSegment(value: true, label: Text('Schedule')),
                ButtonSegment(value: false, label: Text('Event Pattern')),
              ],
              selected: {_isSchedule},
              onSelectionChanged: (v) => setState(() => _isSchedule = v.first),
            ),
            const SizedBox(height: 16),
            if (_isSchedule)
              TextField(
                controller: _scheduleController,
                decoration: const InputDecoration(
                  labelText: 'Schedule Expression',
                  hintText: 'rate(1 day) or cron(0 12 * * ? *)',
                ),
              )
            else
              TextField(
                controller: _patternController,
                decoration: const InputDecoration(
                  labelText: 'Event Pattern (JSON)',
                  hintText: '{"source": ["aws.ec2"]}',
                ),
                maxLines: 3,
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () {
            if (_nameController.text.isNotEmpty) {
              Navigator.pop(context, {
                'name': _nameController.text,
                'scheduleExpression': _isSchedule
                    ? _scheduleController.text
                    : null,
                'eventPattern': !_isSchedule ? _patternController.text : null,
              });
            }
          },
          child: const Text('Create'),
        ),
      ],
    );
  }
}
