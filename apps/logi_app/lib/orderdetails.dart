import 'package:flutter/material.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/auth/auth.dart';

class OrderDetailsPage extends StatelessWidget {
  final String orderId;
  final String status;

  const OrderDetailsPage({super.key, required this.orderId, required this.status});

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
              ),
            ),
            SizedBox(width: 10),
            Text(
              'LogiApp',
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        actions: [
          PopupMenuButton(
              icon: const Icon(Icons.more_vert),
              onSelected: (value){
                switch(value){
                  case 'account':
                    break;
                  case 'logout':
                    AuthService authService = AuthService(apiUrl: apiBaseUrl);
                    authService.logout().then((success) {
                      if (success) {
                        Navigator.pushAndRemoveUntil(
                            context,
                            MaterialPageRoute(builder: (context) => MyApp()),
                                (route) => false
                        );
                      } else {
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(content: Text('Error al cerrar sesión')),
                        );
                      }
                    });
                    break;
                }
              }, itemBuilder: (BuildContext context) => [
            const PopupMenuItem<String>(
              value: 'account',
              child: Text('Adrian Cavero'), // Replace with actual user name when the api is implemented
            ),
            const PopupMenuItem<String>(
              value: 'logout',
              child: Text('Cerrar Sesión'),
            ),
          ]
          ),
        ],
      ),
      body:SingleChildScrollView(
      child: Center(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            children: [
              Container(
                padding: const EdgeInsets.all(16.0),
                margin: const EdgeInsets.only(bottom: 20.0),
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.grey,
                      blurRadius: 6,
                      offset: Offset(0, 2),
                    ),
                  ],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Detalles del Pedido',
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                    SizedBox(height: 20),
                    Text('ID del Pedido: $orderId', style: TextStyle(fontSize: 18)),
                    SizedBox(height: 10),
                    Text('Estado: $status', style: TextStyle(fontSize: 18)),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.all(16.0),
                margin: const EdgeInsets.only(bottom: 20.0),
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.grey.withOpacity(0.2),
                      blurRadius: 6,
                      offset: Offset(0, 2),
                    ),
                  ],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Dirreción de Entrega',
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                    SizedBox(height: 20),
                    Text('Guayacanes', style: TextStyle(fontSize: 18)),
                    SizedBox(height: 10),
                    Text('Ciudad: Rumiñahui', style: TextStyle(fontSize: 18)),
                    SizedBox(height: 10),
                    Text('Estado: Pichincha', style: TextStyle(fontSize: 18)),
                    SizedBox(height: 10),
                    Text('Código Postal: 171101', style: TextStyle(fontSize: 18)),
                    SizedBox(height: 10),
                  ],
                ),
              ),
              // Crear una sección para mostrar los productos del pedido en una tabla con el nombre y la cantidad
              // Sección de productos del pedido
              Container(
                padding: const EdgeInsets.all(16.0),
                margin: const EdgeInsets.only(bottom: 20.0),
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.grey.withOpacity(0.2),
                      blurRadius: 6,
                      offset: Offset(0, 2),
                    ),
                  ],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Productos',
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                    SizedBox(height: 20),
                    Table(
                      border: TableBorder.all(color: Colors.grey.shade300),
                      children: [
                        TableRow(
                          decoration: BoxDecoration(color: Colors.grey.shade200),
                          children: [
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('Nombre', style: TextStyle(fontWeight: FontWeight.bold)),
                            ),
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('Cantidad', style: TextStyle(fontWeight: FontWeight.bold)),
                            ),
                          ],
                        ),
                        // Ejemplo de productos, reemplazar con datos reales
                        TableRow(
                          children: [
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('Producto 1'),
                            ),
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('2'),
                            ),
                          ],
                        ),
                        TableRow(
                          children: [
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('Producto 2'),
                            ),
                            Padding(
                              padding: EdgeInsets.all(8.0),
                              child: Text('1'),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              ElevatedButton(
                style: ButtonStyle(
                  backgroundColor: WidgetStateProperty.all(const Color(0xFF74bfc3)),
                  foregroundColor: WidgetStateProperty.all(const Color(0xFF0b0808)),
                  overlayColor: WidgetStateProperty.resolveWith<Color?>(
                        (Set<WidgetState> states) {
                      if (states.contains(WidgetState.pressed)) {
                        return const Color(0xFF0B212D);
                      }
                      return null;
                    },
                  ),
                  padding: WidgetStateProperty.all(
                    const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                  ),
                  shape: WidgetStateProperty.all(
                    RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20.0),
                    ),
                  ),
                ),
                onPressed: () {
                  // Logic to change the status of the order
                },
                // Change the text if the status is 'Pendiente' or 'En Proceso'
                child: Text(
                  status == 'Pendiente' ? 'Marcar como En Proceso' : 'Marcar como Entregado',
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ]),
        ),
      )
      )
    );
  }
}