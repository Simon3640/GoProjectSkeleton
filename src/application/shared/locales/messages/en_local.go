package messages

var EnMessages = map[MessageKeysEnum]string{
	"SOMETHING_WENT_WRONG":          "Oh sorry, something went wrong with current action!",
	"RESOURCE_NOT_FOUND":            "Resource not found!",
	"UNAUTHORIZED_RESOURCE":         "You are not authorized to access this resource.",
	"SOME_PARAMETERS_ARE_MISSING":   "Some parameters are missing: %s.",
	"UNKNOWN_RESPONSE_STATUS":       "Response status from server unknown.",
	"TOOL_HAS_NOT_BEEN_INITIALIZED": "The %s tool has not been initialized.",
	"PROCESSING_DATA_CLIENT_ERROR":  "Error processing http client data.",
	"DEPENDENCY_NOT_FOUND":          "%s not found in dependencies container.",
	"AUTHORIZATION_REQUIRED":        "Authorization is required.",
	"INVALID_USER_OR_PASSWORD":      "Invalid user or password.",
	"ERROR_CREATING_USER":           "Create user error.",
	"RESOURCE_EXISTS":               "Resource already exists.",
	"INVALID_DATA":                  "Invalid data provided.",
	"NEW_USER_WELCOME":              "Welcome to our platform, %s!",

	"USER_WAS_CREATED":               "User was created.",
	"USER_WITH_EMAIL_ALREADY_EXISTS": "A user has already registered with the email address: %s.",
	"USER_LIST_SUCCESS":              "User list retrieved s uccessfully.",
	"USER_GET_SUCCESS":               "User retrieved successfully.",
	"USER_UPDATE_SUCCESS":            "User updated successfully.",
	"USER_DELETE_SUCCESS":            "User deleted successfully.",
	"INVALID_USER_ID":                "Invalid user ID.",

	"PASSWORD_REQUIRED":            "Password is required.",
	"PASSWORD_IS_SHORT":            "Password is too short.",
	"PASSWORD_UNDERMINED_STRENGTH": "Password strength is undermined.",
	"PASSWORD_CREATED":             "Password created successfully.",
	"PASSWORD_TOKEN_CREATED":       "Reset password token created successfully.",

	"AUTHORIZATION_HEADER_MISSING": "Authorization header is missing.",
	"AUTHORIZATION_HEADER_INVALID": "Authorization header is invalid.",
	"AUTHORIZATION_TOKEN_EXPIRED":  "Authorization token has expired.",
	"AUTHORIZATION_GENERATED":      "Authorization token generated successfully.",
	"INVALID_JWT_TOKEN":            "Invalid JWT token.",

	"INVALID_EMAIL":    "Invalid email.",
	"INVALID_PASSWORD": "Invalid password.",
	"INVALID_SESSION":  "Invalid session.",

	"APPLICATION_STATUS_OK": "Application is running.",
}
