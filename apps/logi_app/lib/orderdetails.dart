import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';

class OrderDetailsPage extends StatefulWidget {
  final String orderId;

  const OrderDetailsPage({
    super.key,
    required this.orderId,
  });

  @override
  State<OrderDetailsPage> createState() => _OrderDetailsPageState();
}

class _OrderDetailsPageState extends State<OrderDetailsPage> {
  final SecureStorageService secureStorage = SecureStorageService();
  late final AuthService authService;
  Map<String, dynamic>? orderData;
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    authService = AuthService(apiUrl: apiBaseUrl);
    _loadOrderDetails();
  }

  Future<void> _loadOrderDetails() async {
    try {
      setState(() {
        isLoading = true;
      });

      final data = await authService.getOrderById(widget.orderId);

      setState(() {
        orderData = data;
        isLoading = false;
      });
    } catch (e) {
      print('Error al cargar detalles de la orden: $e');
      setState(() {
        isLoading = false;
      });
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error al cargar detalles: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<String> getUserFullName() async {
    try {
      final userDataString = await secureStorage.getUserData();
      if (userDataString == null) return '';

      final userData = jsonDecode(userDataString);
      final nombre = userData['profile']?['first_name'] ?? '';
      final apellido = userData['profile']?['last_name'] ?? '';
      return '$nombre $apellido'.trim();
    } catch (e) {
      print('Error al obtener datos del usuario: $e');
      return '';
    }
  }

  Future<void> _handleLogout() async {
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
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error al cerrar sesión: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Color _getStatusColor(String status) {
    switch (status.toLowerCase()) {
      case 'pending':
        return Colors.orange;
      case 'in_progress':
        return Colors.blue;
      case 'delivered':
        return Colors.green;
      case 'cancelled':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  String _getStatusDisplayName(String status) {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'Pendiente';
      case 'in_progress':
        return 'En Progreso';
      case 'delivered':
        return 'Entregado';
      case 'cancelled':
        return 'Cancelado';
      default:
        return status;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(10),
              child: Image.asset(
                'assets/logo.png',
                height: 40,
                width: 40,
                errorBuilder: (context, error, stackTrace) {
                  return Icon(Icons.local_shipping, size: 40);
                },
              ),
            ),
            const SizedBox(width: 10),
            const Text(
              'Detalles de Orden',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
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
                  if (value == 'logout') {
                    _handleLogout();
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
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : orderData == null
              ? const Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(Icons.error_outline, size: 80, color: Colors.red),
                      SizedBox(height: 16),
                      Text(
                        'No se pudieron cargar los detalles',
                        style: TextStyle(fontSize: 18, color: Colors.red),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _loadOrderDetails,
                  child: SingleChildScrollView(
                    physics: const AlwaysScrollableScrollPhysics(),
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Tarjeta principal con la información de la orden
                        Card(
                          elevation: 4,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Padding(
                            padding: const EdgeInsets.all(20),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                // Número de orden y estado
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    Text(
                                      'Orden #${orderData!['order_number']}',
                                      style: const TextStyle(
                                        fontSize: 24,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                    Container(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 16,
                                        vertical: 8,
                                      ),
                                      decoration: BoxDecoration(
                                        color: _getStatusColor(orderData!['status']),
                                        borderRadius: BorderRadius.circular(20),
                                      ),
                                      child: Text(
                                        _getStatusDisplayName(orderData!['status']),
                                        style: const TextStyle(
                                          color: Colors.white,
                                          fontWeight: FontWeight.bold,
                                          fontSize: 14,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                const SizedBox(height: 24),

                                // Dirección de entrega
                                Row(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    const Icon(
                                      Icons.location_on,
                                      size: 24,
                                      color: Colors.red,
                                    ),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          const Text(
                                            'Dirección de Entrega',
                                            style: TextStyle(
                                              fontSize: 16,
                                              fontWeight: FontWeight.bold,
                                              color: Colors.grey,
                                            ),
                                          ),
                                          const SizedBox(height: 4),
                                          Text(
                                            orderData!['delivery_address'],
                                            style: const TextStyle(
                                              fontSize: 18,
                                              height: 1.3,
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
    );
  }
}