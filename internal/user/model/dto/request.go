package dto

// Model DTO create new user
type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100,alpha"`
	LastName  string `json:"last_name" validate:"required,max=100,alpha"`
	Email     string `json:"email" validate:"required,max=100,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

// UserLoginRequest Model DTO login new user
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,max=100,min=8"`
}

// UserUpdateRequest Model DTO update user
type UserUpdateRequest struct {
	FirstName   string `json:"first_name" validate:"omitempty,max=100,alpha"`
	LastName    string `json:"last_name" validate:"omitempty,max=100,alpha"`
	Email       string `json:"email" validate:"omitempty,max=100,email"`
	City        string `json:"city" validate:"omitempty,max=100,alpha"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=13,numeric"`
}

//type UserUploadAvatarRequest struct {
//	Avatar string `form:"file" validate:"required,max=512"`
//}
