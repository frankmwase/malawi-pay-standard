import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';

const SettingsScreen = ({ onBack }) => {

    return (
        <View style={styles.container}>
            <Text style={styles.title}>Settings</Text>

            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Prototype only</Text>
                <Text style={styles.helperText}>
                    No payment authorization, secure PIN storage, or QR camera is available. Never enter a real payment PIN in this prototype.
                </Text>
            </View>

            <TouchableOpacity style={styles.backButton} onPress={onBack}>
                <Text style={styles.backButtonText}>Back</Text>
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        padding: 20,
        backgroundColor: '#F5F5F5',
    },
    title: {
        fontSize: 24,
        fontWeight: 'bold',
        marginBottom: 20,
        color: '#212121',
    },
    section: {
        marginBottom: 30,
        backgroundColor: 'white',
        padding: 15,
        borderRadius: 8,
    },
    sectionTitle: {
        fontSize: 18,
        fontWeight: '600',
        marginBottom: 15,
        color: '#424242',
    },
    helperText: {
        fontSize: 14,
        color: '#424242',
        marginTop: 5,
    },
    backButton: {
        marginTop: 20,
        padding: 15,
        alignItems: 'center',
    },
    backButtonText: {
        color: '#757575',
        fontSize: 16,
    },
});

export default SettingsScreen;
