import 'package:flutter/material.dart';
import 'package:logi_app/driver.dart';
import 'package:logi_app/auth/auth.dart';
import 'package:logi_app/config/constants.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'LogiApp',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // TRY THIS: Try running your application with "flutter run". You'll see
        // the application has a purple toolbar. Then, without quitting the app,
        // try changing the seedColor in the colorScheme below to Colors.green
        // and then invoke "hot reload" (save your changes or press the "hot
        // reload" button in a Flutter-supported IDE, or press "r" if you used
        // the command line to start the app).
        //
        // Notice that the counter didn't reset back to zero; the application
        // state is not lost during the reload. To reset the state, use hot
        // restart instead.
        //
        // This works for code too, not just values: Most code changes can be
        // tested with just a hot reload.
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF74bfc3),
          primary: const Color(0xFF74bfc3),
          secondary: const Color(0xFFc0b5a9),
          background: const Color(0xFFfdfeff),
        ),
        appBarTheme: const AppBarTheme(
          backgroundColor: Color(0xFF0B212D),
          foregroundColor: Color(0xFFfdfeff),
          elevation: 0,
          titleTextStyle: TextStyle(
            color: Color(0xFFfdfeff),
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
          actionsIconTheme: IconThemeData(color: Color(0xFFfdfeff)),
        ),
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {

  final userController = TextEditingController();
  final passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset('assets/logo.png', height: 100),
                const SizedBox(width: 20),
                const Text('LogiApp',
                    style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)
                ),
              ],
            ),

            TextField(
              decoration: const InputDecoration(labelText: 'Usuario'),
              keyboardType: TextInputType.emailAddress,
              controller: userController,
            ),
            TextField(
              decoration: const InputDecoration(labelText: 'Contraseña'),
              obscureText: true,
              controller: passwordController,
            ),
            const SizedBox(height: 20),
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
              onPressed: () async {
                final username = userController.text;
                final password = passwordController.text;

                if (username.isEmpty || password.isEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Por favor, complete todos los campos')),
                  );
                  return;
                }

                AuthService authService = AuthService(apiUrl: apiBaseUrl);

                if(await authService.login(username, password)) {
                  Navigator.pushReplacement(
                    context,
                    MaterialPageRoute(builder: (context) => const driverPage()),
                  );
                } else {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Error al iniciar sesión')),
                  );
                }

              },
              child: const Text('Ingresar'),
            ),
          ],
        ),
      ),
    );
  }
}
