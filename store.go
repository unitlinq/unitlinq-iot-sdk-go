package unitlinq

import (
	"fmt"
	"strconv"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const ver = 1

type messages struct {
	ID      int `gorm:"primaryKey;index;autoIncrement"`
	Topic   string
	Payload []byte
	Sent    bool
	Created time.Time `gorm:"autoCreateTime:milli;index:timeidx"`
}

type config struct {
	Key   string `gorm:"primarykey;index:keyidx,unique"`
	Value string
}

func (c *Client) initStorage(storage string) error {
	var err error
	c.db, err = gorm.Open(sqlite.Open(storage+"unitlinq.db"), &gorm.Config{})
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	c.db.Exec("PRAGMA journal_mode = WAL;")
	c.db.Exec("PRAGMA journal_size_limit = 104857600;")
	if !c.db.Migrator().HasTable(&config{}) {
		// Initialize the DB for first use
		c.db.AutoMigrate(&config{})
		c.db.Create(&config{
			Key:   "ConfigVersion",
			Value: strconv.Itoa(ver)})
	}
	if !c.db.Migrator().HasTable(&messages{}) {
		// Initialize the DB for first use
		c.db.AutoMigrate(&messages{})
	}
	var dbConfig config
	c.db.First(&dbConfig, "Key = ?", "ConfigVersion")
	val, _ := strconv.Atoi(dbConfig.Value)
	if val < ver {
		c.db.Migrator().DropTable(&messages{})
		err = c.db.AutoMigrate(&messages{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) pushMessage(msg messages) (int, error) {
	m := msg
	c.dbLock.Lock()
	result := c.db.Create(&m)
	if result.Error != nil {
		return -1, result.Error
	}
	c.dbLock.Unlock()
	return m.ID, nil
}
