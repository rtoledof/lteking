package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"

	"cubawheeler.io/pkg/ably"
	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/processor"
	"cubawheeler.io/pkg/realtime"
	"cubawheeler.io/pkg/redis"
)

func NewHandler(
	client *redis.Redis,
	db *mongo.DB,
	user cubawheeler.UserService,
	exit chan struct{},
	connectionString string,
	abyKey string,
	pmConfig cubawheeler.PaymentmethodConfig,
	orderServiceURL string,
	authServiceURL string,
	walletService string,
) *handler.Server {
	resolver := &Resolver{
		user:          user,
		ads:           mongo.NewAdsService(db),
		charge:        mongo.NewChargeService(db),
		coupon:        mongo.NewCouponService(db),
		profile:       mongo.NewProfileService(db),
		processor:     processor.NewCharge(pmConfig),
		vehicle:       mongo.NewVehicleService(db),
		location:      mongo.NewLocationService(db),
		plan:          mongo.NewPlanService(db),
		message:       mongo.NewMessageService(db),
		otp:           redis.NewOtpService(client),
		ablyClient:    ably.NewClient(connectionString, exit, abyKey),
		OrderService:  orderServiceURL,
		AuthService:   authServiceURL,
		WalletService: walletService,
	}
	resolver.order = mongo.NewOrderService(db, resolver.processor, client)
	resolver.realTimeLocation = realtime.NewRealTimeService(
		redis.NewRealTimeService(client),
		resolver.ablyClient.Notifier,
		resolver.user,
		client,
		resolver.order,
	)
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}

func makeRequest(ctx context.Context, method string, url string, body url.Values) (*http.Response, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("Making request to %s", url)
	jwtToken := cubawheeler.JWTFromContext(ctx)
	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(body.Encode())
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v: %w", err, cubawheeler.ErrInternal)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v: %w", err, cubawheeler.ErrInternal)
		}
		logger.Info("Error making request to %s: rsp: %s", url, body)
		if bytes.Contains(body, []byte("error")) {
			var e cubawheeler.Error
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				return nil, cubawheeler.NewError(err, http.StatusInternalServerError, "error decoding response body")
			}
			if e.StatusCode != 0 {
				return nil, &e
			}
		}
	}
	logger.Info("Successful request to %s: status code: %d", url, resp.StatusCode)
	return resp, nil
}
