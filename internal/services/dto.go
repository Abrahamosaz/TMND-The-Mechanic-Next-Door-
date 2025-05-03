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
	FullName  			string 		`json:"fullName" validate:"required"`
	PhoneNumber 		string 		`json:"phoneNumber" validate:"required"`
	Address 			*string 	`json:"address"`
	State 				*string 	`json:"state"`
	Lga 				*string 	`json:"lga"`
}


// bookings dtos
type Location struct {
	Lat 	float64		`json:"lat" validate:"required"`
	Lng 	float64		`json:"lng" validate:"required"`	
	Address string		`json:"address" validate:"required"`
}

type ServiceDetails struct {
	ServiceType  	string		`json:"type" validate:"required"`
	Services  		[]string	`json:"services"`
	Description 	*string		`json:"description"`
}

type VehicleDetails struct {
	VehicleType  	string 		`json:"type" validate:"required"`
	Brand 			string		`json:"brand" validate:"required"`
	Size 			string		`json:"size" validate:"required"`
	Model 			int			`json:"model" validate:"required"`
	Description 	*string		`json:"description"`
}

type CreateBooking struct {
	Location 					Location			`json:"location" validate:"required"`
	Date 						string				`json:"date" validate:"required"`
	ServiceDetails 				ServiceDetails 		`json:"servicesDetails" validate:"required"`
	VehicleDetails 				VehicleDetails 		`json:"vehicleDetails" validate:"required"`
	VehicleImagesUrl 			[]string			`json:"vehicleImagesUrl"`
	VehicleImagesFilename 		[]string 			`json:"vehicleImagesFilename"`
	PublicIds 					[]string			`json:"publicIds"`
}