import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/scheduler_provider.dart';
import '../../widgets/common/status_badge.dart';

class SchedulerPage extends ConsumerWidget {
  const SchedulerPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final schedulerState = ref.watch(schedulerListProvider);
    final schedulerNotifier = ref.read(schedulerListProvider.notifier);
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
                  'EventBridge Scheduler',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: schedulerState.isLoading
                      ? null
                      : () async {
                          final result = await showDialog<Map<String, dynamic>>(
                            context: context,
                            builder: (context) => _CreateScheduleDialog(),
                          );
                          if (result != null) {
                            await schedulerNotifier.createSchedule(
                              name: result['name'],
                              scheduleExpression: result['expression'],
                              targetArn: result['targetArn'],
                            );
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Schedule'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage your scheduled tasks',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child:
                  schedulerState.isLoading && schedulerState.schedules.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: schedulerState.schedules.length,
                      itemBuilder: (context, index) {
                        final schedule = schedulerState.schedules[index];
                        return _buildScheduleCard(
                          context,
                          schedule,
                          schedulerNotifier,
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildScheduleCard(
    BuildContext context,
    Schedule schedule,
    SchedulerListNotifier notifier,
  ) {
    final cs = Theme.of(context).colorScheme;
    final isEnabled = schedule.state == 'ENABLED';

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
              child: const Icon(
                Icons.schedule,
                color: Color(0xFFE7157B),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    schedule.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    schedule.arn.isNotEmpty ? schedule.arn : 'No ARN',
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
              label: schedule.state,
              type: isEnabled ? StatusType.success : StatusType.pending,
            ),
            const SizedBox(width: 12),
            PopupMenuButton(
              icon: Icon(
                Icons.more_vert,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
              itemBuilder: (context) => [
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
                if (value == 'delete') {
                  final confirmed = await showDialog<bool>(
                    context: context,
                    builder: (context) => AlertDialog(
                      title: const Text('Delete Schedule'),
                      content: Text(
                        'Are you sure you want to delete "${schedule.name}"?',
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
                    await notifier.deleteSchedule(schedule.name);
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

class _CreateScheduleDialog extends StatefulWidget {
  @override
  State<_CreateScheduleDialog> createState() => _CreateScheduleDialogState();
}

class _CreateScheduleDialogState extends State<_CreateScheduleDialog> {
  final _nameController = TextEditingController();
  final _expressionController = TextEditingController(text: 'rate(1 day)');
  final _targetController = TextEditingController();

  @override
  void dispose() {
    _nameController.dispose();
    _expressionController.dispose();
    _targetController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Create Schedule'),
      content: SizedBox(
        width: 400,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Schedule Name',
                hintText: 'my-schedule',
              ),
              autofocus: true,
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _expressionController,
              decoration: const InputDecoration(
                labelText: 'Schedule Expression',
                hintText: 'rate(1 day) or cron(0 12 * * ? *)',
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _targetController,
              decoration: const InputDecoration(
                labelText: 'Target ARN',
                hintText: 'arn:aws:lambda:...',
              ),
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
            if (_nameController.text.isNotEmpty &&
                _targetController.text.isNotEmpty) {
              Navigator.pop(context, {
                'name': _nameController.text,
                'expression': _expressionController.text,
                'targetArn': _targetController.text,
              });
            }
          },
          child: const Text('Create'),
        ),
      ],
    );
  }
}
