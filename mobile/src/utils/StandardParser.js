/**
 * StandardParser.js
 * Implements the logic from pkg/umqr and pkg/mwjson in TypeScript
 * for the React Native application.
 */

export const parseMalawiQR = (rawTagData) => {
    // Simplified TLV Parsing logic from pkg/umqr
    // Real implementation would loop through the string parsing Tags and Lengths

    // Mock logic to show structure:
    if (!rawTagData.startsWith("mw:")) {
        throw new Error("Invalid Standard Prefix");
    }

    // Example expected structure: mw:1.0:TXN:AIRTEL:0999123456:5000:SIG...
    const parts = rawTagData.split(":");

    if (parts.length < 6) {
        throw new Error("Malformed QR Data");
    }

    return {
        version: parts[1],
        type: parts[2], // e.g., TXN
        provider: parts[3], // e.g., AIRTEL_MONEY
        recipient: parts[4],
        amount: parseFloat(parts[5]),
        signature: parts[6] || null,
    };
};

export const formatAirtelNumber = (num) => {
    // MSISDN Normalization (pkg/mwjson/validator.go)
    if (num.startsWith("+265")) return num.replace("+265", "0");
    if (num.startsWith("265")) return num.replace("265", "0");
    return num;
};
