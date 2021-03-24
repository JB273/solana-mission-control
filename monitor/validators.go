package monitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PrathyushaLakkireddy/solana-prometheus/config"
	"github.com/PrathyushaLakkireddy/solana-prometheus/types"
	"github.com/PrathyushaLakkireddy/solana-prometheus/utils"
)

func GetVoteAccounts(cfg *config.Config, node string) (types.GetVoteAccountsResponse, error) {
	ops := types.HTTPOptions{
		Endpoint: cfg.Endpoints.RPCEndpoint,
		Method:   http.MethodPost,
		Body: types.Payload{Jsonrpc: "2.0", Method: "getVoteAccounts", ID: 1, Params: []interface{}{
			types.Commitment{
				Commitemnt: "recent",
			},
		}},
	}
	if node == utils.Network {
		ops.Endpoint = cfg.Endpoints.NetworkRPC
	} else if node == utils.Validator {
		ops.Endpoint = cfg.Endpoints.RPCEndpoint
	} else {
		ops.Endpoint = cfg.Endpoints.RPCEndpoint
	}

	var result types.GetVoteAccountsResponse

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting leader shedules: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error while unmarshelling leader shedules: %v", err)
		return result, err
	}

	if result.Error.Code != 0 {
		return result, fmt.Errorf("RPC error: %d %v", result.Error.Code, result.Error.Message)
	}

	return result, nil
}

func AlertStatusCountFromPrometheus(cfg *config.Config) (string, error) {
	var result types.DBRes
	var count string
	response, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=solana_val_alert_count", cfg.Prometheus.PrometheusAddress))
	if err != nil {
		log.Printf("Error: %v", err)
		return count, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(responseData, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return count, err
	}
	if len(result.Data.Result) > 0 {
		count = result.Data.Result[0].Metric.AlertCount
	}

	return count, nil
}

func GetValStatusFromDB(cfg *config.Config) (string, error) {
	var result types.DBRes
	var status string
	response, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=solana_val_status", cfg.Prometheus.PrometheusAddress))
	if err != nil {
		log.Printf("Error: %v", err)
		return status, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(responseData, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return status, err
	}
	if len(result.Data.Result) > 0 {
		status = result.Data.Result[0].Metric.SolanaValStatus
	}

	return status, nil
}

// GetCredits returns the vote credits of previous and current epoch
func GetCredits(cfg *config.Config) (string, string, error) {
	var result types.DBRes
	var cCredits, pCredits string
	response, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=solana_vote_credits", cfg.Prometheus.PrometheusAddress))
	if err != nil {
		log.Printf("Error: %v", err)
		return cCredits, pCredits, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error while reading vote credits from db : %v", err)
		return cCredits, pCredits, err
	}
	json.Unmarshal(responseData, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return cCredits, pCredits, err
	}
	if len(result.Data.Result) > 0 {
		cCredits = result.Data.Result[0].Metric.SolanaCurrentCredits
		pCredits = result.Data.Result[0].Metric.SolanaPreviousCredits
	}

	return cCredits, pCredits, nil
}
