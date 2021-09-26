package coinmonserver

func (s *CoinMonServer) configureRouter() {
	s.router.GET("/hello", s.handlerHello)
	s.router.POST("/signin", s.handlerSignIn)
	s.router.POST("/addcoin", s.handlerAddCoin)
}
