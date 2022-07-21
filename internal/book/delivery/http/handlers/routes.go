package handlers

func (h *bookHandlersHTTP) BookMapRoutes() {
	h.group.Use(h.mw.IsLoggedIn())
	h.group.GET("/:subject", h.FindBySubject())
}
