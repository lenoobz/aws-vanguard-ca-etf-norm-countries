package repos

import (
	"context"
	"fmt"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/config"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/consts"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/entities"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/infrastructure/repositories/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BreakdownMongo struct
type BreakdownMongo struct {
	db     *mongo.Database
	client *mongo.Client
	conf   *config.MongoConfig
	log    logger.ContextLog
}

// NewBreakdownMongo create new repository
func NewBreakdownMongo(db *mongo.Database, l logger.ContextLog, conf *config.MongoConfig) (*BreakdownMongo, error) {
	if db != nil {
		return &BreakdownMongo{
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

	return &BreakdownMongo{
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
func (r *BreakdownMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

// createContext create a new context with timeout
func createContext(ctx context.Context, t uint64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(t) * time.Millisecond
	return context.WithTimeout(ctx, timeout*time.Millisecond)
}

///////////////////////////////////////////////////////////
// Implement exposure repo interface
///////////////////////////////////////////////////////////

// FindCountriesBreakdown finds all fund country exposure
func (r *BreakdownMongo) FindCountriesBreakdown(ctx context.Context) ([]*entities.FundBreakdown, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.VANGUARD_FUND_OVERVIEW_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
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
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var funds []*entities.FundBreakdown

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var fund entities.FundBreakdown
		if err = cur.Decode(&fund); err != nil {
			r.log.Error(ctx, "decode failed", "error", err)
			return nil, err
		}

		funds = append(funds, &fund)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over cursor failed", "error", err)
		return nil, err
	}

	return funds, nil
}

// UpdateCountriesBreakdown updates all fund country exposure
func (r *BreakdownMongo) UpdateCountriesBreakdown(ctx context.Context, funds []*entities.FundBreakdown) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.ASSET_COUNTRIES_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	for _, v := range funds {
		m := models.NewFundBreakdownModel(ctx, r.log, v, r.conf.SchemaVersion)

		filter := bson.D{{
			Key:   "ticker",
			Value: v.Ticker,
		}}

		update := bson.D{
			{
				Key:   "$set",
				Value: m,
			},
			{
				Key: "$setOnInsert",
				Value: bson.D{{
					Key:   "createdAt",
					Value: time.Now().UTC().Unix(),
				}},
			},
		}

		opts := options.Update().SetUpsert(true)

		if _, err := col.UpdateOne(ctx, filter, update, opts); err != nil {
			r.log.Error(ctx, "update one failed", "error", err, "ticker", v.Ticker)
			return err
		}
	}

	return nil
}
