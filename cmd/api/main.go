package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tiemfah/lenswear-service/protocol"
)

func main() {
	err := protocol.ServeHTTP()
	if err != nil {
		result := fmt.Errorf("error starting ServeHTTP : %v", err)
		logrus.Error(result)
	}
}
