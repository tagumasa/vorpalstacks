import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../src/providers/config_provider.dart';

class SettingsPage extends ConsumerStatefulWidget {
  const SettingsPage({super.key});

  @override
  ConsumerState<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends ConsumerState<SettingsPage> {
  String? _selectedCategory;
  final _searchController = TextEditingController();
  String _searchQuery = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final configState = ref.watch(configProvider);

    return Scaffold(
      backgroundColor: cs.surface,
      body: Column(
        children: [
          _buildHeader(cs),
          _buildFilters(cs, configState),
          Expanded(
            child: configState.isLoading
                ? const Center(child: CircularProgressIndicator())
                : configState.error != null
                ? _buildError(cs, configState.error!)
                : _buildConfigList(cs, configState),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader(ColorScheme cs) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Settings',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Configure your VorpalStack environment',
                  style: TextStyle(
                    fontSize: 14,
                    color: cs.onSurface.withValues(alpha: 0.6),
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(configProvider.notifier).loadConfig(),
            tooltip: 'Refresh',
          ),
        ],
      ),
    );
  }

  Widget _buildFilters(ColorScheme cs, ConfigState configState) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search settings...',
                prefixIcon: const Icon(Icons.search),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
                contentPadding: const EdgeInsets.symmetric(horizontal: 12),
                isDense: true,
              ),
              onChanged: (value) {
                setState(() => _searchQuery = value);
              },
            ),
          ),
          const SizedBox(width: 16),
          SizedBox(
            width: 200,
            child: DropdownButtonFormField<String>(
              initialValue: _selectedCategory,
              decoration: InputDecoration(
                hintText: 'All Categories',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
                contentPadding: const EdgeInsets.symmetric(horizontal: 12),
                isDense: true,
              ),
              items: [
                const DropdownMenuItem(
                  value: null,
                  child: Text('All Categories'),
                ),
                ...configState.categories.map(
                  (c) => DropdownMenuItem(
                    value: c,
                    child: Text(_formatCategory(c)),
                  ),
                ),
              ],
              onChanged: (value) {
                setState(() => _selectedCategory = value);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildError(ColorScheme cs, String error) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.error_outline, size: 48, color: cs.error),
          const SizedBox(height: 16),
          Text('Failed to load settings', style: TextStyle(color: cs.error)),
          const SizedBox(height: 8),
          Text(
            error,
            style: TextStyle(color: cs.onSurface.withValues(alpha: 0.6)),
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () => ref.read(configProvider.notifier).loadConfig(),
            child: const Text('Retry'),
          ),
        ],
      ),
    );
  }

  Widget _buildConfigList(ColorScheme cs, ConfigState configState) {
    final filteredEntries = _filterEntries(configState.entries);

    if (filteredEntries.isEmpty) {
      return Center(
        child: Text(
          'No settings found',
          style: TextStyle(color: cs.onSurface.withValues(alpha: 0.6)),
        ),
      );
    }

    final grouped = _groupByCategory(filteredEntries);

    return ListView.builder(
      padding: const EdgeInsets.all(24),
      itemCount: grouped.length,
      itemBuilder: (context, index) {
        final category = grouped.keys.elementAt(index);
        final entries = grouped[category]!;

        return _SettingsSection(
          title: _formatCategory(category),
          children: entries
              .map(
                (e) => _ConfigTile(
                  entry: e,
                  onUpdate: (value) async {
                    final messenger = ScaffoldMessenger.of(context);
                    final success = await ref
                        .read(configProvider.notifier)
                        .updateConfig(e.key, value);
                    if (mounted && !success) {
                      messenger.showSnackBar(
                        const SnackBar(
                          content: Text('Failed to update setting'),
                        ),
                      );
                    }
                  },
                  onReset: () async {
                    final messenger = ScaffoldMessenger.of(context);
                    final success = await ref
                        .read(configProvider.notifier)
                        .resetConfig(e.key);
                    if (mounted && !success) {
                      messenger.showSnackBar(
                        const SnackBar(
                          content: Text('Failed to reset setting'),
                        ),
                      );
                    }
                  },
                ),
              )
              .toList(),
        );
      },
    );
  }

  List<ConfigEntry> _filterEntries(List<ConfigEntry> entries) {
    var filtered = entries;

    if (_selectedCategory != null) {
      filtered = filtered
          .where((e) => e.category == _selectedCategory)
          .toList();
    }

    if (_searchQuery.isNotEmpty) {
      final query = _searchQuery.toLowerCase();
      filtered = filtered.where((e) {
        return e.key.toLowerCase().contains(query) ||
            e.description.toLowerCase().contains(query);
      }).toList();
    }

    return filtered;
  }

  Map<String, List<ConfigEntry>> _groupByCategory(List<ConfigEntry> entries) {
    final grouped = <String, List<ConfigEntry>>{};
    for (final entry in entries) {
      grouped.putIfAbsent(entry.category, () => []).add(entry);
    }
    return grouped;
  }

  String _formatCategory(String category) {
    return category
        .split('_')
        .map((word) {
          if (word.isEmpty) return word;
          return word[0].toUpperCase() + word.substring(1);
        })
        .join(' ');
  }
}

class _SettingsSection extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SettingsSection({required this.title, required this.children});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: cs.onSurface.withValues(alpha: 0.6),
          ),
        ),
        const SizedBox(height: 8),
        Card(child: Column(children: children)),
        const SizedBox(height: 24),
      ],
    );
  }
}

class _ConfigTile extends StatelessWidget {
  final ConfigEntry entry;
  final Future<void> Function(String) onUpdate;
  final Future<void> Function() onReset;

  const _ConfigTile({
    required this.entry,
    required this.onUpdate,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;

    return ListTile(
      leading: Icon(_getIcon(), color: cs.primary),
      title: Row(
        children: [
          Text(entry.key),
          const SizedBox(width: 8),
          _SourceBadge(source: entry.source),
          if (entry.readOnly) ...[
            const SizedBox(width: 4),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
              decoration: BoxDecoration(
                color: cs.errorContainer,
                borderRadius: BorderRadius.circular(4),
              ),
              child: Text(
                'Read Only',
                style: TextStyle(fontSize: 10, color: cs.onErrorContainer),
              ),
            ),
          ],
        ],
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(entry.description),
          const SizedBox(height: 4),
          Text(
            'Current: ${entry.value}',
            style: TextStyle(
              fontFamily: 'monospace',
              fontSize: 12,
              color: cs.onSurface.withValues(alpha: 0.7),
            ),
          ),
          if (entry.envVar.isNotEmpty)
            Text(
              'Env: ${entry.envVar}',
              style: TextStyle(
                fontSize: 10,
                color: cs.onSurface.withValues(alpha: 0.5),
              ),
            ),
        ],
      ),
      trailing: entry.readOnly
          ? null
          : Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                IconButton(
                  icon: const Icon(Icons.edit),
                  onPressed: () => _showEditDialog(context),
                  tooltip: 'Edit',
                ),
                if (entry.source == 'STORE')
                  IconButton(
                    icon: const Icon(Icons.restore),
                    onPressed: onReset,
                    tooltip: 'Reset to default',
                  ),
              ],
            ),
    );
  }

  IconData _getIcon() {
    if (entry.key.startsWith('server.')) return Icons.dns;
    if (entry.key.startsWith('aws.')) return Icons.cloud;
    if (entry.key.startsWith('storage.')) return Icons.storage;
    if (entry.key.startsWith('features.')) return Icons.toggle_on;
    if (entry.key.startsWith('endpoints.')) return Icons.link;
    if (entry.key.startsWith('ports.')) return Icons.numbers;
    return Icons.settings;
  }

  void _showEditDialog(BuildContext context) {
    final controller = TextEditingController(text: entry.value);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Edit ${entry.key}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(entry.description),
            const SizedBox(height: 16),
            TextField(
              controller: controller,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Value',
              ),
              maxLines: null,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              onUpdate(controller.text);
              Navigator.pop(context);
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }
}

class _SourceBadge extends StatelessWidget {
  final String source;

  const _SourceBadge({required this.source});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final (color, bgColor) = switch (source) {
      'STORE' => (cs.primary, cs.primaryContainer),
      'ENV' => (cs.tertiary, cs.tertiaryContainer),
      _ => (cs.onSurface.withValues(alpha: 0.7), cs.surfaceContainerHighest),
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(source, style: TextStyle(fontSize: 10, color: color)),
    );
  }
}
