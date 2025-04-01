package authdomain

// AuthenticatedAccount represents an Account that has succeded an
// authentication flow. Code that needs to check who is performing an operation
// can depend on this type.
//
// At the moment this type merely indicatest that an authentication chack has
// succeeded. But it could hold information regarding which kind of
// authentication mechanism was used, e.g., password, passkey. Was 2FA used,
// etc. It this a revisit from a user with "remember me" enabled.
type AuthenticatedAccount struct{ Account }
