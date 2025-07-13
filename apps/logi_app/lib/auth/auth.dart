import 'package:http/http.dart';
import 'dart:convert';
import '/secureStorage/flutter_secure_storage.dart';

class AuthService {
  final String apiUrl;
  final SecureStorageService secureStorage = SecureStorageService();
  AuthService({required this.apiUrl});

  Future<bool> login(String username, String password) async {
    final response = await post(
      Uri.parse('$apiUrl/v1/users/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': username, 'password': password}),
    );

    final data = jsonDecode(response.body);
    if (response.statusCode == 200) {
      await secureStorage.saveToken(data['token']);
      return true;
    } else {
      throw Exception('Failed to login: ${data['error']}');
    }
  }

  Future<bool> logout() async {
    final response = await post(
      Uri.parse('$apiUrl/v1/users/logout'),
      headers: {'Content-Type': 'application/json', 'Authorization': 'Bearer ${await secureStorage.getToken()}'},
    );

    if (response.statusCode == 200) {
      await secureStorage.deleteToken();
      return true;
    } else {
      throw Exception('Failed to logout');
    }
  }
}