import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';

const HomeScreen = ({ onScanPress, onHistoryPress, onSettingsPress }) => {
    return (
        <View style={styles.container}>
            <Text style={styles.title}>Malawi Pay</Text>
            <Text style={styles.subtitle}>Universal Digital Exchange</Text>

            <View style={styles.card}>
                <Text style={styles.balanceLabel}>Wallet Status</Text>
                <Text style={styles.balanceValue}>Ready to Scan</Text>
            </View>

            <TouchableOpacity style={styles.scanButton} onPress={onScanPress}>
                <Text style={styles.scanButtonText}>Scan QR Code</Text>
            </TouchableOpacity>

            <View style={styles.row}>
                <TouchableOpacity style={styles.secondaryButton} onPress={onHistoryPress}>
                    <Text style={styles.secondaryButtonText}>History</Text>
                </TouchableOpacity>

                <TouchableOpacity style={styles.secondaryButton} onPress={onSettingsPress}>
                    <Text style={styles.secondaryButtonText}>Settings</Text>
                </TouchableOpacity>
            </View>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#F5F5F5',
        alignItems: 'center',
        justifyContent: 'center',
        padding: 20,
    },
    title: {
        fontSize: 28,
        fontWeight: 'bold',
        color: '#D32F2F', // Malawian Red
        marginBottom: 5,
    },
    subtitle: {
        fontSize: 16,
        color: '#757575',
        marginBottom: 40,
    },
    card: {
        width: '100%',
        backgroundColor: 'white',
        borderRadius: 12,
        padding: 20,
        marginBottom: 30,
        elevation: 2,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
    },
    balanceLabel: {
        fontSize: 14,
        color: '#9E9E9E',
    },
    balanceValue: {
        fontSize: 24,
        fontWeight: 'bold',
        color: '#212121',
        marginTop: 5,
    },
    scanButton: {
        width: '100%',
        height: 60,
        backgroundColor: '#212121', // Black
        borderRadius: 30,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: 20,
    },
    scanButtonText: {
        color: 'white',
        fontSize: 18,
        fontWeight: 'bold',
    },
    row: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        width: '100%',
    },
    secondaryButton: {
        width: '48%',
        height: 50,
        backgroundColor: 'white',
        borderWidth: 1,
        borderColor: '#E0E0E0',
        borderRadius: 25,
        alignItems: 'center',
        justifyContent: 'center',
    },
    secondaryButtonText: {
        color: '#424242',
        fontWeight: '600',
    },
});

export default HomeScreen;
