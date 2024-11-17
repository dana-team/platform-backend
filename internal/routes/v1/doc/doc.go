package doc

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

const (
	bearerKey    = "bearer"
	basicAuthKey = "basic"
)

// setupSecuritySchemes adds security schemes as components to the Huma config.
func setupSecuritySchemes(config huma.Config) {
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		bearerKey: {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "jwt",
		},
		basicAuthKey: {
			Type:         "http",
			Scheme:       "basic",
			BearerFormat: "Basic Auth",
		},
	}
}

// SetupAPIRegistry sets up the API and Regsitry for Huma-backed docs.
func SetupAPIRegistry(engine *gin.Engine) (huma.API, huma.Registry) {
	config := huma.DefaultConfig("platform-backend", "")
	config.OpenAPI.Info.Description = "This serves as an API Reference for the `platform-backend`.\n\n" +
		"It includes all the API-endpoints supported by the backend, with information on how to use them.\n\n" +
		"Code exists at `https://github.com/dana-team/platform-backend`."

	config.DocsPath = ""
	mapRegistry := huma.NewMapRegistry("#/components/schemas/", SchemaNamer)
	config.OpenAPI.Components.Schemas = mapRegistry
	setupSecuritySchemes(config)

	api := humagin.New(engine, config)

	return api, mapRegistry
}

// deref takes a reflect.Type and returns the underlying type if it is a pointer.
// It recursively dereferences the type until a non-pointer type is reached.
//
// If the provided type is a pointer, the function will follow the pointer chain
// by calling `Elem()` on the type until it encounters a non-pointer type. If the
// type is not a pointer, it is returned as-is.
func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// SchemaNamer generates a schema name for a given type, ensuring uniqueness
// by including the full package path when necessary. The function operates
// by taking the Go type, removing pointer indicators, and transforming the
// type name into a form that is compatible with naming conventions.
//
// If the type is an array (e.g., `[]int`), it converts it to a form that
// replaces `[]` with `List`, so `[]int` becomes `ListInt`. For generic types
// (e.g., `MyType[SubType]`), the brackets are replaced to form a concatenated
// string like `MyTypeSubType`.
//
// If two types share the same name but come from different packages, the
// full package path is added to the schema name to prevent collisions. The
// package path is included as part of the name, replacing slashes (`/`) in
// the path with underscores (`_`).
//
// If the type is unnamed, the provided `hint` string is used as the fallback.
func SchemaNamer(t reflect.Type, hint string) string {
	t = deref(t)

	name := t.Name()
	pkgPath := t.PkgPath()

	if name == "" {
		name = hint
	}

	name = strings.ReplaceAll(name, "[]", "List[")

	result := ""
	for _, part := range strings.FieldsFunc(name, func(r rune) bool {
		return r == '[' || r == ']' || r == '*' || r == ','
	}) {
		if pkgPath != "" {
			result += strings.ReplaceAll(pkgPath, "/", "_") + "_"
		}

		fqn := strings.Split(part, ".")
		base := fqn[len(fqn)-1]

		r, size := utf8.DecodeRuneInString(base)
		result += strings.ToUpper(string(r)) + base[size:]
	}

	name = result

	return name
}
