import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/route53_provider.dart';

class Route53Page extends ConsumerWidget {
  const Route53Page({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final route53State = ref.watch(route53ListProvider);
    final cs = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: cs.surface,
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Route 53',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.w600,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'DNS and domain management',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: route53State.isLoading && route53State.hostedZones.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : ListView.builder(
                      itemCount: route53State.hostedZones.length,
                      itemBuilder: (context, index) {
                        final zone = route53State.hostedZones[index];
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ExpansionTile(
                            leading: Icon(
                              Icons.dns_outlined,
                              color: cs.primary,
                            ),
                            title: Text(zone.name),
                            subtitle: Text(
                              zone.comment ?? '${zone.recordSetCount} records',
                              style: TextStyle(
                                fontSize: 12,
                                color: cs.onSurface.withValues(alpha: 0.6),
                              ),
                            ),
                            trailing: _MetricChip(
                              label: 'Records',
                              value: '${zone.recordSetCount}',
                            ),
                            children: route53State.records
                                .where((r) => r.name.endsWith(zone.name))
                                .map(
                                  (record) => ListTile(
                                    dense: true,
                                    leading: Container(
                                      width: 40,
                                      alignment: Alignment.center,
                                      child: Text(
                                        record.type,
                                        style: TextStyle(
                                          fontSize: 11,
                                          fontWeight: FontWeight.w600,
                                          color: cs.primary,
                                        ),
                                      ),
                                    ),
                                    title: Text(
                                      record.name.replaceAll(zone.name, ''),
                                      style: const TextStyle(fontSize: 13),
                                    ),
                                    subtitle: Text(
                                      record.value.length > 30
                                          ? '${record.value.substring(0, 30)}...'
                                          : record.value,
                                      style: TextStyle(
                                        fontSize: 11,
                                        color: cs.onSurface.withValues(
                                          alpha: 0.5,
                                        ),
                                      ),
                                    ),
                                    trailing: Text(
                                      'TTL: ${record.ttl}',
                                      style: TextStyle(
                                        fontSize: 11,
                                        color: cs.onSurface.withValues(
                                          alpha: 0.5,
                                        ),
                                      ),
                                    ),
                                  ),
                                )
                                .toList(),
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }
}

class _MetricChip extends StatelessWidget {
  final String label;
  final String value;
  const _MetricChip({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.end,
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
