package ports

import "errors"

// ... otros errores
var ErrPersonDocumentExists = errors.New("Ya existe una persona con esa identificación")
