import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/kinesis_provider.dart';
import '../../widgets/common/status_badge.dart';

class KinesisPage extends ConsumerWidget {
  const KinesisPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(kinesisListProvider);
    final notifier = ref.read(kinesisListProvider.notifier);
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
                  'Kinesis Streams',
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
                      : () => notifier.loadStreams(),
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage real-time data streams',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: state.isLoading && state.streams.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : state.streams.isEmpty
                  ? Center(
                      child: Text(
                        'No streams found',
                        style: TextStyle(
                          color: cs.onSurface.withValues(alpha: 0.6),
                        ),
                      ),
                    )
                  : ListView.builder(
                      itemCount: state.streams.length,
                      itemBuilder: (context, index) {
                        final stream = state.streams[index];
                        return _buildStreamCard(context, stream);
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStreamCard(BuildContext context, KinesisStream stream) {
    final cs = Theme.of(context).colorScheme;
    final isActive = stream.status == 'ACTIVE';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: const Color(0xFFFFE0B2).withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(
                Icons.stream,
                color: Color(0xFFE65100),
                size: 24,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    stream.streamName,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: cs.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    stream.streamARN,
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
              label: stream.status,
              type: isActive ? StatusType.success : StatusType.pending,
            ),
          ],
        ),
      ),
    );
  }
}
