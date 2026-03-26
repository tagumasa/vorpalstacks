import 'package:flutter/material.dart';

class Breadcrumb extends StatelessWidget {
  final List<String> items;
  final void Function(int)? onItemTap;

  const Breadcrumb({super.key, required this.items, this.onItemTap});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final widgets = <Widget>[];

    for (int i = 0; i < items.length; i++) {
      final isLast = i == items.length - 1;
      widgets.add(
        InkWell(
          onTap: isLast ? null : () => onItemTap?.call(i),
          child: Text(
            items[i],
            style: TextStyle(
              fontSize: 12,
              color: isLast ? cs.onSurface : cs.primary,
              fontWeight: i == 0 ? FontWeight.w500 : FontWeight.w400,
            ),
          ),
        ),
      );
      if (!isLast) {
        widgets.add(
          Icon(
            Icons.chevron_right,
            size: 16,
            color: cs.onSurface.withValues(alpha: 0.5),
          ),
        );
      }
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
      decoration: BoxDecoration(
        color: cs.surface,
        border: Border(bottom: BorderSide(color: cs.outline)),
      ),
      child: Row(children: widgets),
    );
  }
}
