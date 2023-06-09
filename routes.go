package main

import (
	"ccrd/controller"
	"ccrd/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func serveRoutes(r *gin.Engine) {

	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "intro/intro.html", nil)
	// })

	// frontend_user
	frontend_user := controller.Frontend{}
	frontend_userGroup := r.Group("/")
	frontend_userGroup.GET("/home", frontend_user.UserGetHome) //index.html
	frontend_userGroup.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "intro/intro.html", nil)
	})

	frontend_userGroup.GET("/singin", frontend_user.UserGetSingin)                                        //index.html
	frontend_userGroup.POST("/login", frontend_user.UserGetLogin)                                         //login
	frontend_userGroup.GET("/register", frontend_user.UserGetRegister)                                    //register
	frontend_userGroup.GET("/class", frontend_user.UserGetClass)                                          //register
	frontend_userGroup.GET("/maps", frontend_user.UserGetMaps)                                            //maps
	frontend_userGroup.GET("/maps/map/:id", frontend_user.UserGetMap)                                     //maps
	frontend_userGroup.GET("/maps/mob/:id", frontend_user.UserGetMonster)                                 //maps
	frontend_userGroup.GET("/profile", middleware.UserCheck(), frontend_user.UserGetProfile)              //profile
	frontend_userGroup.POST("/profile/changpass", middleware.UserCheck(), frontend_user.UserGetChangPass) //profile
	frontend_userGroup.POST("/profile/delete", middleware.UserCheck(), frontend_user.UserGetDelete)       //profile

	frontend_userGroup.GET("/email/verify/:code", frontend_user.UserEmailVerify) //mail
	frontend_userGroup.GET("/email", func(ctx *gin.Context) {

		user := "Test001"
		Id := "1111222333"
		email := "siwanat.s@ro-legend.com"

		from := mail.NewEmail("AOK-TH", "yokoyokororog@hotmail.com")
		subject := "AOK-TH Verifying your email address."
		to := mail.NewEmail("AOK-TH", email)

		plainTextContent := `
		hello. %s 
		Please verify email
		You’re almost there! We sent an email to Click here to verify your email address. http://%s/email/verify/%x
		
		Just click on the link in that email to complete your singup. If you don’t see it, you may need to check your spam folder.
	
		`

		htmlContent := `
		<html>
			<head></head>
			<body>
				<p>hello. %s </p>
				<p>Please verify email</p>
				<p>You’re almost there! We sent an email to <a href="http://%s/email/verify/%x"><u>Click here to verify your email address.</u></a></p>
				<p></p>
				<p>Just click on the link in that email to complete your singup. If you don’t see it, you may need to check your spam folder.</p>
			</body>
		</html>
		`

		plainTextContent = fmt.Sprintf(plainTextContent, user, "ageofkhaganth.com", Id)
		htmlContent = fmt.Sprintf(htmlContent, user, "ageofkhaganth.com", Id)

		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)

		if err != nil {
			fmt.Println("ไม่สำเร็จ")
			log.Println(err)
		} else {
			fmt.Println("UserEmailVerifySend สำเร็จ")
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	}) //mail

	frontend_userGroup.GET("/newpage/:id", frontend_user.UserNewPage) //newpage

	frontend_userGroup.GET("/logout", frontend_user.UserGetLogout) //UserGetLogout
	// frontend_userGroup.GET("/download/:id", frontend_user.UserGetDownload)

	//privacypolicy
	frontend_userGroup.GET("/privacypolicy", frontend_user.UserGetPrivacypolicy)

	//service
	frontend_userGroup.GET("/service", frontend_user.UserGetService)

	// // auth https://console.cloud.google.com/
	auth_user := r.Group("/auth")
	// auth_user.GET("/", google.LoginHandler) //google add 	ctx.Redirect(http.StatusTemporaryRedirect, GetLoginURL(state))
	// auth_user.GET("/google/registered", frontend_user.Auth_google_Regis)
	// auth_user.GET("/custom", frontend_user.Auth_custom)
	auth_user.GET("/customregis", frontend_user.Auth_custom_regis)
	// auth_user.GET("/facebooklogin", frontend_user.Auth_facebook_login)
	// auth_user.GET("/facebookCall", frontend_user.Auth_facebook_call) //https://ageofkhaganth.com/auth/facebookCall/
	// auth_user.GET("/facebookRegis", frontend_user.Auth_facebook_regis)

	// private := r.Group("/auth")
	// private.Use(google.Auth())
	// private.GET("/google", frontend_user.Auth_google) //index.html

	// // AIP Razer notify
	topup_user := controller.Topup{}
	topup_Group := r.Group("/topup")
	topup_Group.GET("", topup_user.Paytopup)                 // API notify url
	topup_Group.GET("/processingpay", topup_user.PayProcess) //redirect url

	// //
	r.GET("/topups", controller.Paytopups)                // เปิดหน้าเติมเงิน
	r.POST("/topups/point", controller.PaytopupsAddPoint) // กด order

	// r.GET("/topups/:user", controller.UserCheck)
	// r.GET("/topups/play", controller.Payment) // เมื่อลูกค้า กด ออเดอร์ เข้ามา

	// //admin
	// admin_user := controller.Admin{}
	// admin_userGroup := r.Group("/admin")
	// admin_userGroup.GET("", admin_user.UserGetAdmin)      //index.html
	// admin_userGroup.GET("/items", admin_user.GetItemsAll) //index.html
	// admin_userGroup.GET("/logtopup", admin_user.Logtopup)

	// // จัดอันดับ RANKINGS ทุกอาชีพ
	// Ranking_controller := controller.Rankings{}
	// frontend_RankingGroup := r.Group("/ranking")
	// frontend_RankingGroup.GET("/:class", Ranking_controller.Ranking)

	// //Game Guides
	// guides_controller := controller.Guides{}
	// guides_userGroup := r.Group("/guide/map")
	// guides_userGroup.GET("/:maps", guides_controller.GetMap)
	// character_userGroup := r.Group("/guide/jobs")
	// character_userGroup.GET("/:chars", guides_controller.Character)

	//************
	// //category
	// categoryController := controller.Categroy{}
	// categoryGroup := r.Group("/categorys")
	// categoryGroup.GET("", categoryController.FindAll)
	// categoryGroup.GET("/:id", categoryController.FindOne)
	// categoryGroup.POST("", categoryController.Create)
	// categoryGroup.PATCH("/:id", categoryController.Update)
	// categoryGroup.DELETE("/:id", categoryController.Delete)

	// //Products
	// productController := controller.Products{}
	// productGroup := r.Group("/products")
	// productGroup.GET("", productController.FindAll)
	// productGroup.GET("/:id", productController.FindOne)
	// productGroup.POST("", productController.Create)
	// productGroup.PATCH("/:id", productController.Update)
	// productGroup.DELETE("/:id", productController.Delete)

	// //orders
	// orderController := controller.Order{}
	// orderGroup := r.Group("/orders")
	// orderGroup.GET("", orderController.FindAll)
	// orderGroup.GET("/:id", orderController.FindOne)
	// orderGroup.POST("", orderController.Create)

	// //WeloveKhan
	// userController := controller.User{}
	// userGroup := r.Group("")

	// userGroup.POST("/login", userController.ChaeckLogin_jwt)
	// userGroup.POST("/registered", userController.Registered)

	// userAll := r.Group("/users", middleware.JWTAuth())
	// userAll.GET("/readallprofile", userController.ChaeckLoginAll)
	// userAll.GET("/profile", userController.Profile)

	// userAll.GET("/111", userController.FFFFF)

}
