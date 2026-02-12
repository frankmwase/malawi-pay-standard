import React, { useState, useEffect } from 'react';
import { NativeModules, SafeAreaView, StatusBar, StyleSheet, Alert } from 'react-native';
const { MalawiPayModule } = NativeModules;
import HomeScreen from './src/screens/HomeScreen';
import PaymentScreen from './src/screens/PaymentScreen';
import SettingsScreen from './src/screens/SettingsScreen';
import QRScanner from './src/components/QRScanner';
import { parseMalawiQR } from './src/utils/StandardParser';

export default function App() {
  const [currentScreen, setCurrentScreen] = useState('home');
  const [txnDetails, setTxnDetails] = useState(null);
  const [activeSims, setActiveSims] = useState([]);

  useEffect(() => {
    // Load available SIMs on startup
    const loadSims = async () => {
      try {
        const sims = await MalawiPayModule.getSimOperators();
        setActiveSims(sims);
        console.log("Detected SIMs:", sims);
      } catch (e) {
        console.error("SIM Detection Failed", e);
      }
    };
    loadSims();
  }, []);

  const handleScan = (data) => {
    try {
      const parsed = parseMalawiQR(data);
      setTxnDetails(parsed);
      setCurrentScreen('payment');
    } catch (e) {
      Alert.alert("Invalid QR", "This code is not a valid Malawi Pay standard code.");
      setCurrentScreen('home');
    }
  };

  const handlePaymentConfirm = async (pin) => {
    try {
      // 1. Intelligent SIM Selection
      const targetProvider = txnDetails.provider; // e.g., AIRTEL_MONEY
      const matchingSim = activeSims.find(sim =>
        (sim.carrierName.toUpperCase().includes("AIRTEL") && targetProvider.includes("AIRTEL")) ||
        (sim.carrierName.toUpperCase().includes("TNM") && targetProvider.includes("TNM"))
      );

      if (!matchingSim) {
        Alert.alert("SIM Error", `Please insert a ${targetProvider} SIM to complete this payment.`);
        return;
      }

      // 2. Execute USSD via Native Bridge
      await MalawiPayModule.dialSession(
        txnDetails.provider,
        txnDetails.recipient,
        txnDetails.amount.toString(),
        pin,
        matchingSim.subscriptionId
      );

      Alert.alert("Success", "Transaction Initiated. Please check the USSD dialog.");
      setCurrentScreen('home');
    } catch (e) {
      Alert.alert("Transaction Failed", e.message);
    }
  };

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
        return (
          <QRScanner
            onRead={handleScan}
            onClose={() => setCurrentScreen('home')}
          />
        );
      case 'payment':
        return (
          <PaymentScreen
            recipient={txnDetails?.recipient}
            amount={txnDetails?.amount}
            provider={txnDetails?.provider}
            onConfirm={handlePaymentConfirm}
            onCancel={() => setCurrentScreen('home')}
          />
        );
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
