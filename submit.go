package unitlinq

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (c *Client) SubmitMessagesTask() {
	// First check whether there is exit signal
	var token mqtt.Token
	var msgProg messages
	//var iserror bool
	i := 0
	taken := false
	reRun := false
	for {
		select {
		case <-c.destroy:
			return
		default:
		}
		if !reRun {
			<-time.After(10 * time.Second)
		}
		// Check for unsent messages in payload
		if taken {
			fmt.Println("LOG: TOKEN ALREADY TAKEN")
			err := token.Error()
			if err != nil {
				fmt.Println(err)
				taken = false
				token = nil
				continue
			}
			sent := token.WaitTimeout(12 * time.Second)
			if !sent {
				continue
			}
			// Now update the message as sent
			msgProg.Sent = true
			c.db.Save(&msgProg)
			taken = false
		}
		var msgs []messages
		var recent []messages
		result := c.db.Order("created desc").Limit(100).Where("sent = ?", false).Find(&recent)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		for _, msg := range recent {
			token = c.MQTT.Publish(msg.Topic, 1, true, msg.Payload)
			taken = true
			msgProg = msg
			sent := token.WaitTimeout(12 * time.Second)
			if !sent {
				if !c.IsConnected() {
					fmt.Println("CONNECTION ISSUE DETECTED!")
					//iserror = true
					//taken = false
					//token = nil
				}
				break
			}
			// Now update the message as sent
			fmt.Println("MESSAGE UPLOADED SUCCESSFULLY ID: ", msg.ID, " CLEANING IN:", (100 - i), " STEPS")
			msg.Sent = true
			tx := c.db.Save(&msg)
			if tx.Error != nil {
				fmt.Println(tx.Error)
			}
			taken = false
		}
		result = c.db.Order("created asc").Limit(100).Where("sent = ?", false).Find(&msgs)
		n := result.RowsAffected
		if result.Error != nil {
			fmt.Println(result.Error)
		}

		// Now submit all the data one by one
		for _, msg := range msgs {
			token = c.MQTT.Publish(msg.Topic, 1, true, msg.Payload)
			taken = true
			msgProg = msg
			sent := token.WaitTimeout(12 * time.Second)
			if !sent {
				if !c.IsConnected() {
					fmt.Println("CONNECTION ISSUE DETECTED!")
					//taken = false
					//token = nil
					//iserror = true
				}
				break
			}
			// Now update the message as sent
			fmt.Println("MESSAGE UPLOADED SUCCESSFULLY ID: ", msg.ID, " CLEANING IN:", (100 - i), " STEPS")
			msg.Sent = true
			tx := c.db.Save(&msg)
			if tx.Error != nil {
				fmt.Println(tx.Error)
			}
			taken = false
		}
		if n == 10 {
			fmt.Println("IMMEDIATELY RUN THE SUBMISSION AGAIN!")
			reRun = true
		} else {
			reRun = false
		}
		//
		// Run cleaning periodically
		if !taken {
			if i%10 == 0 {
				c.db.Where("created < ? and sent = ?", time.Now().Add(-24*time.Hour), true).Delete(&messages{})
				c.db.Where("created < ? and sent = ?", time.Now().Add(-180*24*time.Hour), false).Delete(&messages{})
				fmt.Println("DELETED UNNCESSARY ENTRIES!")
			}
			if i >= 100 {
				c.db.Where("created < ? and sent = ?", time.Now().Add(-24*time.Hour), true).Delete(&messages{})
				c.db.Where("created < ? and sent = ?", time.Now().Add(-180*24*time.Hour), false).Delete(&messages{})
				fmt.Println("DELETED UNNCESSARY ENTRIES!")
				// Vacuum database
				c.dbLock.Lock()
				c.db.Exec("VACUUM")
				c.dbLock.Unlock()
				fmt.Println("VACUUM COMPLETED!")
				i = 0
			}
		}
		i++
	}
}
