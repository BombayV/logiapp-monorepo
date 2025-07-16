import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/orderdetails.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/services/location_service.dart';
import 'package:logi_app/models/order.dart';

class DriverPage extends StatefulWidget {
  const DriverPage({super.key});

  @override
  State<DriverPage> createState() => _DriverPageState();
}

class _DriverPageState extends State<DriverPage> with WidgetsBindingObserver {
  final SecureStorageService secureStorage = SecureStorageService();
  late final AuthService authService;
  List<Order> orders = [];
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    authService = AuthService(apiUrl: apiBaseUrl);

    // Agregar observer para detectar cambios en el estado de la aplicación
    WidgetsBinding.instance.addObserver(this);

    // Enviar ubicación cuando se carga la pantalla
    _sendLocationOnResume();

    // Cargar órdenes al inicializar
    getOrdersByDriver();
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

  Future<void> getOrdersByDriver() async {
    try {
      setState(() {
        isLoading = true;
      });

      final ordersData = await authService.getOrdersByDriver();

      if (ordersData != null && ordersData['orders'] is List) {
        final ordersList = ordersData['orders'] as List;
        setState(() {
          orders = ordersList
              .map((orderJson) => Order.fromJson(orderJson as Map<String, dynamic>))
              .toList();
          isLoading = false;
        });
      } else {
        setState(() {
          orders = [];
          isLoading = false;
        });
      }
    } catch (e) {
      print('Error al obtener pedidos del conductor: $e');
      setState(() {
        orders = [];
        isLoading = false;
      });
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error al obtener pedidos: $e'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 4),
          ),
        );
      }
    }
  }

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
                height: 40,
                width: 40,
                errorBuilder: (context, error, stackTrace) {
                  return Icon(Icons.local_shipping, size: 40);
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
      body: RefreshIndicator(
        onRefresh: getOrdersByDriver,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: isLoading
              ? const Center(child: CircularProgressIndicator())
              : orders.isEmpty
                  ? const Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.inbox_outlined,
                            size: 80,
                            color: Colors.grey,
                          ),
                          SizedBox(height: 16),
                          Text(
                            'No hay órdenes asignadas',
                            style: TextStyle(
                              fontSize: 18,
                              color: Colors.grey,
                            ),
                          ),
                        ],
                      ),
                    )
                  : ListView.builder(
                      itemCount: orders.length,
                      itemBuilder: (context, index) {
                        final order = orders[index];
                        return _buildOrderCard(context, order);
                      },
                    ),
        ),
      ),
    );
  }

  Widget _buildOrderCard(BuildContext context, Order order) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 4,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
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
              children: [
                const Icon(Icons.email, size: 16, color: Colors.grey),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    order.email,
                    style: const TextStyle(fontSize: 14),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
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
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => OrderDetailsPage(order: order),
                        ),
                      );
                    },
                    icon: const Icon(Icons.visibility),
                    label: const Text('Ver Detalles'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.blue,
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: order.status == 'delivered' || order.status == 'cancelled'
                        ? null
                        : () {
                            _updateOrderStatus(order);
                          },
                    icon: const Icon(Icons.update),
                    label: const Text('Actualizar'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
  }

  Future<void> _updateOrderStatus(Order order) async {
    // Mostrar dialog para actualizar estado
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Actualizar Orden #${order.orderNumber}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: const Text('En Progreso'),
              leading: Radio(
                value: 'in_progress',
                groupValue: order.status,
                onChanged: (value) {
                  Navigator.of(context).pop();
                  _performStatusUpdate(order, 'in_progress');
                },
              ),
            ),
            ListTile(
              title: const Text('Entregado'),
              leading: Radio(
                value: 'delivered',
                groupValue: order.status,
                onChanged: (value) {
                  Navigator.of(context).pop();
                  _performStatusUpdate(order, 'delivered');
                },
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancelar'),
          ),
        ],
      ),
    );
  }

  Future<void> _performStatusUpdate(Order order, String newStatus) async {
    try {
      // Aquí deberías implementar la llamada al API para actualizar el estado
      // final success = await authService.updateOrderStatus(order.orderId, newStatus);

      // Por ahora, simularemos que la actualización fue exitosa
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Estado actualizado a: ${_getStatusDisplayName(newStatus)}'),
          backgroundColor: Colors.green,
        ),
      );

      // Recargar órdenes después de actualizar
      getOrdersByDriver();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error al actualizar estado: $e'),
          backgroundColor: Colors.red,
        ),
      );
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
}