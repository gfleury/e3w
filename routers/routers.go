package routers

import (
	"errors"
	"net/http"

	"github.com/coreos/etcd/clientv3"
	"github.com/gfleury/e3w/e3ch"
	"github.com/gin-gonic/gin"
	"github.com/soyking/e3ch"
	"github.com/soyking/e3w/conf"
)

type e3chHandler func(*gin.Context, *client.EtcdHRCHYClient) (interface{}, error)

type groupHandler func(e3chHandler) respHandler

func withE3chGroup(e3chClt *client.EtcdHRCHYClient, config *conf.Config) groupHandler {
	return func(h e3chHandler) respHandler {
		return func(c *gin.Context) (interface{}, error) {
			clt := e3chClt
			if config.Auth {
				var err error
				username, password, authOK := c.Request.BasicAuth()
				if authOK == false || (username == "" || password == "") {
					return nil, errors.New("Not authorized")
				}

				clt, err = e3ch.CloneE3chClient(username, password, e3chClt)
				if err != nil {
					if err.Error() == "etcdserver: authentication failed, invalid user ID or password" {
						return nil, errors.New("Not authorized")
					}
					return nil, err
				}
				defer clt.EtcdClient().Close()
			}
			return h(c, clt)
		}
	}
}

type etcdHandler func(*gin.Context, *clientv3.Client) (interface{}, error)

func etcdWrapper(h etcdHandler) e3chHandler {
	return func(c *gin.Context, e3chClt *client.EtcdHRCHYClient) (interface{}, error) {
		return h(c, e3chClt.EtcdClient())
	}
}

func InitRouters(g *gin.Engine, config *conf.Config, e3chClt *client.EtcdHRCHYClient) {
	g.GET("/healthz", func(c *gin.Context) {
		c.String(200, "alface")
	})
	g.GET("/", func(c *gin.Context) {
		username, password, authOK := c.Request.BasicAuth()
		if authOK == false || (username == "" || password == "") {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.AbortWithStatus(401)
			c.JSON(http.StatusUnauthorized, map[string]string{"message": "401 errors"})
			return
		}
		c.File("./static/dist/index.html")
	})
	g.Static("/public", "./static/dist")

	e3chGroup := withE3chGroup(e3chClt, config)

	// key/value actions
	g.GET("/kv/*key", resp(e3chGroup(getKeyHandler)))
	g.POST("/kv/*key", resp(e3chGroup(postKeyHandler)))
	g.PUT("/kv/*key", resp(e3chGroup(putKeyHandler)))
	g.DELETE("/kv/*key", resp(e3chGroup(delKeyHandler)))

	// members actions
	g.GET("/members", resp(e3chGroup(etcdWrapper(getMembersHandler))))

	// roles actions
	g.GET("/roles", resp(e3chGroup(etcdWrapper(getRolesHandler))))
	g.POST("/role", resp(e3chGroup(etcdWrapper(createRoleHandler))))
	g.GET("/role/:name", resp(e3chGroup(getRolePermsHandler)))
	g.DELETE("/role/:name", resp(e3chGroup(etcdWrapper(deleteRoleHandler))))
	g.POST("/role/:name/permission", resp(e3chGroup(createRolePermHandler)))
	g.DELETE("/role/:name/permission", resp(e3chGroup(deleteRolePermHandler)))

	// users actions
	g.GET("/users", resp(e3chGroup(etcdWrapper(getUsersHandler))))
	g.POST("/user", resp(e3chGroup(etcdWrapper(createUserHandler))))
	g.GET("/user/:name", resp(e3chGroup(etcdWrapper(getUserRolesHandler))))
	g.DELETE("/user/:name", resp(e3chGroup(etcdWrapper(deleteUserHandler))))
	g.PUT("/user/:name/password", resp(e3chGroup(etcdWrapper(setUserPasswordHandler))))
	g.PUT("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(grantUserRoleHandler))))
	g.DELETE("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(revokeUserRoleHandler))))
}
