package user

func (um *UserModule) SetupRoutes() {
	api := um.route.Group("/users")

	api.Get("/", um.getAllUsers)
	api.Post("/", um.createUser)
}
