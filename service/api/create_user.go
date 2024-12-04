package api
func (rt *_router) CreateUser(user User) (User, error) {
	dbUser, err := rt.db.CreateUser(user.ToDatabase())
	if err != nil {
		return user, err
	}

	err = user.FromDatabase(dbUser)
	if err != nil {
		return user, err
	}

	return user, nil

}