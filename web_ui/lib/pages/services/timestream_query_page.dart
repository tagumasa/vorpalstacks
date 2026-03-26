import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/timestream_query_provider.dart';
import '../../widgets/common/status_badge.dart';

class TimestreamQueryPage extends ConsumerWidget {
  const TimestreamQueryPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(timestreamQueryListProvider);
    final notifier = ref.read(timestreamQueryListProvider.notifier);
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
                  'Timestream Scheduled Queries',
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
                      : () => notifier.loadScheduledQueries(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage scheduled time-series queries',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: state.isLoading && state.scheduledQueries.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : state.scheduledQueries.isEmpty
                  ? Center(
                      child: Text(
                        'No scheduled queries found',
                        style: TextStyle(
                          color: cs.onSurface.withValues(alpha: 0.6),
                        ),
                      ),
                    )
                  : ListView.builder(
                      itemCount: state.scheduledQueries.length,
                      itemBuilder: (context, index) {
                        final query = state.scheduledQueries[index];
                        return _buildQueryCard(context, query);
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQueryCard(BuildContext context, TimestreamScheduledQuery query) {
    final cs = Theme.of(context).colorScheme;
    final isEnabled = query.state == 'ENABLED';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFE8F5E9).withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.schedule,
                color: Color(0xFF388E3C),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    query.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Last run: ${query.lastRunStatus}',
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
            ),
            StatusBadge(
              label: query.state,
              type: isEnabled ? StatusType.success : StatusType.pending,
            ),
          ],
        ),
      ),
    );
  }
}
