package influx

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"sync"
	"time"
)

type InfluxClient struct {
	initLock  sync.Once
	writeLock sync.Mutex
	client    influxdb2.Client
	writeApi  map[string]api.WriteAPI
	Token     string
	Url       string
}

func (i *InfluxClient) init(ctx context.Context) error {
	var err error
	i.initLock.Do(func() {
		client := influxdb2.NewClient(i.Url, i.Token)
		i.writeApi = make(map[string]api.WriteAPI)
		i.client = client
	})
	return err
}

func (i *InfluxClient) GetWriteClient(ctx context.Context, org string, bucket string) api.WriteAPI {
	i.init(ctx)
	key := fmt.Sprintf("%s-%s", org, bucket)
	if api, ok := i.writeApi[key]; ok {
		return api
	}
	i.writeLock.Lock()
	defer i.writeLock.Unlock()

	if api, ok := i.writeApi[key]; ok {
		return api
	}

	writeApi := i.client.WriteAPI(org, bucket)
	i.writeApi[key] = writeApi
	return writeApi
}

func (i *InfluxClient) WriteMsg(ctx context.Context, org string, bucket string, measurement string, tags map[string]string, fields map[string]interface{}, d *time.Time) {
	i.init(ctx)
	api := i.GetWriteClient(ctx, org, bucket)

	p := influxdb2.NewPoint(measurement, tags, fields, *d)
	api.WritePoint(p)
	api.Flush()
}
