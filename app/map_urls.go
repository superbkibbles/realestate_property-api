package app

const prefix = "/api/property"

func mapURLS() {
	router.GET(prefix, handler.Get)                                    // Get All Properties
	router.GET(prefix+"/:id", handler.GetByID)                         // Get Properties By ID
	router.POST(prefix, handler.Create)                                // Create a property
	router.POST(prefix+"/search", handler.Search)                      // Search for properties
	router.PATCH(prefix+"/update/:id", handler.Update)                 // update for properties
	router.POST(prefix+"/media/:id", handler.UploadMedia)              // Upload Media
	router.POST(prefix+"/property_pic/:id", handler.UploadPropertyPic) // Upload Property Picture
	router.DELETE(prefix+"/media/:id/:media_id", handler.DeleteMedia)  // Delete Media
	router.GET(prefix+"/active", handler.GetActive)                    // Get active properties
	router.GET(prefix+"/deactive", handler.GetDeactive)                // Get Deactive properties
	router.POST(prefix+"/:id/translate", handler.Translate)            // translate by id
	router.GET(prefix+"/:id/translate", handler.GetTranslated)         // translate by id
}
