import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';

// No camera integration is shipped yet. Do not simulate a payment QR in a user build.
const QRScanner = ({ onClose }) => {
    return (
        <View style={styles.container}>
            <View style={styles.cameraPlaceholder}>
                <Text style={styles.placeholderText}>QR scanning unavailable</Text>
                <Text style={styles.placeholderSubText}>This prototype cannot process payments. A verified camera and payment integration are required.</Text>
            </View>

            <TouchableOpacity style={styles.closeButton} onPress={onClose}>
                <Text style={styles.closeButtonText}>Close Camera</Text>
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: 'black',
        justifyContent: 'center',
        alignItems: 'center',
    },
    cameraPlaceholder: {
        width: 300,
        height: 300,
        borderWidth: 2,
        borderColor: '#00E676', // Green scan frame
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#212121',
    },
    placeholderText: {
        color: 'white',
        fontSize: 18,
        marginBottom: 10,
    },
    placeholderSubText: {
        color: '#757575',
        marginBottom: 20,
    },
    closeButton: {
        position: 'absolute',
        bottom: 50,
        padding: 15,
        backgroundColor: 'rgba(255,255,255,0.2)',
        borderRadius: 30,
    },
    closeButtonText: {
        color: 'white',
        fontSize: 16,
    },
});

export default QRScanner;
