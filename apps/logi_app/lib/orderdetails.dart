import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';
import 'package:url_launcher/url_launcher.dart';

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
  final SecureStorageService _secureStorage = SecureStorageService();
  late final AuthService _authService;

  Map<String, dynamic>? _orderData;
  Map<String, dynamic>? _orderItems;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _authService = AuthService(apiUrl: apiBaseUrl);
    _loadOrderDetails();
  }

  // ========== CARGA DE DATOS ==========

  Future<void> _loadOrderDetails() async {
    try {
      setState(() => _isLoading = true);

      final data = await _authService.getOrderById(widget.orderId);
      final items = await _authService.getOrderItems(widget.orderId);

      setState(() {
        _orderData = data;
        _orderItems = items;
        _isLoading = false;
      });
    } catch (e) {
      debugPrint('[OrderDetails] Error cargando detalles: $e');
      setState(() => _isLoading = false);

      if (mounted) {
        _showErrorSnackBar('Error al cargar detalles: $e');
      }
    }
  }

  // ========== GESTIÓN DE USUARIO ==========

  Future<String> _getUserFullName() async {
    try {
      final userDataString = await _secureStorage.getUserData();
      if (userDataString == null) return '';

      final userData = jsonDecode(userDataString) as Map<String, dynamic>;
      final firstName = userData['profile']?['first_name'] ?? '';
      final lastName = userData['profile']?['last_name'] ?? '';
      return '$firstName $lastName'.trim();
    } catch (e) {
      debugPrint('[OrderDetails] Error obteniendo nombre: $e');
      return '';
    }
  }

  Future<void> _handleLogout() async {
    try {
      final success = await _authService.logout();
      if (success && mounted) {
        Navigator.pushAndRemoveUntil(
          context,
          MaterialPageRoute(builder: (context) => const MyApp()),
          (route) => false,
        );
      }
    } catch (e) {
      debugPrint('[OrderDetails] Error en logout: $e');
      if (mounted) {
        _showErrorSnackBar('Error al cerrar sesión: $e');
      }
    }
  }

  // ========== GESTIÓN DE ESTADO ==========

  Color _getStatusColor(String status) {
    switch (status.toLowerCase()) {
      case 'pending':
        return Colors.orange;
      case 'in_progress':
        return Colors.blue;
      case 'delivered':
      case 'completed':
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
      case 'completed':
        return 'Completado';
      case 'cancelled':
        return 'Cancelado';
      default:
        return status.toUpperCase();
    }
  }

  String _getNextStatus(String currentStatus) {
    switch (currentStatus.toLowerCase()) {
      case 'pending':
        return 'in_progress';
      case 'in_progress':
        return 'completed';
      default:
        return currentStatus;
    }
  }

  String _getNextStatusDisplayName(String currentStatus) {
    switch (currentStatus.toLowerCase()) {
      case 'pending':
        return 'En Progreso';
      case 'in_progress':
        return 'Completado';
      default:
        return currentStatus;
    }
  }

  bool _canUpdateStatus(String currentStatus) {
    return currentStatus.toLowerCase() == 'pending' ||
           currentStatus.toLowerCase() == 'in_progress';
  }

  // ========== ACTUALIZACIÓN DE ESTADO ==========

  Future<void> _updateOrderStatus(String newStatus) async {
    try {
      final result = await _authService.updateOrderStatus(widget.orderId, newStatus);

      if (result != null) {
        setState(() => _orderData = result);

        if (mounted) {
          _showSuccessSnackBar('Estado actualizado a ${_getStatusDisplayName(newStatus)}');
        }
      } else {
        throw Exception('No se pudo actualizar el estado');
      }
    } catch (e) {
      debugPrint('[OrderDetails] Error actualizando estado: $e');
      if (mounted) {
        _showErrorSnackBar('Error al actualizar estado: $e');
      }
    }
  }

  Future<void> _showUpdateStatusDialog() async {
    final currentStatus = _orderData!['status'] as String;
    final nextStatus = _getNextStatus(currentStatus);
    final nextStatusDisplay = _getNextStatusDisplayName(currentStatus);

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirmar actualización'),
        content: Text(
          '¿Estás seguro de que quieres cambiar el estado de esta orden a "$nextStatusDisplay"?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancelar'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.blue,
              foregroundColor: Colors.white,
            ),
            child: const Text('Confirmar'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      await _updateOrderStatus(nextStatus);
    }
  }

  // ========== GESTIÓN DE MAPAS ==========

  Future<void> _openMaps(String address) async {
    final encodedAddress = Uri.encodeComponent(address);

    try {
      final mapUrl = Theme.of(context).platform == TargetPlatform.iOS
          ? 'https://maps.apple.com/?q=$encodedAddress'
          : 'https://www.google.com/maps/search/?api=1&query=$encodedAddress';

      final uri = Uri.parse(mapUrl);

      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      } else {
        throw Exception('No se puede abrir la aplicación de mapas');
      }
    } catch (e) {
      debugPrint('[OrderDetails] Error abriendo mapas: $e');
      if (mounted) {
        _showErrorSnackBar('Error al abrir mapas: $e');
      }
    }
  }

  // ========== UTILIDADES ==========

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

  void _showSuccessSnackBar(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.green,
        duration: const Duration(seconds: 3),
      ),
    );
  }

  // ========== CONSTRUCCIÓN DE UI ==========

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: _buildBody(),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return AppBar(
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
            'Detalles de Orden',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
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
                  child: const Row(
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

    if (_orderData == null) {
      return const Center(
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
      );
    }

    return RefreshIndicator(
      onRefresh: _loadOrderDetails,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildOrderInfoCard(),
            const SizedBox(height: 24),
            _buildItemsCard(),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderInfoCard() {
    return Card(
      elevation: 4,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildOrderHeader(),
            const SizedBox(height: 24),
            _buildAddressSection(),
            const SizedBox(height: 16),
            _buildActionButtons(),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderHeader() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          'Orden #${_orderData!['order_number']}',
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
            color: _getStatusColor(_orderData!['status']),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Text(
            _getStatusDisplayName(_orderData!['status']),
            style: const TextStyle(
              color: Colors.white,
              fontWeight: FontWeight.bold,
              fontSize: 14,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildAddressSection() {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Icon(Icons.location_on, size: 24, color: Colors.red),
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
                _orderData!['delivery_address'],
                style: const TextStyle(fontSize: 18, height: 1.3),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildActionButtons() {
    return Column(
      children: [
        // Botón de mapas
        SizedBox(
          width: double.infinity,
          child: ElevatedButton.icon(
            onPressed: () => _openMaps(_orderData!['delivery_address']),
            icon: const Icon(Icons.directions),
            label: const Text('Ver en Mapas'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.green,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 12),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
        ),

        // Botón de actualizar estado
        if (_canUpdateStatus(_orderData!['status'])) ...[
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: _showUpdateStatusDialog,
              icon: const Icon(Icons.update),
              label: Text('Cambiar a ${_getNextStatusDisplayName(_orderData!['status'])}'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blue,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
          ),
        ],
      ],
    );
  }

  Widget _buildItemsCard() {
    return Card(
      elevation: 4,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildItemsHeader(),
            const SizedBox(height: 16),
            _buildItemsTable(),
            if (_orderItems != null && _orderItems!['total'] != null)
              _buildItemsTotal(),
          ],
        ),
      ),
    );
  }

  Widget _buildItemsHeader() {
    return const Row(
      children: [
        Icon(Icons.shopping_cart, size: 24, color: Colors.blue),
        SizedBox(width: 8),
        Text(
          'Productos a Recolectar',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }

  Widget _buildItemsTable() {
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          _buildTableHeader(),
          _buildTableRows(),
        ],
      ),
    );
  }

  Widget _buildTableHeader() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey.shade100,
        borderRadius: const BorderRadius.only(
          topLeft: Radius.circular(8),
          topRight: Radius.circular(8),
        ),
      ),
      child: const Row(
        children: [
          Expanded(
            flex: 3,
            child: Text(
              'Producto',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
          ),
          Expanded(
            flex: 1,
            child: Text(
              'Cantidad',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
              textAlign: TextAlign.center,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTableRows() {
    if (_orderItems == null || _orderItems!['items'] == null || _orderItems!['items'].isEmpty) {
      return const Padding(
        padding: EdgeInsets.all(20),
        child: Text(
          'No hay productos para recolectar',
          style: TextStyle(
            fontSize: 16,
            color: Colors.grey,
            fontStyle: FontStyle.italic,
          ),
          textAlign: TextAlign.center,
        ),
      );
    }

    return Column(
      children: _orderItems!['items'].map<Widget>((item) {
        return Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            border: Border(
              bottom: BorderSide(
                color: Colors.grey.shade200,
                width: 1,
              ),
            ),
          ),
          child: Row(
            children: [
              Expanded(
                flex: 3,
                child: Text(
                  item['product_name'] ?? '',
                  style: const TextStyle(fontSize: 15),
                ),
              ),
              Expanded(
                flex: 1,
                child: Text(
                  '${item['quantity']}',
                  style: const TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildItemsTotal() {
    return Padding(
      padding: const EdgeInsets.only(top: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Text(
            'Total de productos:',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(
              horizontal: 12,
              vertical: 6,
            ),
            decoration: BoxDecoration(
              color: Colors.blue.shade100,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Text(
              '${_orderItems!['total']}',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: Colors.blue.shade700,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
