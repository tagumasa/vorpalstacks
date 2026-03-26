import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/apigateway_provider.dart';

class ApigatewayPage extends ConsumerWidget {
  const ApigatewayPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final apiState = ref.watch(apiGatewayListProvider);
    final apiNotifier = ref.read(apiGatewayListProvider.notifier);
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
                  'API Gateway',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: apiState.isLoading
                      ? null
                      : () => _showCreateDialog(context, apiNotifier),
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create API'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'API creation and management',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: apiState.isLoading && apiState.apis.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.separated(
                      itemCount: apiState.apis.length,
                      separatorBuilder: (_, _) => const SizedBox(height: 12),
                      itemBuilder: (context, index) =>
                          _ApiCard(api: apiState.apis[index]),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _showCreateDialog(
    BuildContext context,
    ApiGatewayNotifier apiNotifier,
  ) async {
    final nameController = TextEditingController();
    final descController = TextEditingController();
    final result = await showDialog<Map<String, String>>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create API'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: const InputDecoration(labelText: 'API Name'),
              autofocus: true,
            ),
            const SizedBox(height: 16),
            TextField(
              controller: descController,
              decoration: const InputDecoration(labelText: 'Description'),
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
              'description': descController.text,
            }),
            child: const Text('Create'),
          ),
        ],
      ),
    );
    if (result != null && result['name']!.isNotEmpty) {
      await apiNotifier.createAPI(result['name']!, result['description'] ?? '');
    }
  }
}

class _ApiCard extends StatelessWidget {
  final ApiGateway api;
  const _ApiCard({required this.api});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: const Color(0xFFE91E63).withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.api_outlined, color: Color(0xFFE91E63)),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    api.name,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    api.endpoint,
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: const Color(0xFF4CAF50).withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: const Text(
                'prod',
                style: TextStyle(
                  fontSize: 12,
                  color: Color(0xFF4CAF50),
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
            const SizedBox(width: 8),
            IconButton(icon: const Icon(Icons.chevron_right), onPressed: () {}),
          ],
        ),
      ),
    );
  }
}
