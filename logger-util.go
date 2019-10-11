package main

import "os"

func LogError(s string, s2 string, s3 string) {
	f, err := os.OpenFile(s3,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	log.Info("[" + s + "-error-logging] Recording error...")
	if _, err := f.WriteString(s + ",ERROR: " + s2 + "\n"); err != nil {
		log.Println(err)
	}
}
