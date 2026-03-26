import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/cloudwatch_provider.dart';

class CloudwatchPage extends ConsumerWidget {
  const CloudwatchPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cwState = ref.watch(cloudwatchListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'CloudWatch',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Monitoring and logging service',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: cwState.isLoading && cwState.logGroups.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : DefaultTabController(
                      length: 2,
                      child: Column(
                        children: [
                          TabBar(
                            labelColor: cs.primary,
                            unselectedLabelColor: cs.onSurface.withValues(
                              alpha: 0.6,
                            ),
                            indicatorColor: cs.primary,
                            tabs: const [
                              Tab(text: 'Log Groups'),
                              Tab(text: 'Alarms'),
                            ],
                          ),
                          const SizedBox(height: 16),
                          Expanded(
                            child: TabBarView(
                              children: [
                                _buildLogGroups(cwState.logGroups, cs),
                                _buildAlarms(cwState.alarms, cs),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLogGroups(List<LogGroup> logGroups, ColorScheme cs) {
    return GridView.builder(
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
        childAspectRatio: 2.5,
      ),
      itemCount: logGroups.length,
      itemBuilder: (context, index) {
        final lg = logGroups[index];
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.article_outlined, color: cs.primary),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        lg.logGroupName,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                Row(
                  children: [
                    _MetricChip(
                      label: 'Retention',
                      value: '${lg.retentionInDays}d',
                    ),
                    const SizedBox(width: 16),
                    _MetricChip(
                      label: 'Size',
                      value: _formatBytes(lg.storedBytes ?? 0),
                    ),
                  ],
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildAlarms(List<MetricAlarm> alarms, ColorScheme cs) {
    return ListView.builder(
      itemCount: alarms.length,
      itemBuilder: (context, index) {
        final alarm = alarms[index];
        Color stateColor;
        switch (alarm.state) {
          case 'OK':
            stateColor = const Color(0xFF4CAF50);
            break;
          case 'ALARM':
            stateColor = const Color(0xFFF44336);
            break;
          default:
            stateColor = const Color(0xFF9E9E9E);
        }
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: Icon(Icons.alarm, color: stateColor),
            title: Text(alarm.alarmName),
            subtitle: Text('${alarm.namespace}/${alarm.metricName}'),
            trailing: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: stateColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: Text(
                alarm.state,
                style: TextStyle(
                  color: stateColor,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
          ),
        );
      },
    );
  }

  String _formatBytes(int bytes) {
    if (bytes < 1024) return '$bytes B';
    if (bytes < 1024 * 1024) return '${(bytes / 1024).toStringAsFixed(1)} KB';
    return '${(bytes / (1024 * 1024)).toStringAsFixed(1)} MB';
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
