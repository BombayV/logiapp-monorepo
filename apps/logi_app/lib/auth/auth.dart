import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';

class AuthService {
  final String _apiUrl;
  final SecureStorageService _secureStorage = SecureStorageService();

  static const String _logPrefix = '[AuthService]';

  AuthService({required String apiUrl}) : _apiUrl = apiUrl;

  // Login del usuario
  Future<Map<String, dynamic>> login(String email, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$_apiUrl/v1/users/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'email': email, 'password': password}),
      );

      final data = jsonDecode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        await _secureStorage.saveToken(data['token'] as String);
        await _fetchUserProfile();

        debugPrint('$_logPrefix Login exitoso para: $email');
        return {
          'success': true,
          'token': data['token'],
          'message': 'Login exitoso',
        };
      } else {
        return {
          'success': false,
          'message': data['error'] ?? 'Error en el login',
        };
      }
    } catch (e) {
      debugPrint('$_logPrefix Error en login: $e');
      return {
        'success': false,
        'message': 'Error de conexión: $e',
      };
    }
  }

  // Obtener perfil del usuario
  Future<void> _fetchUserProfile() async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) throw Exception('No hay token disponible');

      final response = await http.get(
        Uri.parse('$_apiUrl/v1/users/me'),
        headers: _buildHeaders(token),
      );

      if (response.statusCode == 200) {
        final userData = jsonDecode(response.body) as Map<String, dynamic>;

        if (userData['role'] == 'sales') {
          throw Exception('Usuario no autorizado');
        }

        await _secureStorage.saveUserData(userData);
        debugPrint('$_logPrefix Perfil de usuario obtenido exitosamente');
      } else {
        throw Exception('Error al obtener datos del usuario');
      }
    } catch (e) {
      debugPrint('$_logPrefix Error obteniendo perfil: $e');
      rethrow;
    }
  }

  // Enviar ubicación del usuario
  Future<bool> sendLocation(double latitude, double longitude) async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) return false;

      final response = await http.put(
        Uri.parse('$_apiUrl/v1/users/location'),
        headers: _buildHeaders(token),
        body: jsonEncode({
          'latitude': latitude,
          'longitude': longitude,
        }),
      );

      final success = response.statusCode == 200;
      if (success) {
        debugPrint('$_logPrefix Ubicación enviada: $latitude, $longitude');
      } else {
        debugPrint('$_logPrefix Error enviando ubicación: ${response.statusCode}');
      }

      return success;
    } catch (e) {
      debugPrint('$_logPrefix Error enviando ubicación: $e');
      return false;
    }
  }

  // Logout del usuario
  Future<bool> logout() async {
    try {
      final token = await _secureStorage.getToken();

      if (token != null) {
        await http.post(
          Uri.parse('$_apiUrl/v1/users/logout'),
          headers: _buildHeaders(token),
        );
      }

      await _secureStorage.clearAll();
      debugPrint('$_logPrefix Logout exitoso');
      return true;
    } catch (e) {
      debugPrint('$_logPrefix Error en logout: $e');
      await _secureStorage.clearAll();
      return false;
    }
  }

  // Verificar si el usuario está autenticado
  Future<bool> isAuthenticated() async {
    try {
      final token = await _secureStorage.getToken();
      return token != null && token.isNotEmpty;
    } catch (e) {
      debugPrint('$_logPrefix Error verificando autenticación: $e');
      return false;
    }
  }

  // Obtener órdenes del conductor
  Future<dynamic> getOrdersByDriver() async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) return null;

      final userDataString = await _secureStorage.getUserData();
      if (userDataString == null) return null;

      final userData = jsonDecode(userDataString) as Map<String, dynamic>;
      final driverId = userData['user_id'];

      final response = await http.get(
        Uri.parse('$_apiUrl/v1/users/$driverId/orders'),
        headers: _buildHeaders(token),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      }

      debugPrint('$_logPrefix Error obteniendo órdenes: ${response.statusCode}');
      return null;
    } catch (e) {
      debugPrint('$_logPrefix Error obteniendo órdenes: $e');
      return null;
    }
  }

  // Obtener orden por ID
  Future<Map<String, dynamic>?> getOrderById(String orderId) async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) return null;

      final response = await http.get(
        Uri.parse('$_apiUrl/v1/orders/$orderId'),
        headers: _buildHeaders(token),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      }

      debugPrint('$_logPrefix Error obteniendo orden $orderId: ${response.statusCode}');
      return null;
    } catch (e) {
      debugPrint('$_logPrefix Error obteniendo orden: $e');
      return null;
    }
  }

  // Obtener items de una orden
  Future<Map<String, dynamic>?> getOrderItems(String orderId) async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) return null;

      final response = await http.get(
        Uri.parse('$_apiUrl/v1/orders/$orderId/items'),
        headers: _buildHeaders(token),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      }

      debugPrint('$_logPrefix Error obteniendo items de orden $orderId: ${response.statusCode}');
      return null;
    } catch (e) {
      debugPrint('$_logPrefix Error obteniendo items: $e');
      return null;
    }
  }

  // Actualizar estado de una orden
  Future<Map<String, dynamic>?> updateOrderStatus(String orderId, String status) async {
    try {
      final token = await _secureStorage.getToken();
      if (token == null) return null;

      final response = await http.patch(
        Uri.parse('$_apiUrl/v1/orders/$orderId/status'),
        headers: _buildHeaders(token),
        body: jsonEncode({'status': status}),
      );

      if (response.statusCode == 200) {
        debugPrint('$_logPrefix Estado de orden $orderId actualizado a: $status');
        return jsonDecode(response.body) as Map<String, dynamic>;
      }

      debugPrint('$_logPrefix Error actualizando estado: ${response.statusCode}');
      return null;
    } catch (e) {
      debugPrint('$_logPrefix Error actualizando estado: $e');
      return null;
    }
  }

  // Construir headers para las peticiones HTTP
  Map<String, String> _buildHeaders(String token) {
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }
}