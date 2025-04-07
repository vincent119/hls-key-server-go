package configs

// version get
// @Summary Get version
// @Description Get version
// @Tags Version
func Version() (x string) {
	return Conf.App.Version
}
