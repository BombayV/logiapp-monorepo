import 'package:geolocator/geolocator.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/config/constants.dart';

class LocationService {
  static final LocationService _instance = LocationService._internal();
  factory LocationService() => _instance;
  LocationService._internal();

  final AuthService _authService = AuthService(apiUrl: apiBaseUrl);

  /// Solicita permisos de ubicación y envía la ubicación actual
  Future<bool> requestLocationAndSend() async {
    try {
      // Verificar si el servicio de ubicación está habilitado
      bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        print('Los servicios de ubicación están deshabilitados.');
        return false;
      }

      // Verificar permisos
      LocationPermission permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          print('Los permisos de ubicación fueron denegados.');
          return false;
        }
      }

      if (permission == LocationPermission.deniedForever) {
        print('Los permisos de ubicación están permanentemente denegados.');
        return false;
      }

      // Obtener la ubicación actual
      Position position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
        timeLimit: const Duration(seconds: 10),
      );

      // Enviar la ubicación al servidor
      bool success = await _authService.sendLocation(
        position.latitude, 
        position.longitude
      );

      if (success) {
        print('Ubicación enviada exitosamente: ${position.latitude}, ${position.longitude}');
      } else {
        print('Error al enviar la ubicación al servidor');
      }

      return success;
    } catch (e) {
      print('Error al obtener o enviar la ubicación: $e');
      return false;
    }
  }

  /// Obtiene la ubicación actual sin enviarla
  Future<Position?> getCurrentLocation() async {
    try {
      bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        return null;
      }

      LocationPermission permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          return null;
        }
      }

      if (permission == LocationPermission.deniedForever) {
        return null;
      }

      return await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
        timeLimit: const Duration(seconds: 10),
      );
    } catch (e) {
      print('Error al obtener ubicación: $e');
      return null;
    }
  }

  /// Verifica si los permisos de ubicación están concedidos
  Future<bool> hasLocationPermission() async {
    LocationPermission permission = await Geolocator.checkPermission();
    return permission == LocationPermission.always || 
           permission == LocationPermission.whileInUse;
  }
}
