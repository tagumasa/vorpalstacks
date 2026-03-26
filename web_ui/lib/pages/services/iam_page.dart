import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/iam_provider.dart';
import '../../widgets/common/status_badge.dart';

class IAMPage extends ConsumerStatefulWidget {
  const IAMPage({super.key});

  @override
  ConsumerState<IAMPage> createState() => _IAMPageState();
}

class _IAMPageState extends ConsumerState<IAMPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _tabController.addListener(() {
      if (!_tabController.indexIsChanging) {
        ref.read(iamListProvider.notifier).setSelectedTab(_tabController.index);
      }
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final iamState = ref.watch(iamListProvider);
    final iamNotifier = ref.read(iamListProvider.notifier);
    final cs = Theme.of(context).colorScheme;

    if (iamState.selectedTabIndex != _tabController.index) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _tabController.animateTo(iamState.selectedTabIndex);
      });
    }

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
                  'IAM',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: iamState.isLoading
                      ? null
                      : () async {
                          final nameController = TextEditingController();
                          final result = await showDialog<String>(
                            context: context,
                            builder: (context) => AlertDialog(
                              title: const Text('Create User'),
                              content: TextField(
                                controller: nameController,
                                decoration: const InputDecoration(
                                  labelText: 'User Name',
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
                            await iamNotifier.createUser(result);
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create User'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage users, roles, and policies',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            TabBar(
              controller: _tabController,
              labelColor: cs.primary,
              unselectedLabelColor: cs.onSurface.withValues(alpha: 0.6),
              indicatorColor: cs.primary,
              tabs: const [
                Tab(text: 'Users'),
                Tab(text: 'Roles'),
                Tab(text: 'Policies'),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildUserTable(iamState, iamNotifier),
                  _buildRoleTable(iamState, iamNotifier),
                  _buildPolicyTable(iamState),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildUserTable(IAMState iamState, IAMNotifier iamNotifier) {
    final cs = Theme.of(context).colorScheme;

    if (iamState.isLoading && iamState.users.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    return Card(
      child: SingleChildScrollView(
        child: DataTable(
          headingRowColor: WidgetStateProperty.all(cs.surface),
          columns: const [
            DataColumn(label: Text('User Name')),
            DataColumn(label: Text('User ID')),
            DataColumn(label: Text('ARN')),
            DataColumn(label: Text('Created')),
            DataColumn(label: Text('Status')),
            DataColumn(label: Text('Actions')),
          ],
          rows: iamState.users
              .map(
                (user) => DataRow(
                  cells: [
                    DataCell(Text(user.userName)),
                    DataCell(Text(user.userId)),
                    DataCell(Text(user.arn, overflow: TextOverflow.ellipsis)),
                    DataCell(Text(_formatDate(user.created))),
                    DataCell(
                      StatusBadge(
                        label: user.status,
                        type: user.status == 'Active'
                            ? StatusType.success
                            : StatusType.pending,
                      ),
                    ),
                    DataCell(
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          IconButton(
                            icon: const Icon(Icons.edit_outlined, size: 18),
                            onPressed: () {},
                          ),
                          IconButton(
                            icon: const Icon(Icons.delete_outline, size: 18),
                            onPressed: iamState.isLoading
                                ? null
                                : () async {
                                    final confirm = await showDialog<bool>(
                                      context: context,
                                      builder: (context) => AlertDialog(
                                        title: const Text('Delete User'),
                                        content: Text(
                                          'Are you sure you want to delete "${user.userName}"?',
                                        ),
                                        actions: [
                                          TextButton(
                                            onPressed: () =>
                                                Navigator.pop(context, false),
                                            child: const Text('Cancel'),
                                          ),
                                          ElevatedButton(
                                            onPressed: () =>
                                                Navigator.pop(context, true),
                                            child: const Text('Delete'),
                                          ),
                                        ],
                                      ),
                                    );
                                    if (confirm == true) {
                                      await iamNotifier.deleteUser(
                                        user.userName,
                                      );
                                    }
                                  },
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
    );
  }

  Widget _buildRoleTable(IAMState iamState, IAMNotifier iamNotifier) {
    final cs = Theme.of(context).colorScheme;

    if (iamState.isLoading && iamState.roles.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    return Card(
      child: SingleChildScrollView(
        child: DataTable(
          headingRowColor: WidgetStateProperty.all(cs.surface),
          columns: const [
            DataColumn(label: Text('Role Name')),
            DataColumn(label: Text('Role ID')),
            DataColumn(label: Text('ARN')),
            DataColumn(label: Text('Created')),
            DataColumn(label: Text('Actions')),
          ],
          rows: iamState.roles
              .map(
                (role) => DataRow(
                  cells: [
                    DataCell(Text(role.roleName)),
                    DataCell(Text(role.roleId)),
                    DataCell(Text(role.arn, overflow: TextOverflow.ellipsis)),
                    DataCell(Text(_formatDate(role.created))),
                    DataCell(
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          IconButton(
                            icon: const Icon(
                              Icons.visibility_outlined,
                              size: 18,
                            ),
                            onPressed: () {},
                          ),
                          IconButton(
                            icon: const Icon(Icons.delete_outline, size: 18),
                            onPressed: iamState.isLoading
                                ? null
                                : () async {
                                    final confirm = await showDialog<bool>(
                                      context: context,
                                      builder: (context) => AlertDialog(
                                        title: const Text('Delete Role'),
                                        content: Text(
                                          'Are you sure you want to delete "${role.roleName}"?',
                                        ),
                                        actions: [
                                          TextButton(
                                            onPressed: () =>
                                                Navigator.pop(context, false),
                                            child: const Text('Cancel'),
                                          ),
                                          ElevatedButton(
                                            onPressed: () =>
                                                Navigator.pop(context, true),
                                            child: const Text('Delete'),
                                          ),
                                        ],
                                      ),
                                    );
                                    if (confirm == true) {
                                      await iamNotifier.deleteRole(
                                        role.roleName,
                                      );
                                    }
                                  },
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
    );
  }

  Widget _buildPolicyTable(IAMState iamState) {
    final cs = Theme.of(context).colorScheme;

    if (iamState.isLoading && iamState.policies.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    return Card(
      child: SingleChildScrollView(
        child: DataTable(
          headingRowColor: WidgetStateProperty.all(cs.surface),
          columns: const [
            DataColumn(label: Text('Policy Name')),
            DataColumn(label: Text('ARN')),
            DataColumn(label: Text('Type')),
            DataColumn(label: Text('Attached To')),
            DataColumn(label: Text('Actions')),
          ],
          rows: iamState.policies
              .map(
                (policy) => DataRow(
                  cells: [
                    DataCell(Text(policy.policyName)),
                    DataCell(Text(policy.arn, overflow: TextOverflow.ellipsis)),
                    DataCell(Text(policy.type)),
                    DataCell(Text('${policy.attachedCount} entities')),
                    DataCell(
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          IconButton(
                            icon: const Icon(
                              Icons.visibility_outlined,
                              size: 18,
                            ),
                            onPressed: () {},
                          ),
                          IconButton(
                            icon: const Icon(Icons.delete_outline, size: 18),
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
    );
  }

  String _formatDate(DateTime date) {
    return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
  }
}
