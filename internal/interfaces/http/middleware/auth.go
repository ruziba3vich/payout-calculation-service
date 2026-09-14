package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	httpx "github.com/ruziba3vich/payout-calculation-service/internal/interfaces/http"
)

const principalKey = "principal"

type Principal struct {
	ID   uuid.UUID
	Role auth.Role
}

type TokenParser interface {
	Parse(token string) (uuid.UUID, auth.Role, error)
}

func Auth(parser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httpx.Error(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.Error(c, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		id, role, err := parser.Parse(parts[1])
		if err != nil {
			httpx.Error(c, http.StatusUnauthorized, errs.Public(err))
			return
		}

		c.Set(principalKey, Principal{ID: id, Role: role})
		c.Next()
	}
}

func RequireRole(roles ...auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := GetPrincipal(c)
		if !ok {
			httpx.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		for _, r := range roles {
			if p.Role == r {
				c.Next()
				return
			}
		}

		httpx.Error(c, http.StatusForbidden, "forbidden")
	}
}

// RequireSelfOrAdmin lets admins through, and couriers only when the :param
// in the path matches their own id.
func RequireSelfOrAdmin(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := GetPrincipal(c)
		if !ok {
			httpx.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		if p.Role == auth.RoleAdmin {
			c.Next()
			return
		}

		if c.Param(param) != p.ID.String() {
			httpx.Error(c, http.StatusForbidden, "forbidden")
			return
		}

		c.Next()
	}
}

func GetPrincipal(c *gin.Context) (Principal, bool) {
	v, ok := c.Get(principalKey)
	if !ok {
		return Principal{}, false
	}
	p, ok := v.(Principal)
	return p, ok
}
