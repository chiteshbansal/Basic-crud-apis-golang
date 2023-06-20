package model

import "time"

type Post struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatorId"`
	CreatorId   uint      `json:"creator_id"`
	Comments    []Comment `json:"comments"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Comment struct {
	Id        uint      `json:"id"`
	Content   string    `json:"content"`
	Author    User      `json:"author" gorm:"foreignKey:AuthorId"`
	CreatedAt time.Time `json:"createdAt"`
	PostId    uint      `json:"post_id"`
	AuthorId  uint      `json:"author_id"`
}
