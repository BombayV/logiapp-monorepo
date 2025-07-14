import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'dart:convert';

class SecureStorageService {
  final _storage = const FlutterSecureStorage();

  Future<void> saveToken(String token) async {
    await _storage.write(key: 'auth_token', value: token);
  }

  Future<void> saveUserData(Map<String, dynamic> userData) async {
    await _storage.write(key: 'user_data', value: userData.toString());
  }

  Future<Map<String, dynamic>?> getUserData() async {
    String? userDataString = await _storage.read(key: 'user_data');
    if (userDataString != null) {
      return Map<String, dynamic>.from(jsonDecode(userDataString));
    }
    return null;
  }

  Future<void> deleteUserData() async {
    await _storage.delete(key: 'user_data');
  }


  Future<String?> getToken() async {
    return await _storage.read(key: 'auth_token');
  }

  Future<void> deleteToken() async {
    await _storage.delete(key: 'auth_token');
  }
}