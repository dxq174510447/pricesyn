package mock

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"pricesyn/db"
	"pricesyn/tools/taskchain"
	"strconv"
	"sync"
)

/*
CREATE TABLE `kltestdb`.`task_chain` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(200) NULL,
  `version` INT(5) NULL,
  `stage_id` VARCHAR(200) NULL,
  `stage_name` VARCHAR(200) NULL,
  `service_id` VARCHAR(200) NULL,
  `finish` INT(5) NULL,
  PRIMARY KEY (`id`)
  );
*/
type DbTaskChainService struct {
	initLock sync.Mutex
	db       *gorm.DB
	pwd      string
}

var _ taskchain.TaskChainService = (*DbTaskChainService)(nil)

func (d *DbTaskChainService) init(ctx context.Context) error {

	if d.db != nil {
		return nil
	}
	d.initLock.Lock()
	defer d.initLock.Unlock()
	if d.db != nil {
		return nil
	}

	factory := db.DbFactory{
		DbUser:  "root",
		DbPwd:   d.pwd,
		DbHost:  "10.2.10.13",
		DbPort:  3306,
		DbName:  "kltestdb",
		MaxOpen: 200,
		MaxIdle: 10,
	}
	db, err := factory.GetDb(ctx)
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *DbTaskChainService) getInstanceId(ctx context.Context, serviceId string, def *taskchain.TaskChainDef) (int64, error) {
	err := d.init(ctx)
	if err != nil {
		return 0, err
	}
	var result []*TaskChain
	err = d.db.WithContext(ctx).Table(table.TableName()).Where(map[string]interface{}{
		"name":       def.Name,
		"version":    def.Version,
		"service_id": serviceId,
	}).Order("id desc").Find(&result).Error
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("%s can't find chain instance", serviceId)
	}
	return result[0].Id, nil
}

func (d *DbTaskChainService) SaveInstance(ctx context.Context, serviceId string, def *taskchain.TaskChainDef) (string, error) {
	err := d.init(ctx)
	if err != nil {
		return "", err
	}
	t := &TaskChain{
		Name:      def.Name,
		Version:   def.Version,
		StageId:   "",
		StageName: "",
		ServiceId: serviceId,
		Finish:    0,
	}
	err = d.db.WithContext(ctx).Table(table.TableName()).Save(t).Error
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(t.Id, 10), err
}

func (d *DbTaskChainService) SaveTaskStage(ctx context.Context, serviceId string, stageId string,
	stageDef *taskchain.StageDef, def *taskchain.TaskChainDef) error {

	id, err := d.getInstanceId(ctx, serviceId, def)
	if err != nil {
		return err
	}

	err = d.db.WithContext(ctx).Table(table.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
		"stage_id":   stageId,
		"stage_name": stageDef.Name,
		"finish":     0,
	}).Error
	return err
}

func (d *DbTaskChainService) EndInstance(ctx context.Context, serviceId string, def *taskchain.TaskChainDef) error {

	id, err := d.getInstanceId(ctx, serviceId, def)
	if err != nil {
		return err
	}

	err = d.db.WithContext(ctx).Table(table.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
		"finish": "1",
	}).Error
	return err
}

func (d *DbTaskChainService) GetStageId(ctx context.Context, serviceId string, chainName string) (string, int, int, error) {
	err := d.init(ctx)
	if err != nil {
		return "", 0, 0, err
	}
	var result []*TaskChain
	err = d.db.WithContext(ctx).Table(table.TableName()).Where(map[string]interface{}{
		"name":       chainName,
		"service_id": serviceId,
	}).Order("id desc").Find(&result).Error
	if err != nil {
		return "", 0, 0, err
	}
	if len(result) == 0 {
		return "", 0, 0, nil
	}
	return result[0].StageId, result[0].Version, result[0].Finish, nil
}

var table TaskChain = TaskChain{}

type TaskChain struct {
	Id        int64  "gorm:`primaryKey;column:id`"
	Name      string "gorm:`column:name`"
	Version   int    "gorm:`column:version`"
	StageId   string "gorm:`column:stage_id`"
	StageName string "gorm:`column:stage_name`"
	ServiceId string "gorm:`column:service_id`"
	Finish    int    "gorm:`column:finish`"
}

func (m *TaskChain) TableName() string {
	return "task_chain"
}
