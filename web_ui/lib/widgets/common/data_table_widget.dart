import 'package:flutter/material.dart';

class AWSDataTable extends StatelessWidget {
  final List<String> headers;
  final List<List<String>> rows;
  final List<Widget>? actions;
  final String? emptyMessage;
  final bool isLoading;
  final void Function(int)? onRowTap;

  const AWSDataTable({
    super.key,
    required this.headers,
    required this.rows,
    this.actions,
    this.emptyMessage,
    this.isLoading = false,
    this.onRowTap,
  });

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;

    return Container(
      decoration: BoxDecoration(
        color: cs.surface,
        border: Border.all(color: cs.outline),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: cs.surface,
              border: Border(bottom: BorderSide(color: cs.outline)),
            ),
            child: Row(
              children: headers
                  .map(
                    (h) => Expanded(
                      child: Text(
                        h,
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: cs.onSurface.withValues(alpha: 0.8),
                        ),
                      ),
                    ),
                  )
                  .toList(),
            ),
          ),
          if (isLoading)
            const Padding(
              padding: EdgeInsets.all(32),
              child: Center(child: CircularProgressIndicator()),
            )
          else if (rows.isEmpty)
            Padding(
              padding: const EdgeInsets.all(32),
              child: Center(
                child: Text(
                  emptyMessage ?? 'No data available',
                  style: TextStyle(color: cs.onSurface.withValues(alpha: 0.5)),
                ),
              ),
            )
          else
            ListView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: rows.length,
              itemBuilder: (context, index) {
                final row = rows[index];
                return InkWell(
                  onTap: onRowTap != null ? () => onRowTap!(index) : null,
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 12,
                    ),
                    decoration: BoxDecoration(
                      color: index.isEven
                          ? cs.surface
                          : cs.surface.withValues(alpha: 0.3),
                      border: Border(
                        bottom: BorderSide(
                          color: cs.outline.withValues(alpha: 0.5),
                        ),
                      ),
                    ),
                    child: Row(
                      children: row
                          .map(
                            (cell) => Expanded(
                              child: Text(
                                cell,
                                style: TextStyle(
                                  fontSize: 13,
                                  color: cs.onSurface,
                                ),
                              ),
                            ),
                          )
                          .toList(),
                    ),
                  ),
                );
              },
            ),
        ],
      ),
    );
  }
}
