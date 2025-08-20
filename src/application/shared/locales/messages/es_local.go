package messages

var EsMessages = map[MessageKeysEnum]string{
	"SOMETHING_WENT_WRONG":          "Oh, lo sentimos, algo no salió bien con esta acción!",
	"RESOURCE_NOT_FOUND":            "Recurso no encontrado!",
	"SOME_PARAMETERS_ARE_MISSING":   "Faltan algunos parámetros: %s.",
	"UNKNOWN_RESPONSE_STATUS":       "Estado de respuesta del servidor desconocido.",
	"TOOL_HAS_NOT_BEEN_INITIALIZED": "La herramienta %s no ha sido configurada.",
	"PROCESSING_DATA_CLIENT_ERROR":  "Error procesando la data en el cliente http.",
	"DEPENDENCY_NOT_FOUND":          "No se encontró %s en el contenedor de dependencias.",
	"AUTHORIZATION_REQUIRED":        "La autorización es requerida.",
	"INVALID_USER_OR_PASSWORD":      "Usuario o contraseña incorrecto.",
	"ERROR_CREATING_USER":           "Algo no salió bien registrado el usuario.",

	"USER_WAS_CREATED":               "El usuario fue registrado",
	"USER_WITH_EMAIL_ALREADY_EXISTS": "Un usuario ya ha sido registrado con el email: %s.",
	"USER_LIST_SUCCESS":              "Lista de usuarios obtenida con éxito.",
	"USER_GET_SUCCESS":               "Usuario obtenido con éxito.",
	"USER_UPDATE_SUCCESS":            "Usuario actualizado con éxito.",
	"USER_DELETE_SUCCESS":            "Usuario eliminado con éxito.",
	"INVALID_USER_ID":                "El ID de usuario es inválido.",

	"PASSWORD_REQUIRED":            "La contraseña es requerida.",
	"PASSWORD_IS_SHORT":            "La contraseña es demasiado corta.",
	"PASSWORD_UNDERMINED_STRENGTH": "La fortaleza de la contraseña está comprometida.",
	"PASSWORD_CREATED":             "Contraseña creada con éxito.",

	"AUTHORIZATION_HEADER_MISSING": "Falta el encabezado de autorización.",
	"AUTHORIZATION_HEADER_INVALID": "El encabezado de autorización es inválido.",
	"AUTHORIZATION_TOKEN_EXPIRED":  "El token de autorización ha expirado.",
	"AUTHORIZATION_GENERATED":      "Token de autorización generado con éxito.",
	"INVALID_JWT_TOKEN":            "Token JWT inválido.",

	"INVALID_EMAIL":    "El correo electrónico es inválido.",
	"INVALID_PASSWORD": "La contraseña es inválida.",
	"INVALID_SESSION":  "La sesión es inválida.",

	"APPLICATION_STATUS_OK": "La aplicación está en ejecución.",
}
