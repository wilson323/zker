// backend/infra/database/readwrite_split.go
package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ReadWriteSplitConfig 读写分离配置
type ReadWriteSplitConfig struct {
	PrimaryDSN      string
	ReplicaDSNs     []string
	Policy          dbresolver.Policy
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

// InitDBWithReadWriteSplit 初始化数据库（读写分离）
func InitDBWithReadWriteSplit(config *ReadWriteSplitConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.PrimaryDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	if len(config.ReplicaDSNs) > 0 {
		replicas := make([]gorm.Dialector, len(config.ReplicaDSNs))
		for i, dsn := range config.ReplicaDSNs {
			replicas[i] = mysql.Open(dsn)
		}

		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   config.Policy,
		}))

		if err != nil {
			return nil, err
		}

		log.Printf("Database read-write split initialized: %d replicas", len(config.ReplicaDSNs))
	}

	log.Println("Database initialized successfully")
	return db, nil
}
