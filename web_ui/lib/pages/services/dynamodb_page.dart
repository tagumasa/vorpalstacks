import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/table_provider.dart';
import '../../widgets/common/status_badge.dart';

class DynamodbPage extends ConsumerStatefulWidget {
  const DynamodbPage({super.key});

  @override
  ConsumerState<DynamodbPage> createState() => _DynamodbPageState();
}

class _DynamodbPageState extends ConsumerState<DynamodbPage> {
  @override
  Widget build(BuildContext context) {
    final dynamoState = ref.watch(tableListProvider);
    final dynamoNotifier = ref.read(tableListProvider.notifier);
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
                  'DynamoDB',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: dynamoState.isLoading
                      ? null
                      : () async {
                          final nameController = TextEditingController();
                          final result = await showDialog<String>(
                            context: context,
                            builder: (context) => AlertDialog(
                              title: const Text('Create Table'),
                              content: TextField(
                                controller: nameController,
                                decoration: const InputDecoration(
                                  labelText: 'Table Name',
                                ),
                                autofocus: true,
                              ),
                              actions: [
                                TextButton(
                                  onPressed: () => Navigator.pop(context),
                                  child: const Text('Cancel'),
                                ),
                                ElevatedButton(
                                  onPressed: () => Navigator.pop(
                                    context,
                                    nameController.text,
                                  ),
                                  child: const Text('Create'),
                                ),
                              ],
                            ),
                          );
                          if (result != null && result.isNotEmpty) {
                            await dynamoNotifier.createTable(
                              result,
                              '${result}Id (S)',
                            );
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Table'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'NoSQL database service',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(
                    width: 300,
                    child: Card(
                      child: ListView.separated(
                        padding: const EdgeInsets.all(8),
                        itemCount: dynamoState.tables.length,
                        separatorBuilder: (_, _) => const SizedBox(height: 4),
                        itemBuilder: (context, index) {
                          final table = dynamoState.tables[index];
                          final isSelected =
                              dynamoState.selectedTable?.name == table.name;
                          return ListTile(
                            selected: isSelected,
                            selectedTileColor: cs.primary.withValues(
                              alpha: 0.1,
                            ),
                            leading: const Icon(Icons.table_chart_outlined),
                            title: Text(table.name),
                            subtitle: Text('${table.itemCount} items'),
                            trailing: StatusBadge(
                              label: table.status,
                              type: table.status == 'ACTIVE'
                                  ? StatusType.success
                                  : StatusType.pending,
                            ),
                            onTap: () => dynamoNotifier.selectTable(table),
                          );
                        },
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: dynamoState.selectedTable == null
                        ? _buildEmptySelection()
                        : _buildTableDetails(dynamoState, dynamoNotifier),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptySelection() {
    final cs = Theme.of(context).colorScheme;
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.table_chart_outlined,
            size: 64,
            color: cs.onSurface.withValues(alpha: 0.3),
          ),
          const SizedBox(height: 16),
          Text(
            'Select a table to view items',
            style: TextStyle(
              fontSize: 16,
              color: cs.onSurface.withValues(alpha: 0.6),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTableDetails(DynamoDBState state, DynamoDBNotifier notifier) {
    final cs = Theme.of(context).colorScheme;
    final table = state.selectedTable!;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  table.name,
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    _InfoChip(label: 'Items', value: '${table.itemCount}'),
                    const SizedBox(width: 16),
                    _InfoChip(
                      label: 'Size',
                      value: _formatBytes(table.sizeBytes),
                    ),
                    const SizedBox(width: 16),
                    _InfoChip(
                      label: 'Keys',
                      value: table.partitionKey.isEmpty
                          ? 'N/A'
                          : table.partitionKey,
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Text(
              'Items',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const Spacer(),
            ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Add Item'),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Expanded(
          child: Card(
            child: SingleChildScrollView(
              child: DataTable(
                headingRowColor: WidgetStateProperty.all(cs.surface),
                columns: const [
                  DataColumn(label: Text('userId')),
                  DataColumn(label: Text('name')),
                  DataColumn(label: Text('email')),
                  DataColumn(label: Text('age')),
                  DataColumn(label: Text('Actions')),
                ],
                rows: state.items
                    .map(
                      (item) => DataRow(
                        cells: [
                          DataCell(Text(item.data['userId']?.toString() ?? '')),
                          DataCell(Text(item.data['name']?.toString() ?? '')),
                          DataCell(Text(item.data['email']?.toString() ?? '')),
                          DataCell(Text(item.data['age']?.toString() ?? '')),
                          DataCell(
                            Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                IconButton(
                                  icon: const Icon(
                                    Icons.edit_outlined,
                                    size: 18,
                                  ),
                                  onPressed: () {},
                                ),
                                IconButton(
                                  icon: const Icon(
                                    Icons.delete_outline,
                                    size: 18,
                                  ),
                                  onPressed: () {},
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    )
                    .toList(),
              ),
            ),
          ),
        ),
      ],
    );
  }

  String _formatBytes(int bytes) {
    if (bytes < 1024) return '$bytes B';
    if (bytes < 1024 * 1024) return '${(bytes / 1024).toStringAsFixed(1)} KB';
    return '${(bytes / (1024 * 1024)).toStringAsFixed(1)} MB';
  }
}

class _InfoChip extends StatelessWidget {
  final String label;
  final String value;

  const _InfoChip({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(fontSize: 12, color: Color(0xFF879596)),
        ),
        Text(
          value,
          style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
        ),
      ],
    );
  }
}
