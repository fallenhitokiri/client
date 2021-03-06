package stellar1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const (
	KeybaseTransactionIDLen       = 16
	KeybaseTransactionIDSuffix    = 0x30
	KeybaseTransactionIDSuffixHex = "30"
)

func KeybaseTransactionIDFromString(s string) (KeybaseTransactionID, error) {
	if len(s) != hex.EncodedLen(KeybaseTransactionIDLen) {
		return "", fmt.Errorf("bad KeybaseTransactionID %q: must be %d bytes long", s, KeybaseTransactionIDLen)
	}
	suffix := s[len(s)-2:]
	if suffix != KeybaseTransactionIDSuffixHex {
		return "", fmt.Errorf("bad KeybaseTransactionID %q: must end in 0x%x", s, KeybaseTransactionIDSuffix)
	}
	return KeybaseTransactionID(s), nil
}

func (k KeybaseTransactionID) String() string {
	return string(k)
}

func (k KeybaseTransactionID) Eq(b KeybaseTransactionID) bool {
	return k == b
}

func (t TransactionID) String() string {
	return string(t)
}

func (t TransactionID) Eq(b TransactionID) bool {
	return t == b
}

func ToTimeMs(t time.Time) TimeMs {
	// the result of calling UnixNano on the zero Time is undefined.
	// https://golang.org/pkg/time/#Time.UnixNano
	if t.IsZero() {
		return 0
	}
	return TimeMs(t.UnixNano() / 1000000)
}

func FromTimeMs(t TimeMs) time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(0, int64(t)*1000000)
}

func (t TimeMs) Time() time.Time {
	return FromTimeMs(t)
}

func (a AccountID) String() string {
	return string(a)
}

func (a AccountID) Eq(b AccountID) bool {
	return a == b
}

func (a AccountID) IsNil() bool {
	return len(a) == 0
}

func (a AccountID) LossyAbbreviation() string {
	if len(a) != 56 {
		return "[invalid account id]"
	}
	return fmt.Sprintf("%v...%v", a[:2], a[len(a)-4:len(a)])
}

func (s SecretKey) String() string {
	return "[secret key redacted]"
}

func (s SecretKey) SecureNoLogString() string {
	return string(s)
}

// CheckInvariants checks that the bundle satisfies
// 1. No duplicate account IDs
// 2. At most one primary account
func (s Bundle) CheckInvariants() error {
	accountIDs := make(map[AccountID]bool)
	var foundPrimary bool
	for _, entry := range s.Accounts {
		_, found := accountIDs[entry.AccountID]
		if found {
			return fmt.Errorf("duplicate account ID: %v", entry.AccountID)
		}
		accountIDs[entry.AccountID] = true
		if entry.IsPrimary {
			if foundPrimary {
				return errors.New("multiple primary accounts")
			}
			foundPrimary = true
		}
		if entry.Mode == AccountMode_NONE {
			return errors.New("account missing mode")
		}
	}
	if s.Revision < 1 {
		return fmt.Errorf("revision %v < 1", s.Revision)
	}
	return nil
}

func (s Bundle) PrimaryAccount() (BundleEntry, error) {
	for _, entry := range s.Accounts {
		if entry.IsPrimary {
			return entry, nil
		}
	}
	return BundleEntry{}, errors.New("primary stellar account not found")
}

func (a *Asset) IsNativeXLM() bool {
	return a.Type == "native"
}

func AssetNative() Asset {
	return Asset{
		Type:   "native",
		Code:   "",
		Issuer: "",
	}
}

func (t TransactionStatus) ToPaymentStatus() PaymentStatus {
	switch t {
	case TransactionStatus_PENDING:
		return PaymentStatus_PENDING
	case TransactionStatus_SUCCESS:
		return PaymentStatus_COMPLETED
	case TransactionStatus_ERROR_TRANSIENT, TransactionStatus_ERROR_PERMANENT:
		return PaymentStatus_ERROR
	default:
		return PaymentStatus_UNKNOWN
	}

}

func (t TransactionStatus) Details(errMsg string) (status, detail string) {
	switch t {
	case TransactionStatus_PENDING:
		status = "pending"
	case TransactionStatus_SUCCESS:
		status = "completed"
	case TransactionStatus_ERROR_TRANSIENT, TransactionStatus_ERROR_PERMANENT:
		status = "error"
		detail = errMsg
	default:
		status = "unknown"
		detail = errMsg
	}

	return status, detail
}

func NewPaymentLocal(txid TransactionID, ctime TimeMs) *PaymentLocal {
	return &PaymentLocal{
		Id:   txid,
		Time: ctime,
	}
}

func (p *PaymentSummary) ToDetails() *PaymentDetails {
	return &PaymentDetails{
		Summary: *p,
	}
}

func (p *PaymentSummary) TransactionID() (TransactionID, error) {
	t, err := p.Typ()
	if err != nil {
		return "", err
	}

	switch t {
	case PaymentSummaryType_STELLAR:
		s := p.Stellar()
		return s.TxID, nil
	case PaymentSummaryType_DIRECT:
		s := p.Direct()
		return s.TxID, nil
	case PaymentSummaryType_RELAY:
		s := p.Relay()
		return s.TxID, nil
	}

	return "", errors.New("unknown payment summary type")
}

func (d *StellarServerDefinitions) GetCurrencyLocal(code OutsideCurrencyCode) (res CurrencyLocal, ok bool) {
	def, found := d.Currencies[code]
	if found {
		res = CurrencyLocal{
			Description: fmt.Sprintf("%s (%s)", string(code), def.Symbol.Symbol),
			Code:        code,
			Symbol:      def.Symbol.Symbol,
			Name:        def.Name,
		}
		ok = true
	} else {
		res = CurrencyLocal{
			Code: code,
		}
		ok = false
	}
	return res, ok
}

func (c OutsideCurrencyCode) String() string {
	return string(c)
}
