package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Telegram Bot Configuration
const botToken = "--------------------------------"
const chatID =   "--------------------------------"

// Configuration
const (
	MaxWorkers            = 1000    // Maximum concurrent workers
	ProgressNotification  = 1000000 // Send notification every 1M checks
)

// Global variables
var (
	fundedAddresses = make(map[string]bool)
	fundedMutex     sync.RWMutex
	checkedCount    int64
	foundCount      int64
	statsMutex      sync.Mutex
	startTime       time.Time
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("🚀 Enhanced Bitcoin Address Finder")
		fmt.Println("Usage: ./enhanced-bitcoin-finder <threads> <output-file.txt>")
		fmt.Println("Example: ./enhanced-bitcoin-finder 500 found_wallets.txt")
		fmt.Println("")
		fmt.Println("Features:")
		fmt.Println("• Generates random private keys and checks addresses")
		fmt.Println("• Compares against funded addresses list")
		fmt.Println("• Telegram notifications for matches and progress")
		fmt.Println("• Progress tracking and statistics")
		os.Exit(1)
	}

	numThreads, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid number of threads: %s", err)
	}

	if numThreads > MaxWorkers {
		log.Printf("Warning: Limiting threads to %d for stability", MaxWorkers)
		numThreads = MaxWorkers
	}

	outputFile := os.Args[2]

	// Load funded addresses
	if err := loadFundedAddresses(); err != nil {
		log.Fatalf("Failed to load funded addresses: %s", err)
	}

	startTime = time.Now()

	// Send startup message
	sendTelegramMessage(fmt.Sprintf("🚀 Starting Bitcoin Address Finder with %d threads\n📊 Loaded %d funded addresses", numThreads, len(fundedAddresses)))

	fmt.Printf("🔍 Starting Bitcoin Address Discovery...\n")
	fmt.Printf("📊 Loaded %d funded addresses\n", len(fundedAddresses))
	fmt.Printf("🧵 Using %d threads\n", numThreads)
	fmt.Printf("📁 Output file: %s\n", outputFile)
	
	// Start statistics reporting
	go reportStatistics()
	
	// Start searching
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go addressSearcher(i, &wg, &mutex, outputFile)
	}

	wg.Wait()
}

func loadFundedAddresses() error {
	file, err := os.Open("Bitcoin_addresses_LATEST.txt")
	if err != nil {
		return fmt.Errorf("failed to open Bitcoin_addresses_LATEST.txt: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	
	for scanner.Scan() {
		address := strings.TrimSpace(scanner.Text())
		if len(address) > 0 {
			fundedMutex.Lock()
			fundedAddresses[address] = true
			fundedMutex.Unlock()
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}

	log.Printf("✅ Loaded %d funded addresses", count)
	return nil
}

func addressSearcher(id int, wg *sync.WaitGroup, mutex *sync.Mutex, outputFile string) {
	defer wg.Done()
	
	for {
		// Generate random private key
		privateKey, err := generateRandomPrivateKey()
		if err != nil {
			log.Printf("Worker %d: Failed to generate private key: %s", id, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		
		// Convert to address
		address, err := privateKeyToAddress(privateKey)
		if err != nil {
			log.Printf("Worker %d: Failed to convert private key to address: %s", id, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		
		// Check if address is in funded list
		fundedMutex.RLock()
		isFunded := fundedAddresses[address]
		fundedMutex.RUnlock()
		
		// Update statistics
		statsMutex.Lock()
		checkedCount++
		currentChecked := checkedCount
		statsMutex.Unlock()
		
		// Progress notification every 1M checks
		if currentChecked%ProgressNotification == 0 {
			elapsed := time.Since(startTime)
			rate := float64(currentChecked) / elapsed.Seconds()
			message := fmt.Sprintf("📊 Progress Update:\n• Checked: %d addresses\n• Found: %d matches\n• Rate: %.2f checks/sec\n• Elapsed: %s", 
				currentChecked, foundCount, rate, elapsed.Round(time.Second))
			sendTelegramMessage(message)
		}
		
		if isFunded {
			// Found a match!
			statsMutex.Lock()
			foundCount++
			statsMutex.Unlock()
			
			// Convert private key to hex
			privateKeyHex := hex.EncodeToString(privateKey)
			
			// Send Telegram notification
			message := fmt.Sprintf("🎯 FOUND BITCOIN ADDRESS!\n🔑 Private Key: %s\n📍 Address: %s\n📊 Total Found: %d", 
				privateKeyHex, address, foundCount)
			sendTelegramMessage(message)
			
			// Record in file
			mutex.Lock()
			file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Worker %d: Failed to open file: %s", id, err)
				mutex.Unlock()
				continue
			}
			
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			file.WriteString(fmt.Sprintf("[%s] FOUND! PrivateKey: %s Address: %s\n", 
				timestamp, privateKeyHex, address))
			file.Close()
			mutex.Unlock()
			
			fmt.Printf("🎯 Worker %d: FOUND MATCH! PrivateKey: %s Address: %s\n", 
				id, privateKeyHex, address)
		}
		
		// Small delay to avoid overwhelming the system
		time.Sleep(10 * time.Millisecond)
	}
}

func generateRandomPrivateKey() ([]byte, error) {
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func privateKeyToAddress(privateKeyBytes []byte) (string, error) {
	// Convert private key to public key
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(privateKeyBytes)
	
	// Create ECDSA public key
	publicKey := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	
	return publicKeyToAddress(publicKey)
}

func publicKeyToAddress(publicKey ecdsa.PublicKey) (string, error) {
	pubKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	sha256Hash := sha256.New()
	sha256Hash.Write(pubKeyBytes)
	sha256Result := sha256Hash.Sum(nil)

	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Result)
	ripemd160Result := ripemd160Hash.Sum(nil)

	networkVersion := byte(0x00)
	addressBytes := append([]byte{networkVersion}, ripemd160Result...)
	checksum := sha256Checksum(addressBytes)
	fullAddress := append(addressBytes, checksum...)

	return base58.Encode(fullAddress), nil
}

func sha256Checksum(input []byte) []byte {
	firstSHA := sha256.New()
	firstSHA.Write(input)
	result := firstSHA.Sum(nil)

	secondSHA := sha256.New()
	secondSHA.Write(result)
	finalResult := secondSHA.Sum(nil)

	return finalResult[:4]
}

func reportStatistics() {
	for {
		time.Sleep(30 * time.Second)
		
		statsMutex.Lock()
		elapsed := time.Since(startTime)
		rate := float64(checkedCount) / elapsed.Seconds()
		stats := fmt.Sprintf("📊 Progress Report:\n"+
			"• Checked: %d addresses\n"+
			"• Found: %d matches\n"+
			"• Rate: %.2f checks/sec\n"+
			"• Elapsed: %s",
			checkedCount, foundCount, rate, elapsed.Round(time.Second))
		statsMutex.Unlock()
		
		log.Println(stats)
	}
}

func sendTelegramMessage(message string) error {
	// URL encode the message
	encodedMessage := strings.ReplaceAll(message, " ", "%20")
	encodedMessage = strings.ReplaceAll(encodedMessage, "\n", "%0A")
	
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", 
		botToken, chatID, encodedMessage)
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
