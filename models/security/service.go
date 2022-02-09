package security

import (
	"fmt"
	"log"
	"time"

	"statsv0/tools/custerror"

	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

// User es el usuario logueado
type User struct {
	ID          string   `json:"id"  validate:"required"`
	Name        string   `json:"name"  validate:"required"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"  validate:"required"`
}

// Validate valida si el token es valido
func Validate(token string) (*User, error) {
	// Si esta en cache, retornamos el cache
	//declara found y ok como variables resultado de la funcion cache get token
	//si salio todo ok (aca se ejecuta el if con la variable ok recien inicializada)
	//entonces se ejecuta lo de adentro que es lo mismo, pero con el valor user (clave:valor -> token:user)
	// el nil es el error

	if found, ok := cache.Get(token); ok {
		if user, ok := found.(*User); ok {
			return user, nil
		}
	}
	//si por el contrario no se encuentra en cache el token user enttonces utiliza la funcion get remote token del archivo dao que
	//consulta a authgo si lo encuentra lo guarda en cache
	user, err := getRemoteToken(token)
	if err != nil {
		return nil, custerror.Unauthorized
	}

	// Todo bien, se agrega al cache y se retorna

	cache.Set(token, user, gocache.DefaultExpiration)

	return user, nil
}

// Invalidate invalida un token del cache
func Invalidate(token string) {
	cache.Delete(token[7:])
	log.Output(1, fmt.Sprintf("Token invalidado: %s", token))
}
