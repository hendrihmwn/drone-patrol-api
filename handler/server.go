package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
)

type IdPath struct {
	ID string `param:"id" validate:"required,uuid4"`
}

type Server struct {
	Repository repository.RepositoryInterface
	Validator  *validator.Validate
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Validator:  validator.New(),
	}
}
