package http

import (
    "errors"
    "log"
    "net/http"

    "github.com/gin-gonic/contrib/static"
    "github.com/gin-gonic/gin"
    "github.com/gobwas/ws"

    "github.com/forscht/watt/db"
)

type NewChannelHandler func(channel *Channel)

// Serv starts an HTTP server at the specified address and handles incoming requests.
// It takes a callback function handleNewChannel which is called when a new WebSocket channel is established.
func Serv(
    addr string,
    cache *db.Cache,
    placeHolderValues map[string]string,
    chHandler NewChannelHandler,
) error {

    //
    // Init GIN
    //
    // Set Gin Mode
    gin.SetMode(gin.ReleaseMode)

    // Create a new Gin router with default middleware
    r := gin.Default()

    //
    // Setup routes
    //

    // Serve frontend static files
    r.Use(static.Serve("/", Embed("static", placeHolderValues, false)))
    //r.Use(static.Serve("/", static.LocalFile("./http/static", true)))

    // Handle WS connection
    r.GET("/sync/:addr", func(c *gin.Context) {
        // Get the email address from the query parameter "Addr"
        addr := c.Param("addr")
        if addr == "" {
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }
        // Upgrade the HTTP connection to a WebSocket connection
        conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
        if err != nil {
            c.Error(err)
            return
        }
        // Call the callback function with a new Channel object representing the WebSocket channel
        chHandler(&Channel{Addr: addr, Conn: conn})
    })

    // Handle API route which returns messageId
    r.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
            c.AbortWithError(http.StatusBadRequest, errors.New("id missing"))
            return
        }
        message, ok := cache.GetById(id)
        if !ok {
            c.Redirect(http.StatusMovedPermanently, "/")
            return
        }
        if message.HTML != "" {
            c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(message.HTML))
            return
        }
        if message.Text != "" {
            c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(message.Text))
            return
        }
        c.Redirect(http.StatusMovedPermanently, "/")
    })

    // Redirect to root
    r.NoRoute(func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/")
    })

    //
    // Start the server
    //

    // Listen for incoming connections
    log.Println("HTTP server listening on", addr)
    return http.ListenAndServe(addr, r)
}
