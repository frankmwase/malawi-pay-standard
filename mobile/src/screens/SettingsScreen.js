import React, { useState } from 'react';
import { View, Text, Switch, StyleSheet, TouchableOpacity, TextInput } from 'react-native';

const SettingsScreen = ({ onBack }) => {
    const [sessionUnique, setSessionUnique] = useState(false);
    const [defaultSim, setDefaultSim] = useState('SIM 1');
    const [airtelPin, setAirtelPin] = useState('');
    const [tnmPin, setTnmPin] = useState('');

    const toggleSwitch = () => setSessionUnique(previousState => !previousState);

    return (
        <View style={styles.container}>
            <Text style={styles.title}>Settings</Text>

            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Security</Text>
                <View style={styles.row}>
                    <Text style={styles.label}>Session Unique PIN</Text>
                    <Switch
                        trackColor={{ false: "#767577", true: "#81b0ff" }}
                        thumbColor={sessionUnique ? "#2196F3" : "#f4f3f4"}
                        onValueChange={toggleSwitch}
                        value={sessionUnique}
                    />
                </View>
                <Text style={styles.helperText}>
                    If enabled, you will be asked for a PIN every time. If disabled, the app uses the stored PIN.
                </Text>
            </View>

            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Stored PINs</Text>
                <TextInput
                    style={styles.input}
                    placeholder="Airtel Money PIN"
                    secureTextEntry
                    keyboardType="numeric"
                    maxLength={4}
                    value={airtelPin}
                    onChangeText={setAirtelPin}
                />
                <TextInput
                    style={styles.input}
                    placeholder="TNM Mpamba PIN"
                    secureTextEntry
                    keyboardType="numeric"
                    maxLength={4}
                    value={tnmPin}
                    onChangeText={setTnmPin}
                />
                <TouchableOpacity style={styles.saveButton}>
                    <Text style={styles.saveButtonText}>Save PINs securely</Text>
                </TouchableOpacity>
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
    row: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: 5,
    },
    label: {
        fontSize: 16,
        color: '#212121',
    },
    helperText: {
        fontSize: 12,
        color: '#757575',
        marginTop: 5,
    },
    input: {
        height: 50,
        borderWidth: 1,
        borderColor: '#E0E0E0',
        borderRadius: 8,
        paddingHorizontal: 15,
        marginBottom: 15,
        fontSize: 16,
    },
    saveButton: {
        backgroundColor: '#2196F3',
        padding: 15,
        borderRadius: 8,
        alignItems: 'center',
    },
    saveButtonText: {
        color: 'white',
        fontWeight: 'bold',
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
