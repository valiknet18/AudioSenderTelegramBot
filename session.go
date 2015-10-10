package main 

import (
 "gopkg.in/mgo.v2"
 "log"
)

func getSession() *mgo.Session {  
    // Connect to our local mongo
    s, err := mgo.Dial("localhost:27017/audio_sender_telegram")

    // Check if connection error, is mongo running?
    if err != nil {
        log.Printf("Can't connect to database")
    }
    
    return s
}