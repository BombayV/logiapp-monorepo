import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:logi_app/config/constants.dart';
import 'package:logi_app/main.dart';
import 'package:logi_app/orderdetails.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/secureStorage/flutter_secure_storage.dart';

class driverPage extends StatelessWidget {
  final SecureStorageService secureStorage = SecureStorageService();
  driverPage({super.key});

  Future<String> getUserFullName() async {
    final userDataString = await secureStorage.getUserData();
    final userData = jsonDecode(userDataString as String);
    print('User Data: $userData');
    final nombre = userData['profile']?['first_name'] ?? '';
    final apellido = userData['profile']?['last_name'] ?? '';
    return '$nombre $apellido';
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
          FutureBuilder<String>(
            future: getUserFullName(),
            builder: (context, snapshot) {
              final userName = snapshot.data ?? '';
              return PopupMenuButton(
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
                    }
                  }, itemBuilder: (BuildContext context) => [
                PopupMenuItem<String>(
                  value: 'account',
                  child: Text(userName),
                ),
                const PopupMenuItem<String>(
                  value: 'logout',
                  child: Text('Cerrar Sesión'),
                ),
              ]
              );
            },
          ),
        ],
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            children: [
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => OrderDetailsPage(orderId: '12345', status: 'Pendiente',)),
                  );
                },
                child: Container(
                  margin: const EdgeInsets.symmetric(vertical: 8),
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Color(0xFFc0b5a9),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: const [
                      Text('Factura #12345', style: TextStyle(fontSize: 16)),

                      Row(
                        children: [
                          Icon(Icons.access_time, color: Colors.orange),
                          SizedBox(width: 6),
                          Text('Pendiente', style: TextStyle(fontSize: 16, color: Colors.orange)),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => OrderDetailsPage(orderId: '67890', status: 'En progreso',)),
                  );
                },
                child: Container(
                  margin: const EdgeInsets.symmetric(vertical: 8),
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Color(0xFFc0b5a9),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: const [
                      Text('Factura #67890', style: TextStyle(fontSize: 16)),
                      Row(
                        children: [
                          Icon(Icons.work, color: Colors.blue),
                          SizedBox(width: 6),
                          Text('En progreso', style: TextStyle(fontSize: 16, color: Colors.blue)),
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