package services

// authentication
type Signup struct {
	FullName 				string  	`json:"fullName" validate:"required"`
	Email 					string 		`json:"email" validate:"required,email"`
	PhoneNumber 			string 		`json:"phoneNumber" validate:"required,gte=11,lte=11"`
	Password 				string 		`json:"password" validate:"required,gte=8"`
	ConfirmPassword  		string 		`json:"confirmPassword" validate:"required,gte=8,eqfield=Password"`
	RegisterWithGoogle  	bool 		`json:"registerWithGoogle"`
}

type Login struct {
	Email 				string 	`json:"email" validate:"required,email"`
	Password 			string 	`json:"password" validate:"required,gte=8"`
}

type Email struct {
	Email 				string 	`json:"email" validate:"required,email"`
}

type ChangePassword struct {
	OtpCode  			string 	`json:"otpCode" validate:"required"`
	Email 				string 	`json:"email" validate:"required,email"`
	Password 			string 	`json:"password" validate:"required,gte=8"`
	ConfirmPassword  	string 	`json:"confirmPassword" validate:"required,gte=8,eqfield=Password"`
}

type VerifyOtp struct {
	Email 				string 	`json:"email" validate:"required,email"`
	OtpCode  			string 	`json:"otpCode" validate:"required"`
}


type EditProfile struct {
	Name  			string 		`json:"fullName" validate:"required"`
	Address 		*string 	`json:"address"`
	PhoneNumber 	*string 	`json:"phoneNumber"`
	Location 		*string 	`json:"location"`
}