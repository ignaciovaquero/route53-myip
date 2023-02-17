package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"go.uber.org/zap"
)

const (
	envVarPrefix  = "ROUTE53_MYIP"
	defaultName   = "home.ignaciovaquero.es"
	defaultRegion = "eu-south-2"
)

var (
	sugar  *zap.SugaredLogger
	client *route53.Client
)

func getEnv(envName, defaultValue string) string {
	value := os.Getenv(fmt.Sprintf("%s_%s", envVarPrefix, envName))
	if len(value) <= 0 {
		return defaultValue
	}
	return value
}

func pointer[V comparable](value V) *V {
	return &value
}

func init() {
	var err error
	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))
	logFile := getEnv("LOG_PATH", "stdout")
	sugar, err = initLogger(debug, logFile)
	if err != nil {
		fmt.Printf("error initializing logger: %s\n", err.Error())
		os.Exit(1)
	}
	sugar.Debug("logger initialization successful")

	region := getEnv("REGION", defaultRegion)

	sugar.Debugw("creating Route53 client", "region", region)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		sugar.Fatalw("failed to load configuration", "error", err)
	}

	client = route53.NewFromConfig(cfg)
}

func main() {
	name := getEnv("URL", defaultName)
	filePath := getEnv("FILE_PATH", "./ip.txt")
	ipifyURL := getEnv("IPIFY_URL", "https://api.ipify.org")
	strTTL := getEnv("TTL", "1800")
	ttl, err := strconv.ParseInt(strTTL, 10, 64)
	if err != nil {
		sugar.Fatalw("invalid value of TTL provided", "ttl", strTTL)
	}

	sugar.Infow("getting current IP address", "ipify_url", ipifyURL)
	newIP, err := getMyIP(ipifyURL)
	if err != nil {
		sugar.Fatalw("error interacting with ipify API", "error", err.Error())
	}

	oldIP, err := readFile(filePath)
	if err != nil {
		sugar.Fatalw("error reading file", "file", filePath, "error", err.Error())
	}

	if newIP == oldIP {
		sugar.Infow("no changes in IP address", "ip", newIP)
		return
	}

	sugar.Infow("detected changes in IP address", "old_ip", oldIP, "new_ip", newIP)

	hostedZoneName := name[strings.Index(name, ".")+1:]
	ctx := context.Background()

	if err = createRecord(ctx, name, newIP, hostedZoneName, ttl); err != nil {
		sugar.Fatalw("error interacting with Route53", "error", err.Error())
	}

	sugar.Infow("created record in hosted zone",
		"hosted_zone_name", hostedZoneName,
		"name", name,
		"type", "A",
		"value", newIP,
	)

	if err = saveInFile(newIP, filePath); err != nil {
		sugar.Fatalw("error saving ip value in file", "ip", newIP, "file_path", filePath)
	}

	sugar.Infow("saved ip in file", "ip", newIP, "file_path", filePath)
}
