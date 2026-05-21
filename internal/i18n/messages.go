package i18n

// Success messages
const (
	MsgSuccess = "success.ok"
)

// Parameter / validation errors
const (
	MsgParamInvalid = "error.param.invalid"
	MsgIDFormat     = "error.param.id_format"
)

// Auth errors
const (
	MsgAuthNoToken         = "error.auth.no_token"
	MsgAuthBadFormat       = "error.auth.bad_format"
	MsgAuthTokenExpired    = "error.auth.token_expired"
	MsgAuthSessionExpired  = "error.auth.session_expired"
	MsgAuthLoginFailed     = "error.auth.login_failed"
	MsgAuthWrongCreds      = "error.auth.wrong_credentials"
	MsgAuthWrongPassword   = "error.auth.wrong_password"
	MsgAuthAccountLocked   = "error.auth.account_locked"
	MsgAuthAccountDisabled = "error.auth.account_disabled"
	MsgAuthLogoutFailed    = "error.auth.logout_failed"
	MsgAuthGetProfileFail  = "error.auth.get_profile_failed"
)

// Password expiry / permission
const (
	MsgMustChangePassword = "error.auth.must_change_password"
	MsgPasswordExpired    = "error.auth.password_expired"
	MsgForbidden          = "error.auth.forbidden"
)

// User errors
const (
	MsgUserNotFound   = "error.user.not_found"
	MsgUsernameExists = "error.user.username_exists"
	MsgEmailExists    = "error.user.email_exists"
	MsgUserListFailed = "error.user.list_failed"
)

// Role errors
const (
	MsgRoleNotFound   = "error.role.not_found"
	MsgRoleCodeExists = "error.role.code_exists"
	MsgRoleListFailed = "error.role.list_failed"
)

// Menu errors
const (
	MsgMenuNotFound       = "error.menu.not_found"
	MsgMenuParentNotFound = "error.menu.parent_not_found"
	MsgMenuHasChildren    = "error.menu.has_children"
	MsgMenuTreeFailed     = "error.menu.tree_failed"
)

// Store errors
const (
	MsgStoreNotFound     = "error.store.not_found"
	MsgStoreHasTerminals = "error.store.has_terminals"
	MsgStoreListFailed   = "error.store.list_failed"
)

// Terminal errors
const (
	MsgTerminalNotFound      = "error.terminal.not_found"
	MsgTerminalDisabled      = "error.terminal.disabled"
	MsgTerminalNoToken       = "error.terminal.no_token"
	MsgTerminalInvalidToken  = "error.terminal.invalid_token"
	MsgTerminalRotateFailed  = "error.terminal.rotate_failed"
	MsgTerminalListFailed    = "error.terminal.list_failed"
	MsgTerminalInvalidStatus = "error.terminal.invalid_status"
)

// Password validation
const (
	MsgPasswordLength = "error.validation.password_length"
	MsgPasswordUpper  = "error.validation.password_upper"
	MsgPasswordLower  = "error.validation.password_lower"
	MsgPasswordDigit  = "error.validation.password_digit"
)

// General
const (
	MsgEndpointNotFound = "error.general.not_found"
	MsgInternalError    = "error.general.internal_error"
	MsgDBError          = "error.general.db_error"
)

// Log errors
const (
	MsgLogOperationsFailed = "error.log.operations_failed"
	MsgLogTerminalsFailed  = "error.log.terminals_failed"
)

const (
	MsgGetTerminalsStatsFailed = "error.stats.get_terminals_failed"
)
