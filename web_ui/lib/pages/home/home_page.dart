import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/iam_provider.dart';
import '../../providers/table_provider.dart';
import '../../providers/queue_provider.dart';
import '../../providers/bucket_provider.dart';
import '../../providers/user_pool_provider.dart';
import '../../providers/apigateway_provider.dart';

class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final iamState = ref.watch(iamListProvider);
    final dynamoState = ref.watch(tableListProvider);
    final sqsState = ref.watch(queueListProvider);
    final s3State = ref.watch(bucketListProvider);
    final cognitoState = ref.watch(cognitoListProvider);
    final apiState = ref.watch(apiGatewayListProvider);

    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Dashboard',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Overview of your local AWS environment',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: GridView.count(
                crossAxisCount: 4,
                crossAxisSpacing: 16,
                mainAxisSpacing: 16,
                children: [
                  _buildMetricCard(
                    context,
                    'IAM Users',
                    '${iamState.users.length}',
                    Icons.people_outline,
                    const Color(0xFF2196F3),
                  ),
                  _buildMetricCard(
                    context,
                    'DynamoDB Tables',
                    '${dynamoState.tables.length}',
                    Icons.table_chart_outlined,
                    const Color(0xFF4CAF50),
                  ),
                  _buildMetricCard(
                    context,
                    'SQS Queues',
                    '${sqsState.queues.length}',
                    Icons.message_outlined,
                    const Color(0xFFFF9800),
                  ),
                  _buildMetricCard(
                    context,
                    'S3 Buckets',
                    '${s3State.buckets.length}',
                    Icons.folder_outlined,
                    const Color(0xFF9C27B0),
                  ),
                  _buildMetricCard(
                    context,
                    'Cognito Pools',
                    '${cognitoState.userPools.length}',
                    Icons.people,
                    const Color(0xFF00BCD4),
                  ),
                  _buildMetricCard(
                    context,
                    'API Gateways',
                    '${apiState.apis.length}',
                    Icons.api_outlined,
                    const Color(0xFFE91E63),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMetricCard(
    BuildContext context,
    String title,
    String value,
    IconData icon,
    Color color,
  ) {
    final cs = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: color.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(icon, color: color, size: 24),
                ),
                const Spacer(),
              ],
            ),
            const Spacer(),
            Text(
              value,
              style: TextStyle(
                fontSize: 32,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              title,
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
