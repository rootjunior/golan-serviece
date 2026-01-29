package main

type GetPostsQuery struct{}

type GetPostsUseCase struct {
	Repo *PostRepository
}

func (c *GetPostsUseCase) Execute(q GetPostsQuery) ([]PostDB, error) {
	return c.Repo.GetAll()
}
