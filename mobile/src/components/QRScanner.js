import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';

// Placeholder for react-native-vision-camera
// Will implement actual camera logic once dependencies are installed.
const QRScanner = ({ onRead, onClose }) => {
    return (
        <View style={styles.container}>
            <View style={styles.cameraPlaceholder}>
                <Text style={styles.placeholderText}>Camera Feed</Text>
                <Text style={styles.placeholderSubText}>(Simulated)</Text>

                <TouchableOpacity
                    style={styles.simulateButton}
                    onPress={() => onRead("mw:1.0:TXN-MOCK:...")}
                >
                    <Text style={styles.simulateButtonText}>Simulate Scan</Text>
                </TouchableOpacity>
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
    simulateButton: {
        backgroundColor: '#E0E0E0',
        padding: 10,
        borderRadius: 5,
    },
    simulateButtonText: {
        color: 'black',
        fontWeight: '600',
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
