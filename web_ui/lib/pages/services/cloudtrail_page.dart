import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/cloudtrail_provider.dart';
import '../../widgets/common/status_badge.dart';

class CloudTrailPage extends ConsumerWidget {
  const CloudTrailPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(cloudTrailListProvider);
    final notifier = ref.read(cloudTrailListProvider.notifier);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'CloudTrail',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Track API activity and events',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: Text(
                    'Trails',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w500,
                      color: cs.onSurface,
                    ),
                  ),
                ),
                IconButton(
                  onPressed: state.isLoading
                      ? null
                      : () => notifier.loadTrails(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 12),
            if (state.isLoading && state.trails.isEmpty)
              const Center(child: CircularProgressIndicator())
            else if (state.trails.isEmpty)
              Text(
                'No trails found',
                style: TextStyle(color: cs.onSurface.withValues(alpha: 0.6)),
              )
            else
              ...state.trails.map((trail) => _buildTrailCard(context, trail)),
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: Text(
                    'Event Data Stores',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w500,
                      color: cs.onSurface,
                    ),
                  ),
                ),
                IconButton(
                  onPressed: state.isLoading
                      ? null
                      : () => notifier.loadEventDataStores(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Expanded(
              child: state.isLoading && state.eventDataStores.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : state.eventDataStores.isEmpty
                  ? Text(
                      'No event data stores found',
                      style: TextStyle(
                        color: cs.onSurface.withValues(alpha: 0.6),
                      ),
                    )
                  : ListView.builder(
                      itemCount: state.eventDataStores.length,
                      itemBuilder: (context, index) {
                        return _buildEventDataStoreCard(
                          context,
                          state.eventDataStores[index],
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTrailCard(BuildContext context, CloudTrailTrail trail) {
    final cs = Theme.of(context).colorScheme;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFB3E5FC).withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.route,
                color: Color(0xFF0288D1),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    trail.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    trail.homeRegion,
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
            ),
            Text(
              trail.trailARN,
              style: TextStyle(
                fontSize: 11,
                color: cs.onSurface.withValues(alpha: 0.5),
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEventDataStoreCard(
    BuildContext context,
    CloudTrailEventDataStore eds,
  ) {
    final cs = Theme.of(context).colorScheme;
    final isActive = eds.status == 'ENABLED';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFB3E5FC).withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.storage,
                color: Color(0xFF0288D1),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    eds.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    eds.eventDataStoreARN,
                    style: TextStyle(
                      fontSize: 11,
                      color: cs.onSurface.withValues(alpha: 0.5),
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            StatusBadge(
              label: eds.status,
              type: isActive ? StatusType.success : StatusType.pending,
            ),
          ],
        ),
      ),
    );
  }
}
