class StateMachine {
  final String stateMachineArn;
  final String name;
  final String type;
  final String status;
  final DateTime? creationDate;
  final DateTime? lastModified;
  final String? roleArn;
  final int executionsCount;

  StateMachine({
    required this.stateMachineArn,
    required this.name,
    this.type = 'STANDARD',
    this.status = 'ACTIVE',
    this.creationDate,
    this.lastModified,
    this.roleArn,
    this.executionsCount = 0,
  });

  factory StateMachine.fromJson(Map<String, dynamic> json) {
    return StateMachine(
      stateMachineArn: json['stateMachineArn'] ?? '',
      name: json['name'] ?? '',
      type: json['type'] ?? 'STANDARD',
      status: json['status'] ?? 'ACTIVE',
      creationDate: json['creationDate'] != null
          ? DateTime.parse(json['creationDate'])
          : null,
      lastModified: json['lastModified'] != null
          ? DateTime.parse(json['lastModified'])
          : null,
      roleArn: json['roleArn'],
      executionsCount: json['executionsCount'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
    'stateMachineArn': stateMachineArn,
    'name': name,
    'type': type,
    'status': status,
    'creationDate': creationDate?.toIso8601String(),
    'lastModified': lastModified?.toIso8601String(),
    'roleArn': roleArn,
    'executionsCount': executionsCount,
  };
}

class Execution {
  final String executionArn;
  final String name;
  final String status;
  final DateTime? startDate;
  final DateTime? stopDate;
  final String? input;
  final String? output;
  final double? billingDetails;

  Execution({
    required this.executionArn,
    required this.name,
    required this.status,
    this.startDate,
    this.stopDate,
    this.input,
    this.output,
    this.billingDetails,
  });

  factory Execution.fromJson(Map<String, dynamic> json) {
    return Execution(
      executionArn: json['executionArn'] ?? '',
      name: json['name'] ?? '',
      status: json['status'] ?? 'RUNNING',
      startDate: json['startDate'] != null
          ? DateTime.parse(json['startDate'])
          : null,
      stopDate: json['stopDate'] != null
          ? DateTime.parse(json['stopDate'])
          : null,
      input: json['input'],
      output: json['output'],
      billingDetails: json['billingDetails']?.toDouble(),
    );
  }

  Map<String, dynamic> toJson() => {
    'executionArn': executionArn,
    'name': name,
    'status': status,
    'startDate': startDate?.toIso8601String(),
    'stopDate': stopDate?.toIso8601String(),
    'input': input,
    'output': output,
    'billingDetails': billingDetails,
  };
}
