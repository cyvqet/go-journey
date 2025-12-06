import 'dart:convert';

import 'package:frontend/features/login/model/login_models.dart';
import 'package:http/http.dart' as http;


class AuthApi {
  final String baseUrl;

  AuthApi({required this.baseUrl});

  Future<LoginResponse> login(LoginRequest request) async {
    final uri = Uri.parse('$baseUrl/user/login');
    final response = await http.post(
      uri,
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode(request.toJson()),
    );

    final Map<String, dynamic> body = jsonDecode(response.body) as Map<String, dynamic>;
    return LoginResponse.fromJson(body);
  }
}
