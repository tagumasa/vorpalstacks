import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/lambda_provider.dart';
import '../../widgets/common/status_badge.dart';

class LambdaPage extends ConsumerWidget {
  const LambdaPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final lambdaState = ref.watch(lambdaListProvider);
    final lambdaNotifier = ref.read(lambdaListProvider.notifier);
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
                  'Lambda',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: cs.onSurface,
                  ),
                ),
                const Spacer(),
                ElevatedButton.icon(
                  onPressed: lambdaState.isLoading
                      ? null
                      : () async {
                          final result = await showDialog<Map<String, dynamic>>(
                            context: context,
                            builder: (context) => _CreateFunctionDialog(),
                          );
                          if (result != null) {
                            await lambdaNotifier.createFunction(
                              name: result['name'],
                              runtime: result['runtime'],
                              role: result['role'],
                              memorySize: result['memorySize'] ?? 128,
                              timeout: result['timeout'] ?? 30,
                            );
                          }
                        },
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Create Function'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              'Manage your serverless functions',
              style: TextStyle(
                fontSize: 14,
                color: cs.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: lambdaState.isLoading && lambdaState.functions.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : GridView.builder(
                      gridDelegate:
                          const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2,
                            crossAxisSpacing: 16,
                            mainAxisSpacing: 16,
                            childAspectRatio: 2.0,
                          ),
                      itemCount: lambdaState.functions.length,
                      itemBuilder: (context, index) {
                        final function = lambdaState.functions[index];
                        return _buildFunctionCard(
                          context,
                          function,
                          lambdaNotifier,
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFunctionCard(
    BuildContext context,
    LambdaFunction function,
    LambdaListNotifier notifier,
  ) {
    final cs = Theme.of(context).colorScheme;
    final isActive = function.status == 'Active';

    return Card(
      child: InkWell(
        onTap: () => _showEditDialog(context, function, notifier),
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: const Color(0xFFFF9900).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: const Icon(
                      Icons.functions,
                      color: Color(0xFFFF9900),
                      size: 20,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          function.functionName,
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                            color: cs.onSurface,
                          ),
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 2),
                        StatusBadge(
                          label: function.status,
                          type: isActive
                              ? StatusType.success
                              : StatusType.pending,
                        ),
                      ],
                    ),
                  ),
                  PopupMenuButton(
                    icon: Icon(
                      Icons.more_vert,
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                    itemBuilder: (context) => [
                      const PopupMenuItem(
                        value: 'invoke',
                        child: Row(
                          children: [
                            Icon(Icons.play_arrow, size: 18),
                            SizedBox(width: 12),
                            Text('Invoke'),
                          ],
                        ),
                      ),
                      const PopupMenuItem(
                        value: 'edit',
                        child: Row(
                          children: [
                            Icon(Icons.edit_outlined, size: 18),
                            SizedBox(width: 12),
                            Text('Edit'),
                          ],
                        ),
                      ),
                      const PopupMenuItem(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete_outline, size: 18),
                            SizedBox(width: 12),
                            Text('Delete'),
                          ],
                        ),
                      ),
                    ],
                    onSelected: (value) async {
                      if (value == 'invoke') {
                        _showInvokeDialog(context, function, notifier);
                      } else if (value == 'edit') {
                        _showEditDialog(context, function, notifier);
                      } else if (value == 'delete') {
                        final confirmed = await showDialog<bool>(
                          context: context,
                          builder: (context) => AlertDialog(
                            title: const Text('Delete Function'),
                            content: Text(
                              'Are you sure you want to delete "${function.functionName}"?',
                            ),
                            actions: [
                              TextButton(
                                onPressed: () => Navigator.pop(context, false),
                                child: const Text('Cancel'),
                              ),
                              ElevatedButton(
                                onPressed: () => Navigator.pop(context, true),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: Colors.red,
                                ),
                                child: const Text('Delete'),
                              ),
                            ],
                          ),
                        );
                        if (confirmed == true) {
                          await notifier.deleteFunction(function.functionName);
                        }
                      }
                    },
                  ),
                ],
              ),
              const Spacer(),
              Row(
                children: [
                  _buildInfoChip(context, function.runtime, Icons.code),
                  const SizedBox(width: 8),
                  _buildInfoChip(
                    context,
                    '${function.memorySize} MB',
                    Icons.memory,
                  ),
                  const SizedBox(width: 8),
                  _buildInfoChip(
                    context,
                    '${function.timeout}s',
                    Icons.timer_outlined,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showInvokeDialog(
    BuildContext context,
    LambdaFunction function,
    LambdaListNotifier notifier,
  ) {
    showDialog(
      context: context,
      builder: (context) =>
          _InvokeFunctionDialog(function: function, notifier: notifier),
    );
  }

  void _showEditDialog(
    BuildContext context,
    LambdaFunction function,
    LambdaListNotifier notifier,
  ) {
    showDialog(
      context: context,
      builder: (context) =>
          _EditFunctionDialog(function: function, notifier: notifier),
    );
  }

  Widget _buildInfoChip(BuildContext context, String label, IconData icon) {
    final cs = Theme.of(context).colorScheme;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: cs.surface,
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: cs.outline),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: cs.onSurface.withValues(alpha: 0.6)),
          const SizedBox(width: 4),
          Text(
            label,
            style: TextStyle(
              fontSize: 11,
              color: cs.onSurface.withValues(alpha: 0.8),
            ),
          ),
        ],
      ),
    );
  }
}

class _CreateFunctionDialog extends StatefulWidget {
  @override
  State<_CreateFunctionDialog> createState() => _CreateFunctionDialogState();
}

class _CreateFunctionDialogState extends State<_CreateFunctionDialog> {
  final _nameController = TextEditingController();
  final _roleController = TextEditingController(
    text: 'arn:aws:iam::123456789012:role/lambda-role',
  );
  String _runtime = 'python3.11';
  int _memorySize = 128;
  int _timeout = 30;

  final List<String> _runtimes = [
    'python3.11',
    'python3.10',
    'python3.9',
    'nodejs20.x',
    'nodejs18.x',
    'nodejs16.x',
    'java17',
    'java11',
    'go1.x',
    'ruby3.2',
    'dotnet8',
  ];

  final List<int> _memoryOptions = [128, 256, 512, 1024, 2048];
  final List<int> _timeoutOptions = [30, 60, 120, 300, 600];

  @override
  void dispose() {
    _nameController.dispose();
    _roleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Create Function'),
      content: SizedBox(
        width: 400,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Function Name',
                hintText: 'my-function',
              ),
              autofocus: true,
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              initialValue: _runtime,
              decoration: const InputDecoration(labelText: 'Runtime'),
              items: _runtimes
                  .map((r) => DropdownMenuItem(value: r, child: Text(r)))
                  .toList(),
              onChanged: (v) => setState(() => _runtime = v!),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _roleController,
              decoration: const InputDecoration(labelText: 'IAM Role ARN'),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<int>(
                    initialValue: _memorySize,
                    decoration: const InputDecoration(labelText: 'Memory (MB)'),
                    items: _memoryOptions
                        .map(
                          (m) => DropdownMenuItem(value: m, child: Text('$m')),
                        )
                        .toList(),
                    onChanged: (v) => setState(() => _memorySize = v!),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: DropdownButtonFormField<int>(
                    initialValue: _timeout,
                    decoration: const InputDecoration(labelText: 'Timeout (s)'),
                    items: _timeoutOptions
                        .map(
                          (t) => DropdownMenuItem(value: t, child: Text('$t')),
                        )
                        .toList(),
                    onChanged: (v) => setState(() => _timeout = v!),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () {
            if (_nameController.text.isNotEmpty) {
              Navigator.pop(context, {
                'name': _nameController.text,
                'runtime': _runtime,
                'role': _roleController.text,
                'memorySize': _memorySize,
                'timeout': _timeout,
              });
            }
          },
          child: const Text('Create'),
        ),
      ],
    );
  }
}

class _EditFunctionDialog extends StatefulWidget {
  final LambdaFunction function;
  final LambdaListNotifier notifier;

  const _EditFunctionDialog({required this.function, required this.notifier});

  @override
  State<_EditFunctionDialog> createState() => _EditFunctionDialogState();
}

class _EditFunctionDialogState extends State<_EditFunctionDialog> {
  late int _memorySize;
  late int _timeout;
  late String _status;

  final List<int> _memoryOptions = [128, 256, 512, 1024, 2048];
  final List<int> _timeoutOptions = [30, 60, 120, 300, 600];
  final List<String> _statusOptions = ['Active', 'Inactive'];

  @override
  void initState() {
    super.initState();
    _memorySize = widget.function.memorySize;
    _timeout = widget.function.timeout;
    _status = widget.function.status;
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text('Edit ${widget.function.functionName}'),
      content: SizedBox(
        width: 350,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Runtime: ${widget.function.runtime}',
              style: TextStyle(
                color: Theme.of(
                  context,
                ).colorScheme.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<int>(
              initialValue: _memorySize,
              decoration: const InputDecoration(labelText: 'Memory (MB)'),
              items: _memoryOptions
                  .map((m) => DropdownMenuItem(value: m, child: Text('$m')))
                  .toList(),
              onChanged: (v) => setState(() => _memorySize = v!),
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<int>(
              initialValue: _timeout,
              decoration: const InputDecoration(labelText: 'Timeout (s)'),
              items: _timeoutOptions
                  .map((t) => DropdownMenuItem(value: t, child: Text('$t')))
                  .toList(),
              onChanged: (v) => setState(() => _timeout = v!),
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              initialValue: _status,
              decoration: const InputDecoration(labelText: 'Status'),
              items: _statusOptions
                  .map((s) => DropdownMenuItem(value: s, child: Text(s)))
                  .toList(),
              onChanged: (v) => setState(() => _status = v!),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () async {
            final navigator = Navigator.of(context);
            await widget.notifier.updateFunction(
              name: widget.function.functionName,
              memorySize: _memorySize,
              timeout: _timeout,
              status: _status,
            );
            if (mounted) navigator.pop();
          },
          child: const Text('Save'),
        ),
      ],
    );
  }
}

class _InvokeFunctionDialog extends StatefulWidget {
  final LambdaFunction function;
  final LambdaListNotifier notifier;

  const _InvokeFunctionDialog({required this.function, required this.notifier});

  @override
  State<_InvokeFunctionDialog> createState() => _InvokeFunctionDialogState();
}

class _InvokeFunctionDialogState extends State<_InvokeFunctionDialog> {
  final _payloadController = TextEditingController(text: '{"key": "value"}');
  InvocationResult? _result;
  bool _isInvoking = false;

  @override
  void dispose() {
    _payloadController.dispose();
    super.dispose();
  }

  Future<void> _invoke() async {
    setState(() => _isInvoking = true);
    final result = await widget.notifier.invokeFunction(
      widget.function.functionName,
      payload: _payloadController.text,
    );
    setState(() {
      _result = result;
      _isInvoking = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;

    return AlertDialog(
      title: Text('Invoke ${widget.function.functionName}'),
      content: SizedBox(
        width: 500,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Payload (JSON)',
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: cs.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: _payloadController,
              maxLines: 4,
              decoration: InputDecoration(
                border: const OutlineInputBorder(),
                hintText: '{"key": "value"}',
                fillColor: cs.surfaceContainerHighest.withValues(alpha: 0.5),
                filled: true,
              ),
              style: const TextStyle(fontFamily: 'monospace'),
            ),
            const SizedBox(height: 16),
            if (_isInvoking)
              const Center(child: CircularProgressIndicator())
            else if (_result != null) ...[
              Row(
                children: [
                  Icon(
                    _result!.isSuccessful ? Icons.check_circle : Icons.error,
                    color: _result!.isSuccessful ? Colors.green : Colors.red,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'Status: ${_result!.statusCode}',
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: cs.onSurface,
                    ),
                  ),
                  const Spacer(),
                  Text(
                    '${_result!.duration.inMilliseconds}ms',
                    style: TextStyle(
                      color: cs.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                'Response',
                style: TextStyle(
                  fontWeight: FontWeight.w500,
                  color: cs.onSurface,
                ),
              ),
              const SizedBox(height: 8),
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: cs.surfaceContainerHighest.withValues(alpha: 0.5),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  _result!.body,
                  style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
                ),
              ),
            ],
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Close'),
        ),
        ElevatedButton(
          onPressed: _isInvoking ? null : _invoke,
          child: const Text('Invoke'),
        ),
      ],
    );
  }
}
