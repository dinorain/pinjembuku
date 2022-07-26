package handlers

func (h *librarianHandlersHTTP) LibrarianMapRoutes() {
	h.group.POST("/refresh", h.RefreshToken())
	h.group.POST("/login", h.Login())

	h.group.Use(h.mw.IsLoggedIn())
	h.group.PUT("/:id", h.UpdateById(), h.mw.IsLibrarian)

	h.group.GET("/:id", h.FindById())
	h.group.GET("/me", h.GetMe(), h.mw.IsLibrarian)
	h.group.POST("/logout", h.Logout(), h.mw.IsLibrarian)

	h.group.POST("", h.Register(), h.mw.IsAdmin)
	h.group.GET("", h.FindAll(), h.mw.IsAdmin)
	h.group.DELETE("/:id", h.DeleteById(), h.mw.IsAdmin)
}
