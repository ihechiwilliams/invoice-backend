// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package constants

import (
	"errors"
	"fmt"
)

const (
	// ErrorCodeApiError is a ErrorCode of type api_error.
	ErrorCodeApiError ErrorCode = "api_error"
	// ErrorCodeAccountIsNotActive is a ErrorCode of type account_is_not_active.
	ErrorCodeAccountIsNotActive ErrorCode = "account_is_not_active"
	// ErrorCodeAccountIdMissing is a ErrorCode of type account_id_missing.
	ErrorCodeAccountIdMissing ErrorCode = "account_id_missing"
	// ErrorCodeFileUploadError is a ErrorCode of type file_upload_error.
	ErrorCodeFileUploadError ErrorCode = "file_upload_error"
	// ErrorCodeFileDownloadError is a ErrorCode of type file_download_error.
	ErrorCodeFileDownloadError ErrorCode = "file_download_error"
	// ErrorCodeInvalidAccountStatus is a ErrorCode of type invalid_account_status.
	ErrorCodeInvalidAccountStatus ErrorCode = "invalid_account_status"
	// ErrorCodeInvalidProvider is a ErrorCode of type invalid_provider.
	ErrorCodeInvalidProvider ErrorCode = "invalid_provider"
	// ErrorCodeInvalidProviderData is a ErrorCode of type invalid_provider_data.
	ErrorCodeInvalidProviderData ErrorCode = "invalid_provider_data"
	// ErrorCodeInvalidTransactionDate is a ErrorCode of type invalid_transaction_date.
	ErrorCodeInvalidTransactionDate ErrorCode = "invalid_transaction_date"
	// ErrorCodeInvalidTransactionId is a ErrorCode of type invalid_transaction_id.
	ErrorCodeInvalidTransactionId ErrorCode = "invalid_transaction_id"
	// ErrorCodeMarshalingError is a ErrorCode of type marshaling_error.
	ErrorCodeMarshalingError ErrorCode = "marshaling_error"
	// ErrorCodeMonoAccountIsNotActive is a ErrorCode of type mono_account_is_not_active.
	ErrorCodeMonoAccountIsNotActive ErrorCode = "mono_account_is_not_active"
	// ErrorCodeMonoAccountDataMissing is a ErrorCode of type mono_account_data_missing.
	ErrorCodeMonoAccountDataMissing ErrorCode = "mono_account_data_missing"
	// ErrorCodeMonoExchangeTokenDataMissing is a ErrorCode of type mono_exchange_token_data_missing.
	ErrorCodeMonoExchangeTokenDataMissing ErrorCode = "mono_exchange_token_data_missing"
	// ErrorCodeMonoTransactionsDataMissing is a ErrorCode of type mono_transactions_data_missing.
	ErrorCodeMonoTransactionsDataMissing ErrorCode = "mono_transactions_data_missing"
	// ErrorCodeNoAccountFoundForSync is a ErrorCode of type no_account_found_for_sync.
	ErrorCodeNoAccountFoundForSync ErrorCode = "no_account_found_for_sync"
	// ErrorCodeNoProviderFoundForSync is a ErrorCode of type no_provider_found_for_sync.
	ErrorCodeNoProviderFoundForSync ErrorCode = "no_provider_found_for_sync"
	// ErrorCodeTransactionCreationFailed is a ErrorCode of type transaction_creation_failed.
	ErrorCodeTransactionCreationFailed ErrorCode = "transaction_creation_failed"
	// ErrorCodeValidationError is a ErrorCode of type validation_error.
	ErrorCodeValidationError ErrorCode = "validation_error"
)

var ErrInvalidErrorCode = errors.New("not a valid ErrorCode")

// String implements the Stringer interface.
func (x ErrorCode) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x ErrorCode) IsValid() bool {
	_, err := ParseErrorCode(string(x))
	return err == nil
}

var _ErrorCodeValue = map[string]ErrorCode{
	"api_error":                        ErrorCodeApiError,
	"account_is_not_active":            ErrorCodeAccountIsNotActive,
	"account_id_missing":               ErrorCodeAccountIdMissing,
	"file_upload_error":                ErrorCodeFileUploadError,
	"file_download_error":              ErrorCodeFileDownloadError,
	"invalid_account_status":           ErrorCodeInvalidAccountStatus,
	"invalid_provider":                 ErrorCodeInvalidProvider,
	"invalid_provider_data":            ErrorCodeInvalidProviderData,
	"invalid_transaction_date":         ErrorCodeInvalidTransactionDate,
	"invalid_transaction_id":           ErrorCodeInvalidTransactionId,
	"marshaling_error":                 ErrorCodeMarshalingError,
	"mono_account_is_not_active":       ErrorCodeMonoAccountIsNotActive,
	"mono_account_data_missing":        ErrorCodeMonoAccountDataMissing,
	"mono_exchange_token_data_missing": ErrorCodeMonoExchangeTokenDataMissing,
	"mono_transactions_data_missing":   ErrorCodeMonoTransactionsDataMissing,
	"no_account_found_for_sync":        ErrorCodeNoAccountFoundForSync,
	"no_provider_found_for_sync":       ErrorCodeNoProviderFoundForSync,
	"transaction_creation_failed":      ErrorCodeTransactionCreationFailed,
	"validation_error":                 ErrorCodeValidationError,
}

// ParseErrorCode attempts to convert a string to a ErrorCode.
func ParseErrorCode(name string) (ErrorCode, error) {
	if x, ok := _ErrorCodeValue[name]; ok {
		return x, nil
	}
	return ErrorCode(""), fmt.Errorf("%s is %w", name, ErrInvalidErrorCode)
}