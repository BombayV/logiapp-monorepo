import 'package:geolocator/geolocator.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:flutter/foundation.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/config/constants.dart';

class LocationService {
  static final LocationService _instance = LocationService._internal();
  factory LocationService() => _instance;
  LocationService._internal();

  final AuthService _authService = AuthService(apiUrl: apiBaseUrl);
  static const String _logPrefix = '[LocationService]';

  // Configuración de ubicación
  static const LocationSettings _locationSettings = LocationSettings(
    accuracy: LocationAccuracy.high,
    distanceFilter: 10, // Mínimo 10 metros de diferencia
    timeLimit: Duration(seconds: 15),
  );

  /// Solicitar permisos y enviar ubicación actual
  Future<bool> requestLocationAndSend() async {
    try {
      // Verificar servicios de ubicación
      if (!await _isLocationServiceEnabled()) {
        debugPrint('$_logPrefix Servicios de ubicación deshabilitados');
        return false;
      }

      // Verificar y solicitar permisos
      if (!await _checkAndRequestPermissions()) {
        debugPrint('$_logPrefix Permisos de ubicación denegados');
        return false;
      }

      // Obtener ubicación actual
      final position = await _getCurrentPosition();
      if (position == null) {
        debugPrint('$_logPrefix No se pudo obtener la ubicación');
        return false;
      }

      // Enviar ubicación al servidor
      final success = await _authService.sendLocation(
        position.latitude,
        position.longitude,
      );

      if (success) {
        debugPrint('$_logPrefix Ubicación enviada: ${position.latitude}, ${position.longitude}');
      } else {
        debugPrint('$_logPrefix Error enviando ubicación al servidor');
      }

      return success;
    } catch (e) {
      debugPrint('$_logPrefix Error: $e');
      return false;
    }
  }

  /// Verificar si los servicios de ubicación están habilitados
  Future<bool> _isLocationServiceEnabled() async {
    try {
      return await Geolocator.isLocationServiceEnabled();
    } catch (e) {
      debugPrint('$_logPrefix Error verificando servicios: $e');
      return false;
    }
  }

  /// Verificar y solicitar permisos de ubicación
  Future<bool> _checkAndRequestPermissions() async {
    try {
      // Verificar permisos de ubicación con geolocator
      LocationPermission permission = await Geolocator.checkPermission();

      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
      }

      if (permission == LocationPermission.deniedForever) {
        debugPrint('$_logPrefix Permisos permanentemente denegados');
        return false;
      }

      if (permission == LocationPermission.denied) {
        debugPrint('$_logPrefix Permisos denegados');
        return false;
      }

      // Para Android, verificar permisos adicionales si es necesario
      if (defaultTargetPlatform == TargetPlatform.android) {
        return await _checkAndroidPermissions();
      }

      return true;
    } catch (e) {
      debugPrint('$_logPrefix Error verificando permisos: $e');
      return false;
    }
  }

  /// Verificar permisos específicos de Android
  Future<bool> _checkAndroidPermissions() async {
    try {
      final locationPermission = await Permission.location.status;

      if (locationPermission.isDenied) {
        final result = await Permission.location.request();
        if (!result.isGranted) {
          return false;
        }
      }

      return true;
    } catch (e) {
      debugPrint('$_logPrefix Error verificando permisos Android: $e');
      return true; // Continuar si hay error con permission_handler
    }
  }

  /// Obtener la posición actual del dispositivo
  Future<Position?> _getCurrentPosition() async {
    try {
      return await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
        timeLimit: const Duration(seconds: 15),
      );
    } catch (e) {
      debugPrint('$_logPrefix Error obteniendo posición: $e');
      return null;
    }
  }

  /// Verificar si los permisos están concedidos
  Future<bool> hasLocationPermission() async {
    try {
      final permission = await Geolocator.checkPermission();
      return permission == LocationPermission.whileInUse ||
             permission == LocationPermission.always;
    } catch (e) {
      debugPrint('$_logPrefix Error verificando permisos: $e');
      return false;
    }
  }

  /// Obtener la distancia entre dos puntos
  double getDistanceBetween(
    double startLatitude,
    double startLongitude,
    double endLatitude,
    double endLongitude,
  ) {
    return Geolocator.distanceBetween(
      startLatitude,
      startLongitude,
      endLatitude,
      endLongitude,
    );
  }
}
