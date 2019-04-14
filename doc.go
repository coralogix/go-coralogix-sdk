/*
Copyright 2019 Coralogix Ltd. All rights reserved.
Use of this source code is governed by a Apache 2.0
license that can be found in the LICENSE file.


Package coralogix provides SDK for sending logs to Coralogix.


The simple example:

	package main

	import (
		coralogix "github.com/coralogix/go-coralogix-sdk"
	)

	func main() {
		coralogix.SetDebug(true)

		logger := coralogix.NewCoralogixLogger(
			"YOUR_PRIVATE_KEY_HERE",
			"YOUR_APPLICATION_NAME",
			"YOUR_SUBSYSTEM_NAME",
		)
		defer logger.Destroy()

		logger.Debug("Test message 1")
		logger.Info(map[string]string{
			"text":  "Test message 2",
			"extra": "additional",
		})
		logger.Warning("Test message 4")
	}


If you want to use Coralogix SDK with Logrus logging library:

package main

import (
	coralogix "github.com/coralogix/go-coralogix-sdk"
	"github.com/sirupsen/logrus"
)

func main() {
	CoralogixHook := coralogix.NewCoralogixHook(
		"YOUR_PRIVATE_KEY_HERE",
		"YOUR_APPLICATION_NAME",
		"YOUR_SUBSYSTEM_NAME",
	)
	defer CoralogixHook.Close()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	log.AddHook(CoralogixHook)

	log.Info("Test message!")
	log.WithFields(logrus.Fields{
		"Category":   "MyCategory",
		"ClassName":  "MyClassName",
		"MethodName": "MyMethodName",
		"ThreadId":   "MyThreadId",
	}).Info("Test message 2!")
	log.WithFields(logrus.Fields{
		"extra": "additional",
	}).Info("Test message 3!")
	log.Debug("Test message 4!")
	log.Fatal("Test message 5!")
}


For a source watch https://github.com/coralogix/go-coralogix-sdk
*/
package go_coralogix_sdk
