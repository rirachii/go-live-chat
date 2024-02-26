package auth

// import "log"

// //best short lived access token + refresh token
// //using jwt with http-only cookie

// func main() {
// 	// init db
// 	dbConn, err = db.NewDatabase()
// 	if err != nil {
// 		log.Fatalf("Could not initialize Database")
// 	}

// 	userRep := user.NewRepository(dbConn.GetDB())
// 	userSvc := user.NewService(userRep)
// 	userHandler := user.NewHandler(userSvc)

// 	router.InitRouter(userHandler)
// 	router.Start("0.0.0.0:8080")
// }
