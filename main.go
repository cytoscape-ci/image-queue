package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"os"
)

func getGinRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func createDatabase() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("queue"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		q := b.Get([]byte("network_ids"))
		if q == nil {
			b.Put([]byte("network_ids"), []string)
		}
		return nil
	})
	return db
}

func setRoutes(r *gin.Engine, d *bolt.DB) {
	r.POST("/network/:network_id", func(c *gin.Context) {
		id := c.Param("network_id")
		d.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("queue"))
			q := b.Get([]byte("network_ids"))
			append(*q, id)
			err = b.Put([]byte("network_ids"), q)
			return err
		})
		c.JSON(200, gin.H{"message": "Network submission successful."})
	})

	r.GET("/network", func(c *gin.Context) {
		d.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("queue"))
			q := b.Get([]byte("network_ids"))
			len := q.Len() - 1
			network_id = (*q)[len]
			*q = (*q)[:x]
			err = b.Put([]byte("network_ids"), q)
			return err
		})
		c.JSON(200, gin.H{"network_id": network_id})
	})
}

func main() {
	database := getDatabase()
	defer database.Close()
	router := getGinRouter(database)
	setRoutes(router, database)
	router.Run(":8080")
}
