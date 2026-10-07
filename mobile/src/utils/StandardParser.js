// Payment parsing is disabled until the EMV-style TLV/CRC parser and recipient
// verification are implemented and independently tested. Never treat the old
// `mw:` demo string as a payment instruction.
export const parseMalawiQR = () => {
    throw new Error('QR payment parsing is not available in this prototype');
};

export const formatAirtelNumber = (num) => {
    // MSISDN Normalization (pkg/mwjson/validator.go)
    if (num.startsWith("+265")) return num.replace("+265", "0");
    if (num.startsWith("265")) return num.replace("265", "0");
    return num;
};
