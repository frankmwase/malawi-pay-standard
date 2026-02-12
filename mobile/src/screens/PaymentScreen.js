import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, ActivityIndicator, Alert } from 'react-native';

const PaymentScreen = ({ recipient, amount, provider, onConfirm, onCancel }) => {
    const [pin, setPin] = useState('');
    const [isProcessing, setIsProcessing] = useState(false);

    const handleConfirm = () => {
        if (!pin) {
            Alert.alert('Error', 'Please enter your PIN');
            return;
        }
        setIsProcessing(true);
        // Simulate processing delay
        setTimeout(() => {
            onConfirm(pin);
            setIsProcessing(false);
        }, 100);
    };

    return (
        <View style={style.container}>
            <Text style={style.title}>Confirm Payment</Text>

            <View style={style.detailRow}>
                <Text style={style.label}>To:</Text>
                <Text style={style.value}>{recipient}</Text>
            </View>

            <View style={style.detailRow}>
                <Text style={style.label}>Provider:</Text>
                <Text style={style.value}>{provider}</Text>
            </View>

            <View style={style.amountContainer}>
                <Text style={style.currency}>MWK</Text>
                <Text style={style.amount}>{amount}</Text>
            </View>

            <Text style={style.pinLabel}>Enter PIN to Authorize:</Text>
            <TextInput
                style={style.pinInput}
                secureTextEntry
                keyboardType="numeric"
                maxLength={4}
                value={pin}
                onChangeText={setPin}
                placeholder="****"
            />

            <TouchableOpacity
                style={[style.confirmButton, isProcessing && style.disabledButton]}
                onPress={handleConfirm}
                disabled={isProcessing}
            >
                {isProcessing ? (
                    <ActivityIndicator color="white" />
                ) : (
                    <Text style={style.confirmButtonText}>Pay Now</Text>
                )}
            </TouchableOpacity>

            <TouchableOpacity style={style.cancelButton} onPress={onCancel} disabled={isProcessing}>
                <Text style={style.cancelButtonText}>Cancel</Text>
            </TouchableOpacity>
        </View>
    );
};

const style = StyleSheet.create({
    container: {
        flex: 1,
        padding: 20,
        backgroundColor: 'white',
        alignItems: 'center',
        justifyContent: 'center',
    },
    title: {
        fontSize: 22,
        fontWeight: 'bold',
        marginBottom: 30,
        color: '#212121',
    },
    detailRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        width: '100%',
        marginBottom: 10,
        paddingHorizontal: 10,
    },
    label: {
        fontSize: 16,
        color: '#757575',
    },
    value: {
        fontSize: 16,
        fontWeight: '600',
        color: '#212121',
    },
    amountContainer: {
        flexDirection: 'row',
        alignItems: 'flex-start',
        marginVertical: 30,
    },
    currency: {
        fontSize: 18,
        marginTop: 4,
        fontWeight: '600',
        color: '#424242',
        marginRight: 4,
    },
    amount: {
        fontSize: 36,
        fontWeight: 'bold',
        color: '#2E7D32', // Green
    },
    pinLabel: {
        fontSize: 14,
        color: '#757575',
        marginBottom: 10,
        alignSelf: 'flex-start',
        marginLeft: 10,
    },
    pinInput: {
        width: '100%',
        height: 50,
        borderWidth: 1,
        borderColor: '#E0E0E0',
        borderRadius: 8,
        paddingHorizontal: 15,
        fontSize: 18,
        textAlign: 'center',
        marginBottom: 20,
        color: '#212121',
    },
    confirmButton: {
        width: '100%',
        height: 55,
        backgroundColor: '#2E7D32',
        borderRadius: 8,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: 15,
    },
    disabledButton: {
        backgroundColor: '#A5D6A7',
    },
    confirmButtonText: {
        color: 'white',
        fontSize: 18,
        fontWeight: 'bold',
    },
    cancelButton: {
        padding: 15,
    },
    cancelButtonText: {
        color: '#D32F2F',
        fontSize: 16,
    },
});

export default PaymentScreen;
