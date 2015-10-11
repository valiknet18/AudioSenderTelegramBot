package main 

import (
 "gopkg.in/mgo.v2"
 "log"
 "github.com/valiknet18/AudioSenderTelegramBot/config"
)

func getDatabase() *mgo.Database {  
	configFile := config.ParseConfig()
    // Connect to our local mongo
    s, err := mgo.Dial("localhost:27017")

    // Check if connection error, is mongo running?
    if err != nil {
        log.Printf("Can't connect to database")
    }

    db := s.DB(configFile.Database)
    
    return db
}