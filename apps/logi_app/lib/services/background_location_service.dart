import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:logi_app/services/location_service.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/config/constants.dart';

class BackgroundLocationService {
  static Timer? _backgroundTimer;
  static const Duration _backgroundInterval = Duration(minutes: 5);
  static const String _logPrefix = '[BackgroundLocationService]';

  // Inicializar el servicio
  static Future<void> initialize() async {
    debugPrint('$_logPrefix Servicio inicializado');
  }

  // Iniciar el envío periódico de localización
  static Future<void> startLocationUpdates() async {
    await stopLocationUpdates();

    _backgroundTimer = Timer.periodic(_backgroundInterval, (_) async {
      await _sendLocationIfAuthenticated();
    });

    debugPrint('$_logPrefix Actualizaciones de ubicación iniciadas');
  }

  // Detener el envío de localización
  static Future<void> stopLocationUpdates() async {
    _backgroundTimer?.cancel();
    _backgroundTimer = null;
    debugPrint('$_logPrefix Actualizaciones de ubicación detenidas');
  }

  // Cancelar todas las tareas
  static Future<void> cancelAllTasks() async {
    await stopLocationUpdates();
  }

  // Verificar autenticación y enviar ubicación
  static Future<void> _sendLocationIfAuthenticated() async {
    try {
      final isAuthenticated = await _isUserAuthenticated();

      if (!isAuthenticated) {
        debugPrint('$_logPrefix Usuario no autenticado, deteniendo servicio');
        await stopLocationUpdates();
        return;
      }

      await LocationService().requestLocationAndSend();
      debugPrint('$_logPrefix Ubicación enviada exitosamente');
    } catch (e) {
      debugPrint('$_logPrefix Error: $e');

      if (_isAuthenticationError(e)) {
        debugPrint('$_logPrefix Error de autenticación detectado, deteniendo servicio');
        await stopLocationUpdates();
      }
    }
  }

  // Verificar si el usuario está autenticado
  static Future<bool> _isUserAuthenticated() async {
    try {
      final secureStorage = SecureStorageService();
      final token = await secureStorage.getToken();

      if (token == null || token.isEmpty) {
        debugPrint('$_logPrefix No hay token disponible');
        return false;
      }

      final authService = AuthService(apiUrl: apiBaseUrl);
      final isValid = await authService.isAuthenticated();

      if (!isValid) {
        debugPrint('$_logPrefix Token inválido o expirado');
        return false;
      }

      return true;
    } catch (e) {
      debugPrint('$_logPrefix Error verificando autenticación: $e');
      return false;
    }
  }

  // Verificar si es un error de autenticación
  static bool _isAuthenticationError(dynamic error) {
    final errorString = error.toString().toLowerCase();
    return errorString.contains('auth') ||
        errorString.contains('token') ||
        errorString.contains('unauthorized');
  }
}
