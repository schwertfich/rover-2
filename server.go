package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net"
)

func (r *rover) startServer(ipPort string) error {
	// Erstellt eine neue Gin-Instanz
	router := gin.Default()

	// Aktiviert CORS
	router.Use(cors.Default())
	// API-Gruppe unter /api/v1/ bereitstellen
	api := router.Group("/api")
	{
		api.GET("/:fileType", func(c *gin.Context) {
			fileType := c.Param("fileType")
			var response interface{}
			var err error

			switch fileType {
			case "plan":
				response = r.Plan
			case "rso":
				response = r.RSO
			case "map":
				response = r.Map
			case "graph":
				response = r.Graph
			default:
				c.String(400, "Please enter a valid file type: plan, rso, map, graph")
				return
			}

			if err != nil {
				c.JSON(500, gin.H{"error": "Error producing JSON", "details": err.Error()})
				return
			}
			c.JSON(200, response)
		})
	}

	// Log-Ausgabe
	log.Printf("Rover is running on %s", ipPort)

	// Listener erstellen
	l, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatal(err)
	}

	// Falls Screenshot-Feature erforderlich ist
	if r.GenImage {
		log.Printf("Starting screenshot generation for server on %s...", ipPort)
		go screenshot(ipPort)
		select {} // blockiert, um Goroutines laufen zu lassen
	}

	// Server starten
	return router.RunListener(l)
}
