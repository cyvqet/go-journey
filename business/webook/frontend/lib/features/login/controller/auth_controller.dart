import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/auth_api.dart';
import '../model/login_models.dart';

part 'auth_controller.g.dart';

@riverpod
class AuthController extends _$AuthController {
  late final AuthApi _api;

  @override
  bool build() {
    _api = AuthApi(baseUrl: 'http://localhost:8080');
    return false;
  }

  Future<bool> login({required String email, required String password}) async {
    state = true;
    try {
      final req = LoginRequest(email: email, password: password);
      final res = await _api.login(req);
      final msg = res.message ?? '';
      return msg == '登陆成功' || msg == '登录成功';
    } finally {
      state = false;
    }
  }
}
