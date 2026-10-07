import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';

// Not reachable from App: retained as a read-only prototype layout.
const PaymentScreen = ({ recipient, amount, provider, onCancel }) => {

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

            <Text style={style.pinLabel}>Payments are unavailable in this prototype. Do not enter a real PIN.</Text>

            <TouchableOpacity style={style.cancelButton} onPress={onCancel}>
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
    cancelButton: {
        padding: 15,
    },
    cancelButtonText: {
        color: '#D32F2F',
        fontSize: 16,
    },
});

export default PaymentScreen;
