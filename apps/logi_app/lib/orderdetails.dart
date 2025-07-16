import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:logi_app/models/order.dart';

class OrderDetailsPage extends StatefulWidget {
  final Order order;

  const OrderDetailsPage({
    super.key,
    required this.order,
  });

  @override
  State<OrderDetailsPage> createState() => _OrderDetailsPageState();
}

class _OrderDetailsPageState extends State<OrderDetailsPage> {
  final SecureStorageService secureStorage = SecureStorageService();
  late final AuthService authService;
  late String currentStatus;
  bool _isUpdatingStatus = false;

  @override
  void initState() {
    super.initState();
    authService = AuthService(apiUrl: apiBaseUrl);
    currentStatus = widget.order.status;
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
        if (mounted) {
          Navigator.pushAndRemoveUntil(
            context,
            MaterialPageRoute(builder: (context) => const MyApp()),
            (route) => false,
          );
        }
      } else {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Error al cerrar sesión')),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    }
  }

  Future<void> _updateOrderStatus() async {
    setState(() {
      _isUpdatingStatus = true;
    });

    try {
      // Aquí deberías implementar la llamada al API para actualizar el estado
      // final success = await authService.updateOrderStatus(widget.order.orderId, newStatus);

      // Simular delay de red
      await Future.delayed(const Duration(seconds: 1));

      setState(() {
        currentStatus = currentStatus == 'Pendiente' ? 'En Proceso' : 'Entregado';
        _isUpdatingStatus = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Estado actualizado a: ${_getStatusDisplayName(currentStatus)}'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      setState(() {
        _isUpdatingStatus = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error al actualizar estado: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  Future<void> _sendCurrentLocation() async {
    try {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Enviando ubicación...'),
          backgroundColor: Colors.blue,
        ),
      );

      // Aquí deberías implementar el envío de ubicación
      // await LocationService().requestLocationAndSend();

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Ubicación enviada correctamente'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error al enviar ubicación: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  Color _getStatusColor() {
    switch (currentStatus) {
      case 'Pendiente':
        return Colors.orange;
      case 'En Proceso':
        return Colors.blue;
      case 'Entregado':
        return Colors.green;
      default:
        return Colors.grey;
    }
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
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
        title: Text('Orden #${widget.order.orderNumber}'),
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
                      _handleLogout();
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
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Información básica de la orden
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Información de la Orden',
                          style: TextStyle(
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
                            color: widget.order.statusColor,
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Text(
                            widget.order.statusDisplayName,
                            style: const TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    _buildInfoRow('ID de Orden', widget.order.orderId),
                    _buildInfoRow('Número de Orden', widget.order.orderNumber),
                    _buildInfoRow('Email del Cliente', widget.order.email),
                    _buildInfoRow('Dirección', widget.order.address),
                    _buildInfoRow('Fecha de Creación', _formatDate(widget.order.createdAt)),
                    _buildInfoRow('Última Actualización', _formatDate(widget.order.updatedAt)),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),

            // Botones de acción
            if (currentStatus != 'delivered' && currentStatus != 'cancelled') ...[
              const Text(
                'Actualizar Estado',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 16),
              _buildStatusUpdateButtons(),
            ],

            const SizedBox(height: 20),

            // Botón de ubicación
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Acciones',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: _sendCurrentLocation,
                        icon: const Icon(Icons.my_location),
                        label: const Text('Enviar Mi Ubicación'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.blue,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 12),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              '$label:',
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                color: Colors.grey,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 16),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusUpdateButtons() {
    return Column(
      children: [
        if (currentStatus == 'pending')
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: _isUpdatingStatus ? null : () => _updateOrderStatus(),
              icon: _isUpdatingStatus
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.play_arrow),
              label: const Text('Marcar como En Progreso'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.orange,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
        if (currentStatus == 'in_progress') ...[
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: _isUpdatingStatus ? null : () => _updateOrderStatus(),
              icon: _isUpdatingStatus
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.check_circle),
              label: const Text('Marcar como Entregado'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.green,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
        ],
      ],
    );
  }
}