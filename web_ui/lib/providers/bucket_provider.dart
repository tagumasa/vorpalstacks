import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/s3.pbgrpc.dart';
import 'api_providers.dart' show s3Provider;

class S3Bucket {
  final String name;
  final String region;
  final DateTime creationDate;
  final int objectCount;
  final int sizeBytes;

  S3Bucket({
    required this.name,
    this.region = 'us-east-1',
    required this.creationDate,
    this.objectCount = 0,
    this.sizeBytes = 0,
  });
}

class S3Object {
  final String key;
  final int size;
  final DateTime lastModified;
  final String eTag;
  final String storageClass;

  S3Object({
    required this.key,
    required this.size,
    required this.lastModified,
    this.eTag = '',
    this.storageClass = 'STANDARD',
  });
}

class S3State {
  final List<S3Bucket> buckets;
  final S3Bucket? selectedBucket;
  final List<S3Object> objects;
  final bool isLoading;
  final String? error;

  const S3State({
    this.buckets = const [],
    this.selectedBucket,
    this.objects = const [],
    this.isLoading = false,
    this.error,
  });

  S3State copyWith({
    List<S3Bucket>? buckets,
    S3Bucket? selectedBucket,
    bool clearSelectedBucket = false,
    List<S3Object>? objects,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return S3State(
      buckets: buckets ?? this.buckets,
      selectedBucket: clearSelectedBucket
          ? null
          : (selectedBucket ?? this.selectedBucket),
      objects: objects ?? this.objects,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class S3Notifier extends Notifier<S3State> {
  @override
  S3State build() => const S3State();

  S3ServiceClient get _client => ref.read(s3Provider);

  Future<void> loadBuckets() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.listBuckets(ListBucketsRequest());
      final buckets = response.buckets.map((b) {
        DateTime creationDate;
        try {
          creationDate = b.creationdate.isNotEmpty
              ? DateTime.parse(b.creationdate)
              : DateTime.now();
        } catch (_) {
          creationDate = DateTime.now();
        }
        return S3Bucket(
          name: b.name,
          region: b.bucketregion.isNotEmpty ? b.bucketregion : 'us-east-1',
          creationDate: creationDate,
        );
      }).toList();
      state = state.copyWith(buckets: buckets, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  void selectBucket(S3Bucket bucket) {
    state = state.copyWith(selectedBucket: bucket, objects: []);
  }

  Future<void> createBucket(String name, {String region = 'us-east-1'}) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final config = region != 'us-east-1'
          ? CreateBucketConfiguration(
              locationconstraint: _parseLocationConstraint(region),
            )
          : null;
      await _client.createBucket(
        CreateBucketRequest(bucket: name, createbucketconfiguration: config),
      );
      await loadBuckets();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  BucketLocationConstraint _parseLocationConstraint(String region) {
    final name =
        'BUCKET_LOCATION_CONSTRAINT_${region.toUpperCase().replaceAll('-', '_')}';
    for (final value in BucketLocationConstraint.values) {
      if (value.name == name) return value;
    }
    return BucketLocationConstraint.values.first;
  }

  Future<void> deleteBucket(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.deleteBucket(DeleteBucketRequest(bucket: name));
      state = state.copyWith(
        buckets: state.buckets.where((b) => b.name != name).toList(),
        clearSelectedBucket: state.selectedBucket?.name == name,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> loadObjects(String bucketName) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.listObjectsV2(
        ListObjectsV2Request(bucket: bucketName),
      );
      final objects = response.contents.map((obj) {
        DateTime lastModified;
        try {
          lastModified = obj.lastmodified.isNotEmpty
              ? DateTime.parse(obj.lastmodified)
              : DateTime.now();
        } catch (_) {
          lastModified = DateTime.now();
        }
        return S3Object(
          key: obj.key,
          size: obj.size.toInt(),
          lastModified: lastModified,
          eTag: obj.etag,
          storageClass: obj.storageclass.name.replaceAll(
            'OBJECT_STORAGE_CLASS_',
            '',
          ),
        );
      }).toList();
      state = state.copyWith(objects: objects, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final bucketListProvider = NotifierProvider<S3Notifier, S3State>(() {
  return S3Notifier();
});
