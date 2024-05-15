package git

func (r *Repository) GetUsername() (string, error) {
	output, err := r.Execute("config", "--get", "user.name")
	return output, err
}

func (r *Repository) GetEmail() (string, error) {
	output, err := r.Execute("config", "--get", "user.email")
	return output, err
}
