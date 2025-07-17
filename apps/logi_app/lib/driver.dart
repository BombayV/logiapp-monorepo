import 'dart:convert';
import 'dart:async';
import 'package:flutter/material.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/orderdetails.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/services/location_service.dart';
import 'package:logi_app/services/background_location_service.dart';
import 'package:logi_app/models/order.dart';

class DriverPage extends StatefulWidget {
  const DriverPage({super.key});

  @override
  State<DriverPage> createState() => _DriverPageState();
}

class _DriverPageState extends State<DriverPage> with WidgetsBindingObserver {
  final SecureStorageService _secureStorage = SecureStorageService();
  late final AuthService _authService;
  List<Order> _orders = [];
  bool _isLoading = true;
  Timer? _locationTimer;

  static const Duration _locationInterval = Duration(minutes: 5);

  @override
  void initState() {
    super.initState();
    _authService = AuthService(apiUrl: apiBaseUrl);
    WidgetsBinding.instance.addObserver(this);

    _initialize();
  }

  @override
  void dispose() {
    _stopLocationTimer();
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);

    switch (state) {
      case AppLifecycleState.resumed:
        _sendLocationImmediate();
        _startLocationTimer();
        break;
      case AppLifecycleState.paused:
        _stopLocationTimer();
        break;
      default:
        break;
    }
  }

  // Inicialización
  Future<void> _initialize() async {
    // Inicializar el servicio de localización en segundo plano
    await BackgroundLocationService.initialize();

    // Iniciar el servicio de localización en segundo plano
    await BackgroundLocationService.startLocationUpdates();

    // Enviar ubicación inmediatamente
    await _sendLocationImmediate();

    // Iniciar timer para cuando la app esté activa (cada 5 minutos)
    _startLocationTimer();

    // Cargar órdenes
    await _loadOrders();
  }

  // Gestión de localización
  void _startLocationTimer() {
    _stopLocationTimer();
    _locationTimer = Timer.periodic(_locationInterval, (_) => _sendLocationPeriodic());
  }

  void _stopLocationTimer() {
    _locationTimer?.cancel();
    _locationTimer = null;
  }

  Future<void> _sendLocationImmediate() async {
    try {
      await LocationService().requestLocationAndSend();
    } catch (e) {
      debugPrint('Error al enviar ubicación inmediata: $e');
    }
  }

  Future<void> _sendLocationPeriodic() async {
    try {
      debugPrint('Enviando ubicación periódicamente...');
      await LocationService().requestLocationAndSend();
    } catch (e) {
      debugPrint('Error al enviar ubicación periódica: $e');
    }
  }

  // Gestión de órdenes
  Future<void> _loadOrders() async {
    try {
      setState(() => _isLoading = true);

      final ordersData = await _authService.getOrdersByDriver();
      final ordersList = _parseOrdersData(ordersData);

      setState(() {
        _orders = ordersList
            .map((json) => Order.fromJson(json as Map<String, dynamic>))
            .toList();
        _isLoading = false;
      });
    } catch (e) {
      debugPrint('Error al cargar órdenes: $e');
      setState(() {
        _orders = [];
        _isLoading = false;
      });
      _showErrorSnackBar('Error al obtener pedidos: $e');
    }
  }

  List<dynamic> _parseOrdersData(dynamic ordersData) {
    if (ordersData == null) return [];

    if (ordersData is List) return ordersData;

    if (ordersData is Map<String, dynamic> && ordersData['orders'] is List) {
      return ordersData['orders'] as List;
    }

    return [];
  }

  // Gestión de usuario
  Future<String> _getUserFullName() async {
    try {
      final userDataString = await _secureStorage.getUserData();
      if (userDataString == null) return '';

      final userData = jsonDecode(userDataString) as Map<String, dynamic>;
      final firstName = userData['profile']?['first_name'] ?? '';
      final lastName = userData['profile']?['last_name'] ?? '';
      return '$firstName $lastName'.trim();
    } catch (e) {
      debugPrint('Error al obtener nombre del usuario: $e');
      return '';
    }
  }

  Future<void> _handleLogout() async {
    try {
      // Detener el servicio de localización en segundo plano
      await BackgroundLocationService.stopLocationUpdates();

      final success = await _authService.logout();
      if (success && mounted) {
        Navigator.pushAndRemoveUntil(
          context,
          MaterialPageRoute(builder: (context) => const MyApp()),
          (route) => false,
        );
      } else if (mounted) {
        _showErrorSnackBar('Error al cerrar sesión en el servidor');
      }
    } catch (e) {
      debugPrint('Error en logout: $e');
      if (mounted) {
        _showErrorSnackBar('Error de conexión al cerrar sesión: $e');
      }
    }
  }

  // Utilidades
  void _showErrorSnackBar(String message) {
    if (!mounted) return;

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
        duration: const Duration(seconds: 4),
      ),
    );
  }

  void _navigateToOrderDetails(Order order) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => OrderDetailsPage(orderId: order.orderId),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} '
           '${date.hour}:${date.minute.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: RefreshIndicator(
        onRefresh: _loadOrders,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: _buildBody(),
        ),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      automaticallyImplyLeading: false,
      title: Row(
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(10),
            child: Image.asset(
              'assets/logo.png',
              height: 40,
              width: 40,
              errorBuilder: (context, error, stackTrace) =>
                  const Icon(Icons.local_shipping, size: 40),
            ),
          ),
          const SizedBox(width: 10),
          const Text(
            'LogiApp',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
          ),
        ],
      ),
      actions: [
        FutureBuilder<String>(
          future: _getUserFullName(),
          builder: (context, snapshot) {
            final userName = snapshot.data ?? 'Usuario';
            return PopupMenuButton<String>(
              icon: const Icon(Icons.more_vert),
              onSelected: (value) {
                if (value == 'logout') _handleLogout();
              },
              itemBuilder: (context) => [
                PopupMenuItem<String>(
                  value: 'account',
                  child: Row(
                    children: [
                      const Icon(Icons.person),
                      const SizedBox(width: 8),
                      Text(userName),
                    ],
                  ),
                ),
                const PopupMenuItem<String>(
                  value: 'logout',
                  child: Row(
                    children: [
                      Icon(Icons.exit_to_app),
                      SizedBox(width: 8),
                      Text('Cerrar Sesión'),
                    ],
                  ),
                ),
              ],
            );
          },
        ),
      ],
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_orders.isEmpty) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.inbox_outlined, size: 80, color: Colors.grey),
            SizedBox(height: 16),
            Text(
              'No hay órdenes asignadas',
              style: TextStyle(fontSize: 18, color: Colors.grey),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      itemCount: _orders.length,
      itemBuilder: (context, index) => _buildOrderCard(_orders[index]),
    );
  }

  Widget _buildOrderCard(Order order) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 4,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: InkWell(
        onTap: () => _navigateToOrderDetails(order),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Orden #${order.orderNumber}',
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 6,
                    ),
                    decoration: BoxDecoration(
                      color: order.statusColor,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      order.statusDisplayName,
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 12,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Icon(Icons.location_on, size: 16, color: Colors.grey),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      order.address,
                      style: const TextStyle(fontSize: 14),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  const Icon(Icons.access_time, size: 16, color: Colors.grey),
                  const SizedBox(width: 8),
                  Text(
                    'Creado: ${_formatDate(order.createdAt)}',
                    style: const TextStyle(fontSize: 12, color: Colors.grey),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}