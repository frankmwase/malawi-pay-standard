import React, { useState } from 'react';
import { SafeAreaView, StatusBar, StyleSheet, Alert } from 'react-native';
import HomeScreen from './src/screens/HomeScreen';
import SettingsScreen from './src/screens/SettingsScreen';
import QRScanner from './src/components/QRScanner';

// UI prototype only: there is no native payment bridge or production QR scanner.
export default function App() {
  const [currentScreen, setCurrentScreen] = useState('home');

  const renderScreen = () => {
    switch (currentScreen) {
      case 'home':
        return (
          <HomeScreen
            onScanPress={() => setCurrentScreen('scan')}
            onHistoryPress={() => Alert.alert("History", "Coming Soon")}
            onSettingsPress={() => setCurrentScreen('settings')}
          />
        );
      case 'scan':
        return <QRScanner onClose={() => setCurrentScreen('home')} />;
      case 'settings':
        return (
          <SettingsScreen
            onBack={() => setCurrentScreen('home')}
          />
        );
      default:
        return <HomeScreen />;
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="dark-content" />
      {renderScreen()}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F5F5F5',
  },
});
