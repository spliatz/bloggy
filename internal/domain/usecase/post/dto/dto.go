package dto

type CreatePostDTO struct {
	Content string `json:"content" binding:"required"`
}

type GetById struct {
	Id int `json:"id" binding:"required"`
}

type DeleteById struct {
	Id       int `json:"id" binding:"required"`
	AuthorId int `json:"author_id" binding:"required"`
}
