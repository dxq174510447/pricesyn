package mock

import (
	"context"
	"gorm.io/gorm"
	"pricesyn/db"
	"sync"
)

/*
CREATE TABLE `klraildb`.`new_table` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(200) NULL,
  `version` INT(5) NULL,
  `stage_id` VARCHAR(200) NULL,
  `finish` INT(5) NULL,
  PRIMARY KEY (`id`));
 */
type DbTaskChainService struct {
	initLock sync.Mutex
	db *gorm.DB
}

func (d DbTaskChainService) init(ctx context.Context) error {
	if d.db != nil {
		return nil
	}
	d.initLock.Lock()
	defer d.initLock.Unlock()
	if d.db != nil {
		return nil
	}
	err := db.DbFactory{
		DbUser: "root",
		DbPwd: ""
		DbHost: ""
		DbPort     int
		DbName     string
		DbLocation string
		MaxOpen    int
		MaxIdle    int
	}
}

func (d DbTaskChainService) SaveInstance(ctx context.Context, serviceId string, chainName string, chainVersion int) error {
	panic("implement me")
}

func (d DbTaskChainService) SaveTaskStage(ctx context.Context, serviceId string, chainName string, chainVersion int, stageId string) error {
	panic("implement me")
}

func (d DbTaskChainService) GetTaskId(ctx context.Context, serviceId string, chainId string) error {
	panic("implement me")
}

