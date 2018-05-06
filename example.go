package example
import (
        "github.com/gin-gonic/gin"
        "net/http"
        "github.com/patrickmn/go-cache"
        "time"
)
const (
        Enabled         = "enabled"
        Disabled        = "disabled"
)
type FeatureFlag struct {
        Name    string
        State string
}
func GetFlags(ch *cache.Cache) []*FeatureFlag {
        m := make([]*FeatureFlag, ch.ItemCount())
        i := 0
        for _, v := range ch.Items() {
                ff, ok := v.Object.(*FeatureFlag)
                if ok {
                        m[i] = ff
                        i = i + 1
                }
        }
        return m
}

func init() {
        r := gin.New()
        ch := cache.New(5*time.Minute, 10*time.Minute)

        r.POST("/featureflag/:name", func(c *gin.Context) {
                name := c.Param("name")
                state := c.DefaultPostForm("state", "disabled")

                if state != Enabled && state != Disabled {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "unknown state: `" + state + "`. Valid states are `" + Enabled + "` and `" + Disabled + "`"})
                } else {
                        ff := FeatureFlag{Name: name, State: state}
                        ch.Set("ff-" + name, &ff, cache.NoExpiration)
                        c.JSON(http.StatusOK, gin.H{"name": ff.Name, "state": ff.State})
                }
        })

        r.DELETE("/featureflag/:name", func(c *gin.Context) {
                name := c.Param("name")
                ch.Delete("ff-" + name)
                c.JSON(http.StatusOK, gin.H{"name": name})
        })

        r.GET("/featureflag/:name", func(c *gin.Context) {
                name := c.Param("name")
                ff, exist := ch.Get("ff-" + name)
                if exist {
                        c.JSON(http.StatusOK, ff)
                } else {
                        c.JSON(http.StatusNotFound, gin.H{"error": name + " does not exist"})
                }
        })

        r.GET("/featureflag", func(c *gin.Context) {
                c.JSON(http.StatusOK, GetFlags(ch))
        })

        r.GET("/", func(c *gin.Context) {
                c.String(http.StatusOK, "Hello World!")
        })

        http.Handle("/", r)
}
