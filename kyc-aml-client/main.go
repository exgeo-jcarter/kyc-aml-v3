package main

import (
	kyc_aml_client "./KycAmlClient"
	"os"
	"fmt"
	"time"
	"log"
)

func main() {
	
	if len(os.Args) < 2 {
		fmt.Printf("\n"+`usage: kyc-aml-client "a name query goes here" ["address or postal code query"]`+"\n\n")
		return
	}
	
	client, err := kyc_aml_client.NewKycAmlClient("KycAmlClient/config.json")
	if err != nil {
		return
	}
	
	_, err = client.QueryDataServer("load_sdn_list", "")
	if err != nil {
		return
	}
	
	sdn_list, err := client.QueryDataServer("get_sdn_list", "")
	if err != nil {
		return
	}
	
	num_query_servers := 3
	wait_for_training_ch := make(chan int)
	
	go (func() {
		fuzzy_train_sdn_res, err := client.QueryFuzzyServer("train_sdn", sdn_list)
		if err != nil {
			return
		}
		_ = fuzzy_train_sdn_res
		wait_for_training_ch <- 1
	})()
	
	go (func() {
		metaphone_train_sdn_res, err := client.QueryMetaphoneServer("train_sdn", sdn_list)
		if err != nil {
			return
		}
		_ = metaphone_train_sdn_res
		wait_for_training_ch <- 1
	})()
	
	go (func() {
		doublemetaphone_train_sdn_res, err := client.QueryDoubleMetaphoneServer("train_sdn", sdn_list)
		if err != nil {
			return
		}
		_ = doublemetaphone_train_sdn_res
		wait_for_training_ch <- 1
	})()
	
	for i := 0; i < num_query_servers; i++ {
		<- wait_for_training_ch
	}
	
	num_queries := 0
	wait_for_queries_ch := make(chan int)
	
	aq := ""
	
	fuzzy_name_res := "{}"
	fuzzy_address_res := "{}"
	metaphone_name_res := "{}"
	metaphone_address_res := "{}"
	doublemetaphone_name_res := "{}"
	doublemetaphone_address_res := "{}"
	
	logfile, err := os.OpenFile("kyc-aml.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer logfile.Close()
	
	if os.Args[1] != "" {
		
		num_queries += num_query_servers
		
		go (func() {
			fuzzy_name_res, err = client.QueryFuzzyServer("query_name", os.Args[1])
			if err != nil {
				return
			}
			if fuzzy_name_res != "{}" {
				fmt.Printf("Fuzzy name results: %s\n", fuzzy_name_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - Fuzzy name results: " + fuzzy_name_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
		
		go (func() {
			metaphone_name_res, err = client.QueryMetaphoneServer("query_name", os.Args[1])
			if err != nil {
				return
			}
			if metaphone_name_res != "{}" {
				fmt.Printf("Metaphone name results: %s\n", metaphone_name_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - Metaphone name results: " + metaphone_name_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
		
		go (func() {
			doublemetaphone_name_res, err = client.QueryDoubleMetaphoneServer("query_name", os.Args[1])
			if err != nil {
				return
			}
			if doublemetaphone_name_res != "{}" {
				fmt.Printf("DoubleMetaphone name results: %s\n", doublemetaphone_name_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - DoubleMetaphone name results: " + doublemetaphone_name_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
	}
	
	if len(os.Args) > 2 {
		
		num_queries += num_query_servers
		aq = os.Args[2]
		
		go (func() {
			fuzzy_address_res, err = client.QueryFuzzyServer("query_address", os.Args[2])
			if err != nil {
				return
			}
			if fuzzy_address_res != "{}" {
				fmt.Printf("Fuzzy address results: %s\n", fuzzy_address_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - Fuzzy address results: " + fuzzy_address_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
		
		go (func() {
			metaphone_address_res, err = client.QueryMetaphoneServer("query_address", os.Args[2])
			if err != nil {
				return
			}
			if metaphone_address_res != "{}" {
				fmt.Printf("Metaphone address results: %s\n", metaphone_address_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - Metaphone address results: " + metaphone_address_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
		
		go (func() {
			doublemetaphone_address_res, err = client.QueryDoubleMetaphoneServer("query_address", os.Args[2])
			if err != nil {
				return
			}
			if doublemetaphone_address_res != "{}" {
				fmt.Printf("DoubleMetaphone address results: %s\n", doublemetaphone_address_res)
				
				t := time.Now().UTC().Format("2006-01-02 15:04:05")
				_, err = logfile.WriteString(t + " - DoubleMetaphone address results: " + doublemetaphone_address_res + "\n")
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
			}
			wait_for_queries_ch <- 1
		})()
	}
	
	for i := 0; i < num_queries; i++ {
		<- wait_for_queries_ch
	}
	
	risk_score, err := client.CalculateRiskScore(os.Args[1], aq, fuzzy_name_res, fuzzy_address_res, metaphone_name_res, metaphone_address_res, doublemetaphone_name_res, doublemetaphone_address_res)
	if err != nil {
		return
	}
	
	fmt.Printf("Risk score: %v\n", risk_score)
	t := time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err = logfile.WriteString(t + " - Risk score: " + fmt.Sprintf("%v", risk_score) + "\n")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	sdn_entry_res, err := client.LookupSdnEntry(fuzzy_name_res, fuzzy_address_res, metaphone_name_res, metaphone_address_res, doublemetaphone_name_res, doublemetaphone_address_res)
	if err != nil {
		return
	}
	
	fmt.Printf("SDN entry: %s\n", sdn_entry_res)
	t = time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err = logfile.WriteString(t + " - SDN entry: " + sdn_entry_res + "\n")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	
	t = time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err = logfile.WriteString(t + " ----------\n")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
}
