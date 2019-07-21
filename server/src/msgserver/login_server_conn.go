package main

import "time"

var ticker *time.Ticker

func initLoginServerConn(serverIp string, port int) {

	ticker := time.NewTicker(time.Duration(1) * time.Second)


}

func handleTimer() {
	for {
		select {
		case <-ticker.C:
			{
				break
			}
		}
	}
}
