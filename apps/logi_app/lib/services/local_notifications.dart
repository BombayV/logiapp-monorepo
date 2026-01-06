import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter/foundation.dart';

class LocalNotificationsService {
  static final LocalNotificationsService _instance = LocalNotificationsService._();
  factory LocalNotificationsService() => _instance;
  LocalNotificationsService._();

  final FlutterLocalNotificationsPlugin _plugin = FlutterLocalNotificationsPlugin();
  bool _initialized = false;

  static const AndroidNotificationChannel _channel = AndroidNotificationChannel(
    'logiapp_orders_channel',
    'Actualizaciones de órdenes',
    description: 'Notificaciones cuando se asignan nuevas órdenes o cambian de estado',
    importance: Importance.high,
  );

  Future<void> initialize() async {
    if (_initialized) return;

    try {
      const AndroidInitializationSettings androidInit = AndroidInitializationSettings('@mipmap/ic_launcher');
      const InitializationSettings initSettings = InitializationSettings(android: androidInit);

      await _plugin.initialize(initSettings);

      // Solicitar permisos de notificación (Android 13+)
      await _plugin
          .resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()
          ?.requestNotificationsPermission();

      // Crear canal en Android
      await _plugin.resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()
          ?.createNotificationChannel(_channel);

      _initialized = true;
      debugPrint('[LocalNotifications] Inicialización completa');
    } catch (e) {
      debugPrint('[LocalNotifications] Error en inicialización: $e');
      rethrow;
    }
  }

  Future<void> showNewOrders(int count) async {
    if (!_initialized) await initialize();

    try {
      const AndroidNotificationDetails androidDetails = AndroidNotificationDetails(
        'logiapp_orders_channel',
        'Actualizaciones de órdenes',
        channelDescription: 'Notificaciones cuando se asignan nuevas órdenes o cambian de estado',
        importance: Importance.high,
        priority: Priority.high,
        showWhen: true,
      );

      const NotificationDetails details = NotificationDetails(android: androidDetails);

      await _plugin.show(
        1001,
        'Nueva(s) orden(es) asignada(s)',
        count == 1 ? 'Tienes 1 nueva orden' : 'Tienes $count nuevas órdenes',
        details,
      );
      
      debugPrint('[LocalNotifications] Notificación mostrada: $count orden(es)');
    } catch (e) {
      debugPrint('[LocalNotifications] Error mostrando notificación: $e');
    }
  }
}