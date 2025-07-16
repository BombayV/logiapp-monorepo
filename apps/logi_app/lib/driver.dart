import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/orderdetails.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/services/location_service.dart';

class DriverPage extends StatefulWidget {
  const DriverPage({super.key});

  @override
  State<DriverPage> createState() => _DriverPageState();
}

class _DriverPageState extends State<DriverPage> with WidgetsBindingObserver {
  final SecureStorageService secureStorage = SecureStorageService();
  late final AuthService authService;

  @override
  void initState() {
    super.initState();
    authService = AuthService(apiUrl: apiBaseUrl);

    // Agregar observer para detectar cambios en el estado de la aplicación
    WidgetsBinding.instance.addObserver(this);

    // Enviar ubicación cuando se carga la pantalla
    _sendLocationOnResume();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);

    // Enviar ubicación cuando la aplicación se reanude
    if (state == AppLifecycleState.resumed) {
      _sendLocationOnResume();
    }
  }

  Future<void> _sendLocationOnResume() async {
    try {
      // Enviar ubicación automáticamente cuando se reanude la aplicación
      await LocationService().requestLocationAndSend();
    } catch (e) {
      print('Error al enviar ubicación en reanudación: $e');
    }
  }

  // Future<Object> get

  Future<String> getUserFullName() async {
    try {
      final userDataString = await secureStorage.getUserData();
      if (userDataString == null) return '';

      final userData = jsonDecode(userDataString) as Map<String, dynamic>;
      final nombre = userData['profile']?['first_name'] ?? '';
      final apellido = userData['profile']?['last_name'] ?? '';
      return '$nombre $apellido'.trim();
    } catch (e) {
      print('Error al obtener datos del usuario: $e');
      return '';
    }
  }

  Future<void> _handleLogout(BuildContext context) async {
    try {
      final success = await authService.logout();
      if (success) {
        if (context.mounted) {
          Navigator.pushAndRemoveUntil(
            context,
            MaterialPageRoute(builder: (context) => const MyApp()),
            (route) => false,
          );
        }
      } else {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Error al cerrar sesión en el servidor'),
              backgroundColor: Colors.red,
              duration: Duration(seconds: 4),
            ),
          );
        }
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error de conexión al cerrar sesión: $e'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 4),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        title: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(10),
              child: Image.asset(
                'assets/logo.png',
                height: 32,
                width: 32,
                errorBuilder: (context, error, stackTrace) {
                  return Icon(Icons.local_shipping, size: 32);
                },
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
            future: getUserFullName(),
            builder: (context, snapshot) {
              final userName = snapshot.data ?? 'Usuario';
              return PopupMenuButton<String>(
                icon: const Icon(Icons.more_vert),
                onSelected: (value) {
                  switch (value) {
                    case 'account':
                      // TODO: Implementar navegación a cuenta
                      break;
                    case 'logout':
                      _handleLogout(context);
                      break;
                  }
                },
                itemBuilder: (BuildContext context) => [
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
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            children: [
              _buildOrderCard(context),
              // Agregar más contenido aquí según sea necesario
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildOrderCard(BuildContext context) {
    return Card(
      elevation: 4,
      child: InkWell(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => OrderDetailsPage(
                orderId: '12345',
                status: 'Pendiente',
              ),
            ),
          );
        },
        child: Container(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Pedido #12345',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.orange,
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: const Text(
                      'Pendiente',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              const Text('Toca para ver detalles'),
              const SizedBox(height: 8),
              const Row(
                children: [
                  Icon(Icons.location_on, size: 16),
                  SizedBox(width: 4),
                  Text('Dirección de entrega'),
                ],
              ),
            ],
          ),
        ),

      ),
    );
  }
}