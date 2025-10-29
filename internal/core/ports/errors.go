package ports

import "errors"

// ... otros errores
var ErrPersonDocumentExists = errors.New("ya existe una persona con esa identificación")
