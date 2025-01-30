package templates

import "fmt"


func OtpEmail(fullName string, otpCode string, page string) string {

	title := "Verify Your Account"
	if page == "login" {
		title = "Verify Your Email"
	} else if page == "forgotPassword" {
		title = "Reset Your Password"
	}

	fmt.Println("title", title)
	
	return `
	<html lang='en'>
	<head>
		<meta charset='UTF-8' />
		<meta name='viewport' content='width=device-width, initial-scale=1.0' />
		<title> ` + title + ` </title>
		<style>
		body { font-family: 'Trebuchet MS', 'Lucida Sans Unicode', Arial,
		sans-serif; background-color: #fff; margin: 0; padding: 0; display: flex;
		justify-content: center; align-items: center; height: 100%; } .container {
		width: 100%; } .header { padding: 10px; border-bottom: 1px solid #C4C4C44D
		; } .header img { width: 150px; height: 160px;} .content { width: 100%;
		padding: 10px; font-weight: 500; color: #8C8C8C;} h1 { font-size: 20px;
		font-weight: 600; color: #000; margin-bottom: 10px; text-align: left; } p
		{text-align: left; margin: 10px 0; font-size: 16px; line-height: 1.6;
		color: } .otp { font-size: 24px; font-weight: bold; color: #114084;
		margin: 10px 0; } .otp-validity { font-size: 16px; } .footer {
		background-color: #114084; padding: 10px; width: 100%; } .footer p {
		color: white; font-size: 12px; margin-top: 5px; } .footer-icons img {
		width: 20px; height: 20px; } .support-text { border-top: 1px solid
		#C4C4C44D; width: 100%; font-size: 12px; font-weight: 500; color: #8C8C8C;
		text-align: left; padding: 10px 10px; } .support-text a { color: #114084;
		text-decoration: underline; } strong{ color: #000; }
		</style>
	</head>
	<body>
		<div class='container'>
		<div class='header'>
			<img
			src='https://firebasestorage.googleapis.com/v0/b/themechanicnextdoor-f7e12.appspot.com/o/tmnd-logo.jpg?alt=media&token=13214a2f-30d9-4815-b19d-730f9d6f74f8'
			alt='Logo'
			/>
		</div>

		<div class='content'>
			<h1>Verify Your Account</h1>
			<p>Hello ` + fullName + `,</p>

			<p>Thanks for signing up!</p>
			<p>Please use the following One Time Password (OTP)</p>

			<div class='otp'> `+ otpCode + ` </div>
			<p class='otp-validity'>
			This passcode will only be valid for the next
			<strong>15 minutes</strong>.
			</p>
		</div>

		<div class='support-text'>
			Need help? Contact our support team at
			<a
			href='mailto:themechanicnextdoor@gmail.com'
			>themechanicnextdoor@gmail.com</a>
		</div>
		<div class='footer'>
			<div class='footer-icons'>
			<a
				target='_blank'
				href='https://www.instagram.com/themechanicnextdoor/profilecard/?igsh=em9uZHppNWdldWM5'
			><img
				src='https://firebasestorage.googleapis.com/v0/b/themechanicnextdoor-f7e12.appspot.com/o/ri_instagram-fill.png?alt=media&token=b43c3b89-f0f3-4a3b-aa2f-e1c53149c4e3'
				alt='Instagram'
				/></a>
			</div>
			<p>&copy; 2025 SOSA. All rights reserved.</p>

		</div>
		</div>
	</body>
	</html>
	`
}
