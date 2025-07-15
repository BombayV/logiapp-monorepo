import 'package:http/http.dart';
import 'dart:convert';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';

class AuthService {
  final String apiUrl;
  final SecureStorageService secureStorage = SecureStorageService();

  AuthService({required this.apiUrl});

  Future<Map<String, dynamic>> login(String username, String password) async {
    try {
      final response = await post(
        Uri.parse('$apiUrl/v1/users/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'email': username, 'password': password}),
      );

      final data = jsonDecode(response.body);

      if (response.statusCode == 200) {
        await secureStorage.saveToken(data['token']);
        await me();
        return {
          'success': true,
          'token': data['token'],
          'message': 'Login successful',
        };
      } else {
        return {
          'success': false,
          'message': data['error'] ?? 'Login failed',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': 'Error de conexión: $e',
      };
    }
  }

  Future<void> me() async {
    try {
      final token = await secureStorage.getToken();
      if (token == null) throw Exception('No hay token disponible');

      final response = await get(
        Uri.parse('$apiUrl/v1/users/me'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token'
        },
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if(data['role'] == 'sales') {
          throw Exception('Usuario no autorizado');
        }
        await secureStorage.saveUserData(data);
      } else {
        throw Exception('Error al obtener datos del usuario');
      }
    } catch (e) {
      print('Error en me(): $e');
      rethrow;
    }
  }

  Future <bool> sendLocation(double latitude, double longitude) async {
    try {
      final token = await secureStorage.getToken();
      if (token == null) return false;

      final response = await post(
        Uri.parse('$apiUrl/v1/users/location'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token'
        },
        body: jsonEncode({
          'latitude': latitude,
          'longitude': longitude,
        }),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Error al enviar ubicación: $e');
      return false;
    }
  }

  Future<bool> logout() async {
    try {
      final token = await secureStorage.getToken();
      if (token == null) return false;

      final response = await post(
        Uri.parse('$apiUrl/v1/users/logout'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token'
        },
      );

      // Limpiar datos locales independientemente del resultado del servidor
      await secureStorage.clearAll();

      return response.statusCode == 200;
    } catch (e) {
      print('Error en logout: $e');
      // Limpiar datos locales aunque haya error
      await secureStorage.clearAll();
      return false;
    }
  }

  Future<bool> isAuthenticated() async {
    try {
      final token = await secureStorage.getToken();
      return token != null && token.isNotEmpty;
    } catch (e) {
      return false;
    }
  }
}