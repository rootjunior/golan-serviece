package main

import "gorm.io/gorm"

type AppState struct {
	DB              *gorm.DB
	PostRepo        *PostRepository
	GetPostsUseCase *GetPostsUseCase
	Mediator        *Mediator
}

func NewAppState() *AppState {
	db := InitDB()
	repo := &PostRepository{DB: db}
	useCase := &GetPostsUseCase{Repo: repo}
	mediator := NewMediator()

	state := &AppState{
		DB:              db,
		PostRepo:        repo,
		GetPostsUseCase: useCase,
		Mediator:        mediator,
	}

	mediator.RegisterQuery(GetPostsQuery{}, func(q interface{}) (interface{}, error) {
		return useCase.Execute(q.(GetPostsQuery))
	})

	return state
}
