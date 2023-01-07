import 'package:dio/dio.dart';
import 'package:postapic/src/config.dart';
import 'package:retrofit/retrofit.dart';

import 'client.dart';
import 'models.dart';

part 'rest_client.g.dart';

@RestApi()
abstract class RestApiClient implements ApiClient {
  factory RestApiClient(Dio dio, {String baseUrl = apiBaseUrl}) =>
      _RestApiClient(dio, baseUrl: baseUrl);

  @override
  @GET('/api/posts')
  Future<List<Post>> getPosts({
    @Query('offset') required int offset,
    @Query('limit') required int limit,
  });
}
