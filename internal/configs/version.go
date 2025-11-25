package configs

// Version returns the application version.
// @Summary Get version
// @Description Get version
// @Tags Version
func Version() (x string) {
	return Conf.App.Version
}
