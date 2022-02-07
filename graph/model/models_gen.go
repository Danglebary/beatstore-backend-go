// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Beat struct {
	ID        int32     `json:"id"`
	CreatorID int32     `json:"creatorId"`
	Title     string    `json:"title"`
	Genre     string    `json:"genre"`
	Key       string    `json:"key"`
	Bpm       int32     `json:"bpm"`
	Tags      []string  `json:"tags"`
	S3Key     string    `json:"s3Key"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateBeatInput struct {
	CreatorID int32    `json:"creatorId"`
	Title     string   `json:"title"`
	Genre     string   `json:"genre"`
	Key       string   `json:"key"`
	Bpm       int32    `json:"bpm"`
	Tags      []string `json:"tags"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type DeleteBeatInput struct {
	ID int32 `json:"id"`
}

type DeleteBeatResponse struct {
	Message string `json:"message"`
}

type DeleteUserInput struct {
	ID int32 `json:"id"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type Like struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"userId"`
	BeatID int32 `json:"beatId"`
}

type UpdateBeatInput struct {
	ID    int32    `json:"id"`
	Title string   `json:"title"`
	Genre string   `json:"genre"`
	Key   string   `json:"key"`
	Bpm   int32    `json:"bpm"`
	Tags  []string `json:"tags"`
}

type UpdateUserInput struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Beats     []*Beat   `json:"beats"`
	CreatedAt time.Time `json:"createdAt"`
}
