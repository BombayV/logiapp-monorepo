import 'dart:async';
import 'package:logi_app/services/location_service.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/config/constants.dart';

class BackgroundLocationService {
  static Timer? _backgroundTimer;
  static const Duration _backgroundInterval = Duration(minutes: 5);

  // Inicializar el servicio (sin WorkManager)
  static Future<void> initialize() async {
    // No necesita inicialización especial
    print('Background location service initialized');
  }

  // Iniciar el envío periódico de localización
  static Future<void> startLocationUpdates() async {
    await stopLocationUpdates(); // Detener cualquier timer existente

    _backgroundTimer = Timer.periodic(_backgroundInterval, (timer) async {
      await _sendLocationIfAuthenticated();
    });

    print('Background location updates started');
  }

  // Detener el envío de localización
  static Future<void> stopLocationUpdates() async {
    _backgroundTimer?.cancel();
    _backgroundTimer = null;
    print('Background location updates stopped');
  }

  // Cancelar todas las tareas
  static Future<void> cancelAllTasks() async {
    await stopLocationUpdates();
  }

  // Verificar si el usuario está autenticado y enviar ubicación
  static Future<void> _sendLocationIfAuthenticated() async {
    try {
      final isAuthenticated = await _isUserAuthenticated();

      if (!isAuthenticated) {
        print('Usuario no autenticado, deteniendo background service');
        await stopLocationUpdates();
        return;
      }

      // Enviar ubicación solo si está autenticado
      await LocationService().requestLocationAndSend();
      print('Ubicación enviada desde background service');
    } catch (e) {
      print('Error en background location service: $e');

      // Si hay error de autenticación, detener el servicio
      if (e.toString().contains('auth') || e.toString().contains('token')) {
        print('Error de autenticación detectado, deteniendo servicio');
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
        print('No hay token disponible');
        return false;
      }

      // Verificar que el token sea válido
      final authService = AuthService(apiUrl: apiBaseUrl);
      final isValid = await authService.isAuthenticated();

      if (!isValid) {
        print('Token inválido o expirado');
        return false;
      }

      return true;
    } catch (e) {
      print('Error verificando autenticación: $e');
      return false;
    }
  }
}
