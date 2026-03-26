import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/apigateway.pbgrpc.dart';
import 'api_providers.dart' show apiGatewayProvider;

class ApiGateway {
  final String id;
  final String name;
  final String description;
  final String created;
  final String endpoint;

  ApiGateway({
    required this.id,
    required this.name,
    this.description = '',
    this.created = '',
    this.endpoint = '',
  });
}

class ApiGatewayState {
  final List<ApiGateway> apis;
  final bool isLoading;
  final String? error;

  const ApiGatewayState({
    this.apis = const [],
    this.isLoading = false,
    this.error,
  });

  ApiGatewayState copyWith({
    List<ApiGateway>? apis,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return ApiGatewayState(
      apis: apis ?? this.apis,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

final apiGatewayListProvider =
    NotifierProvider<ApiGatewayNotifier, ApiGatewayState>(() {
      return ApiGatewayNotifier();
    });

class ApiGatewayNotifier extends Notifier<ApiGatewayState> {
  @override
  ApiGatewayState build() => const ApiGatewayState();

  APIGatewayServiceClient get _client => ref.read(apiGatewayProvider);

  Future<void> loadAPIs() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.getRestApis(GetRestApisRequest());
      final apis = response.items
          .map(
            (api) => ApiGateway(
              id: api.id,
              name: api.name,
              description: api.description,
              created: api.createddate,
              endpoint: 'https://${api.id}.execute-api.local.amazonaws.com',
            ),
          )
          .toList();
      state = state.copyWith(apis: apis, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createAPI(String name, String description) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.createRestApi(
        CreateRestApiRequest(name: name, description: description),
      );
      await loadAPIs();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}
