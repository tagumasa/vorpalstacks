import 'package:flutter/material.dart';

enum StatusType { success, warning, error, info, pending }

class StatusBadge extends StatelessWidget {
  final String label;
  final StatusType type;

  const StatusBadge({super.key, required this.label, required this.type});

  @override
  Widget build(BuildContext context) {
    final colors = _getColors();
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: colors.bg,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: colors.border),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 6,
            height: 6,
            decoration: BoxDecoration(
              color: colors.dot,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            label,
            style: TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.w600,
              color: colors.text,
            ),
          ),
        ],
      ),
    );
  }

  _Colors _getColors() {
    switch (type) {
      case StatusType.success:
        return _Colors(
          const Color(0xFF4CAF50),
          const Color(0xFF2E7D32),
          const Color(0xFF4CAF50),
        );
      case StatusType.warning:
        return _Colors(
          const Color(0xFFFF9800),
          const Color(0xFFE65100),
          const Color(0xFFFF9800),
        );
      case StatusType.error:
        return _Colors(
          const Color(0xFFD32F2F),
          const Color(0xFFC62828),
          const Color(0xFFD32F2F),
        );
      case StatusType.info:
        return _Colors(
          const Color(0xFF2196F3),
          const Color(0xFF1565C0),
          const Color(0xFF2196F3),
        );
      case StatusType.pending:
        return _Colors(
          const Color(0xFF9E9E9E),
          const Color(0xFF616161),
          const Color(0xFF9E9E9E),
        );
    }
  }
}

class _Colors {
  final Color border;
  final Color text;
  final Color dot;
  Color get bg => border.withValues(alpha: 0.1);

  const _Colors(this.border, this.text, this.dot);
}
