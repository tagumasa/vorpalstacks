import 'package:flutter/material.dart';

class Sidebar extends StatelessWidget {
  final int selectedIndex;
  final Function(int) onTap;

  const Sidebar({super.key, required this.selectedIndex, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;

    return Container(
      width: 250,
      decoration: BoxDecoration(
        color: cs.surface,
        border: Border(right: BorderSide(color: cs.outline)),
      ),
      child: Column(
        children: [
          _buildHeader(context),
          const Divider(height: 1),
          Expanded(
            child: ListView(
              padding: const EdgeInsets.symmetric(vertical: 8),
              children: [
                _buildSectionHeader('GENERAL', context),
                _buildNavItem(
                  context,
                  Icons.dashboard_outlined,
                  'Dashboard',
                  0,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('COMPUTE', context),
                _buildNavItem(context, Icons.functions, 'Lambda', 6),
                _buildNavItem(
                  context,
                  Icons.account_tree_outlined,
                  'Step Functions',
                  7,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('DATABASE', context),
                _buildNavItem(
                  context,
                  Icons.table_chart_outlined,
                  'DynamoDB',
                  2,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('STORAGE', context),
                _buildNavItem(context, Icons.folder_outlined, 'S3', 5),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('INTEGRATION', context),
                _buildNavItem(context, Icons.message_outlined, 'SQS', 3),
                _buildNavItem(context, Icons.campaign_outlined, 'SNS', 4),
                _buildNavItem(context, Icons.api_outlined, 'API Gateway', 11),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('SECURITY & IDENTITY', context),
                _buildNavItem(context, Icons.security_outlined, 'IAM', 1),
                _buildNavItem(context, Icons.people_outline, 'Cognito', 8),
                _buildNavItem(
                  context,
                  Icons.key_outlined,
                  'Secrets Manager',
                  9,
                ),
                _buildNavItem(context, Icons.tune, 'SSM Parameters', 10),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('NETWORKING', context),
                _buildNavItem(context, Icons.cloud_outlined, 'CloudFront', 12),
                _buildNavItem(context, Icons.dns_outlined, 'Route 53', 13),
                _buildNavItem(context, Icons.verified_outlined, 'ACM', 14),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('MANAGEMENT', context),
                _buildNavItem(
                  context,
                  Icons.monitor_heart_outlined,
                  'CloudWatch',
                  15,
                ),
                _buildNavItem(
                  context,
                  Icons.query_stats_outlined,
                  'Athena',
                  16,
                ),
                _buildNavItem(
                  context,
                  Icons.history_outlined,
                  'CloudTrail',
                  20,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('SECURITY', context),
                _buildNavItem(context, Icons.vpn_key_outlined, 'KMS', 17),
                _buildNavItem(context, Icons.security_outlined, 'WAF', 27),
                _buildNavItem(context, Icons.shield_outlined, 'WAFv2', 28),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('EVENTS & SCHEDULING', context),
                _buildNavItem(context, Icons.event_outlined, 'EventBridge', 18),
                _buildNavItem(
                  context,
                  Icons.schedule_outlined,
                  'Scheduler',
                  19,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('STREAMING', context),
                _buildNavItem(context, Icons.waves_outlined, 'Kinesis', 22),
                _buildNavItem(context, Icons.email_outlined, 'SES', 23),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('TIMESTREAM', context),
                _buildNavItem(
                  context,
                  Icons.timeline_outlined,
                  'Timestream Query',
                  25,
                ),
                _buildNavItem(
                  context,
                  Icons.edit_note_outlined,
                  'Timestream Write',
                  26,
                ),
                const SizedBox(height: 8),
                const Divider(),
                _buildSectionHeader('IDENTITY', context),
                _buildNavItem(
                  context,
                  Icons.people_outline,
                  'Cognito Identity',
                  21,
                ),
                _buildNavItem(context, Icons.swap_horiz_outlined, 'STS', 24),
              ],
            ),
          ),
          _buildFooter(context),
        ],
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Container(
            width: 32,
            height: 32,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF2196F3), Color(0xFF1976D2)],
              ),
              borderRadius: BorderRadius.circular(6),
            ),
            child: const Icon(
              Icons.cloud_outlined,
              color: Colors.white,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'VorpalStack',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: cs.onSurface,
                  letterSpacing: 0.5,
                ),
              ),
              Text(
                'Local AWS',
                style: TextStyle(
                  fontSize: 10,
                  color: cs.onSurface.withValues(alpha: 0.6),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title, BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
      child: Text(
        title,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.5),
          letterSpacing: 0.8,
        ),
      ),
    );
  }

  Widget _buildNavItem(
    BuildContext context,
    IconData icon,
    String label,
    int index,
  ) {
    final isSelected = selectedIndex == index;
    final cs = Theme.of(context).colorScheme;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: isSelected
            ? cs.primary.withValues(alpha: 0.1)
            : Colors.transparent,
        borderRadius: BorderRadius.circular(6),
        border: isSelected
            ? Border.all(color: cs.primary.withValues(alpha: 0.3))
            : null,
      ),
      child: ListTile(
        dense: true,
        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        leading: Icon(
          icon,
          size: 20,
          color: isSelected ? cs.primary : cs.onSurface.withValues(alpha: 0.7),
        ),
        title: Text(
          label,
          style: TextStyle(
            color: isSelected
                ? cs.primary
                : cs.onSurface.withValues(alpha: 0.85),
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
            fontSize: 14,
          ),
        ),
        onTap: () => onTap(index),
      ),
    );
  }

  Widget _buildFooter(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        border: Border(top: BorderSide(color: cs.outline)),
      ),
      child: Column(
        children: [
          ListTile(
            dense: true,
            contentPadding: const EdgeInsets.symmetric(horizontal: 12),
            leading: Icon(
              Icons.settings_outlined,
              size: 20,
              color: cs.onSurface.withValues(alpha: 0.7),
            ),
            title: Text(
              'Settings',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.85),
              ),
            ),
            onTap: () => onTap(29),
          ),
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: cs.primary.withValues(alpha: 0.05),
              borderRadius: BorderRadius.circular(4),
              border: Border.all(color: cs.outline),
            ),
            child: Row(
              children: [
                Container(
                  width: 8,
                  height: 8,
                  decoration: const BoxDecoration(
                    color: Color(0xFF4CAF50),
                    shape: BoxShape.circle,
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  'localhost:4566',
                  style: TextStyle(
                    fontSize: 11,
                    color: cs.onSurface.withValues(alpha: 0.7),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
