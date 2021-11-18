package app

func mapURLS() {
	router.GET("/api/property", handler.Get)                                // Get All Properties
	router.GET("/api/property/:id", handler.GetByID)                        // Get Properties By ID
	router.POST("/api/property", handler.Create)                            // Create a property
	router.GET("/api/property/search", handler.Search)                      // Search for properties
	router.PATCH("/api/property/update/:id", handler.Update)                // update for properties
	router.POST("/api/property/media/:id", handler.UploadMedia)             // Upload Media
	router.DELETE("/api/property/media/:id/:media_id", handler.DeleteMedia) // Delete Media
}
