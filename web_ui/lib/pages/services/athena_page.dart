import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/athena_provider.dart';

class AthenaPage extends ConsumerWidget {
  const AthenaPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final athenaState = ref.watch(athenaListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Athena',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Serverless interactive query service',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: athenaState.isLoading && athenaState.workgroups.isEmpty
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
                              Tab(text: 'Workgroups'),
                              Tab(text: 'Named Queries'),
                            ],
                          ),
                          const SizedBox(height: 16),
                          Expanded(
                            child: TabBarView(
                              children: [
                                _buildWorkgroups(athenaState.workgroups, cs),
                                _buildNamedQueries(
                                  athenaState.namedQueries,
                                  cs,
                                ),
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

  Widget _buildWorkgroups(List<Workgroup> workgroups, ColorScheme cs) {
    return ListView.builder(
      itemCount: workgroups.length,
      itemBuilder: (context, index) {
        final wg = workgroups[index];
        Color stateColor = wg.state == 'ENABLED'
            ? const Color(0xFF4CAF50)
            : const Color(0xFF9E9E9E);
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: Icon(Icons.workspaces_outline, color: cs.primary),
            title: Text(wg.name),
            subtitle: wg.description != null ? Text(wg.description!) : null,
            trailing: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: stateColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: Text(
                wg.state,
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

  Widget _buildNamedQueries(List<NamedQuery> queries, ColorScheme cs) {
    return ListView.builder(
      itemCount: queries.length,
      itemBuilder: (context, index) {
        final q = queries[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: Icon(Icons.code_outlined, color: cs.primary),
            title: Text(q.name),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Database: ${q.database}'),
                if (q.description != null)
                  Text(
                    q.description!,
                    style: TextStyle(
                      fontSize: 12,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
              ],
            ),
            trailing: Text(
              q.workGroup,
              style: TextStyle(
                fontSize: 12,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
          ),
        );
      },
    );
  }
}
