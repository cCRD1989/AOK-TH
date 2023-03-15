package controller

type User struct{}

// func (u User) ChaeckLoginAll(ctx *gin.Context) {

// 	var users []model.USER_CHECK_LOGIN

// 	qurel := db.Conn_WLK.Find(&users)
// 	if qurel.Error != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"Status": "error", "message": "ChaeckLoginAll Failed"})
// 	}

// 	ctx.JSON(http.StatusNotFound, gin.H{"Status": "ok", "message": "ChaeckLoginAll", "users": users})

// }

// func (u User) Profile(ctx *gin.Context) {

// 	userId := ctx.MustGet("userId").(float64)
// 	var user model.USER_CHECK_LOGIN

// 	qurel := db.Conn_WLK.First(&user, userId)
// 	if qurel.Error != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"Status": "error", "message": "Profile Failed"})
// 	}

// 	ctx.JSON(http.StatusNotFound, gin.H{"Status": "ok", "message": "A Profile", "users": user})

// }

// func (u User) ChaeckLogin_jwt(ctx *gin.Context) {

// 	var form dto.UserRequest
// 	if err := ctx.ShouldBindJSON(&form); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "Status": "error"})
// 		return
// 	}

// 	//Query Sql เช็คไอดีพาส ถ้าเจอก็ สร้าง jwt ส่งกลับไป
// 	var user model.USER_CHECK_LOGIN
// 	query := db.Conn_WLK.Where("username = ?", form.Username).First(&user)

// 	if err := query.Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"Status": "error", "message": "ID Failed"})
// 		return
// 	}
// 	//check pass
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"Status": "error", "message": "Password Failed"})
// 		return
// 	}

// 	// Create a new token object, specifying signing method and the claims
// 	// you would like it to contain.
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userId": user.ID,
// 		"exp":    time.Now().Add(time.Minute * 1).Unix(),
// 	})
// 	// Sign and get the complete encoded token as a string using the secret
// 	//tokenString, err := token.SignedString(os.Getenv("TOKEN_KEY"))

// 	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"Status": "error", "message": "TokenString Failed"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, dto.UserResponse{
// 		Username: form.Username,
// 		Status:   "ok",
// 		JWT:      tokenString,
// 	})

// }

// func (u User) Registered(ctx *gin.Context) {

// 	var form dto.RegiseredRequest
// 	if err := ctx.ShouldBindJSON(&form); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "Status": "error"})
// 		return
// 	}

// 	bcrypt_pass, _ := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
// 	user := model.USER_CHECK_LOGIN{
// 		Username: form.Username,
// 		Password: string(bcrypt_pass),
// 		Idcode:   form.Idcode,
// 	}
// 	if err := db.Conn_WLK.Create(&user).Error; err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "Status": "error"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, dto.RegiseredResponse{
// 		Username: user.Username,
// 		Password: user.Password,
// 		Idcode:   user.Idcode,
// 		Status:   "ok",
// 	})
// }

// func (u User) FFFFF(ctx *gin.Context) {

// 	resp, err := http.Get("https://sea-api.gold-sandbox.razer.com/ewallet/pay?for=_AGEOFKHAGAN-9988772106163201&channel=truewallet&orderid=9988772106163201&sid=8325&uid=ccrd001&denotype=&mid=&price=50THB&sig=fa52423320f0e37ae180ebf93e5ae6c5")
// 	if err != nil {
// 		fmt.Println("Error", err.Error())
// 	}

// 	fmt.Println("ตอบ", resp)
// }

//https://sea-api.gold-sandbox.razer.com/ewallet/pay?for=_AGEOFKHAGAN-9988772106163201&channel=truewallet&orderid=9988772106163201&sid=8325&uid=ccrd001&denotype=&mid=&price=50THB&sig=fa52423320f0e37ae180ebf93e5ae6c5
