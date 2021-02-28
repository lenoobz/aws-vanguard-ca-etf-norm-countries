package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/config"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/infrastructure/repositories/mongodb/models"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/usecase/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExposureMongo struct
type ExposureMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.IAppLogger
	conf   *config.MongoConfig
}

// NewExposureMongo create new repository
func NewExposureMongo(db *mongo.Database, l logger.IAppLogger, conf *config.MongoConfig) (*ExposureMongo, error) {
	if db != nil {
		return &ExposureMongo{
			db:   db,
			log:  l,
			conf: conf,
		}, nil
	}

	// set context with timeout from the config
	// create new context for the query
	ctx, cancel := createContext(context.Background(), conf.TimeoutMS)
	defer cancel()

	// set mongo client options
	clientOptions := options.Client()

	// set min pool size
	if conf.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(conf.MinPoolSize)
	}

	// set max pool size
	if conf.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(conf.MaxPoolSize)
	}

	// set max idle time ms
	if conf.MaxIdleTimeMS > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(conf.MaxIdleTimeMS) * time.Millisecond)
	}

	// construct a connection string from mongo config object
	cxnString := fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.Username, conf.Password, conf.Host)

	// create mongo client by making new connection
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cxnString))
	if err != nil {
		return nil, err
	}

	return &ExposureMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    l,
		conf:   conf,
	}, nil
}

///////////////////////////////////////////////////////////
// Implement helper function
///////////////////////////////////////////////////////////

// Close disconnect from database
func (r *ExposureMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close sector mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", err)
	}
}

// createContext create a new context with timeout
func createContext(ctx context.Context, t uint64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(t) * time.Millisecond
	return context.WithTimeout(ctx, timeout*time.Millisecond)
}

func getCountryCode(name string, countryConsts []entities.Country) string {
	for _, v := range countryConsts {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			return v.Code
		}
	}
	return ""
}

///////////////////////////////////////////////////////////
// Implement exposure repo interface
///////////////////////////////////////////////////////////

// FindAllExposure finds all fund country exposure
func (r *ExposureMongo) FindAllExposure(ctx context.Context) ([]*entities.Exposure, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames["overview"]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	findOptions := options.Find()

	cur, err := col.Find(ctx, filter, findOptions)

	// only run defer function when find success
	if cur != nil {
		defer func() {
			if deferErr := cur.Close(ctx); deferErr != nil {
				err = deferErr
			}
		}()
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", err)
		return nil, err
	}

	var countryConsts []entities.Country
	if err := json.Unmarshal([]byte(consts.Countries), &countryConsts); err != nil {
		r.log.Error(ctx, "unmarshal countries failed", err)
		return nil, err
	}

	var funds []*entities.Exposure

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var fundOverviewModel entities.Exposure
		if err = cur.Decode(&fundOverviewModel); err != nil {
			r.log.Error(ctx, "decode fund overview failed")
			return nil, err
		}

		for _, v := range fundOverviewModel.CountryExposure {
			v.CountryCode = getCountryCode(v.CountryName, countryConsts)
		}

		funds = append(funds, &fundOverviewModel)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over the exposure list failed", err)
		return nil, err
	}

	return funds, nil
}

// UpdateAllExposure updates all fund country exposure
func (r *ExposureMongo) UpdateAllExposure(ctx context.Context, funds []*entities.Exposure) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames["exposure"]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
	}
	col := r.db.Collection(colname)

	for _, fund := range funds {
		fundModel := models.NewFundExposureModel(fund)

		fundModel.IsActive = true
		fundModel.Schema = r.conf.SchemaVersion
		fundModel.ModifiedAt = time.Now().UTC().Unix()

		filter := bson.D{{
			Key:   "ticker",
			Value: fund.Ticker,
		}}

		update := bson.D{{
			Key:   "$set",
			Value: fundModel,
		}, {
			Key: "$setOnInsert",
			Value: bson.D{{
				Key:   "createdAt",
				Value: time.Now().UTC().Unix(),
			}},
		}}

		opts := options.Update().SetUpsert(true)

		if _, err := col.UpdateOne(ctx, filter, update, opts); err != nil {
			r.log.Error(ctx, "update exposure failed", "ticker", fund.Ticker)
			return err
		}
	}

	return nil
}
