# 🚀 Enhanced Bitcoin Address Finder

A high-performance Bitcoin private key generator and address checker that searches for private keys corresponding to funded Bitcoin addresses from a predefined list.

## ⚠️ **IMPORTANT DISCLAIMER**

This tool is for **EDUCATIONAL PURPOSES ONLY**. Finding a private key that corresponds to a funded address is computationally infeasible due to the enormous size of the private key space (2^256). The probability of success is astronomically low.

**DO NOT USE THIS TOOL FOR MALICIOUS PURPOSES.** Only use it on addresses you own or have explicit permission to test.

## 🌟 Features

- **High Performance**: Multi-threaded architecture for maximum speed
- **Funded Address Checking**: Compares generated addresses against a list of known funded addresses
- **Real-time Notifications**: Telegram bot integration for instant alerts
- **Progress Tracking**: Detailed statistics and progress updates every 1,000,000 checks
- **Persistent Logging**: Saves all matches to a timestamped log file
- **Memory Efficient**: Loads funded addresses into memory for fast lookup

## 📋 Prerequisites

- Go 1.16 or higher
- A Telegram bot token and chat ID
- A list of funded Bitcoin addresses in `Bitcoin_addresses_LATEST.txt`

## 🛠️ Installation

1. **Clone or download the project:**
   ```bash
   git clone <repository-url>
   cd bitcoin-bruteforce-main
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Download funded addresses list:**
   - Visit [http://addresses.loyce.club](http://addresses.loyce.club)
   - Download the **LATEST** file (approximately 1600 MB)
   - The file contains all funded Bitcoin addresses sorted by balance
   - Rename the downloaded file to `Bitcoin_addresses_LATEST.txt`
   - Place it in your project directory

4. **Configure Telegram bot:**
   - Create a Telegram bot via [@BotFather](https://t.me/botfather)
   - Get your bot token
   - Get your chat ID
   - Update the constants in `enhanced-bitcoin-finder.go`:
     ```go
     const botToken = "YOUR_BOT_TOKEN_HERE"
     const chatID = "YOUR_CHAT_ID_HERE"
     ```

5. **Verify your setup:**
   - Ensure `Bitcoin_addresses_LATEST.txt` exists in the project directory
   - Each line should contain one Bitcoin address
   - Addresses should be in standard Bitcoin format (starting with 1, 3, or bc1)

## 🚀 Usage

### Basic Usage

```bash
go build enhanced-bitcoin-finder.go
./enhanced-bitcoin-finder <threads> <output-file.txt>
```

### Examples

**Run with 500 threads, save results to `found_wallets.txt`:**
```bash
./enhanced-bitcoin-finder 500 found_wallets.txt
```

**Run with 1000 threads, save results to `matches.txt`:**
```bash
./enhanced-bitcoin-finder 1000 matches.txt
```

**Run with 100 threads for testing:**
```bash
./enhanced-bitcoin-finder 100 test_results.txt
```

## 📊 How It Works

1. **Initialization**: Loads all funded addresses from `Bitcoin_addresses_LATEST.txt` into memory
2. **Key Generation**: Each worker thread generates random 256-bit private keys
3. **Address Derivation**: Converts private keys to Bitcoin addresses using ECDSA
4. **Matching**: Checks if generated addresses exist in the funded addresses list
5. **Notification**: Sends Telegram alerts for matches and progress updates
6. **Logging**: Records all matches with timestamps to the output file

## 🔧 Configuration

### Thread Count
- **Recommended**: 500-1000 threads for optimal performance
- **Maximum**: 1000 threads (hardcoded limit for stability)
- **Testing**: Start with 100 threads to verify setup

### Progress Notifications
- **Frequency**: Every 1,000,000 address checks
- **Content**: Check count, matches found, rate, elapsed time

### Output Format
```
[2024-01-15 14:30:25] FOUND! PrivateKey: abc123... Address: 1ABC123...
[2024-01-15 14:35:10] FOUND! PrivateKey: def456... Address: 1DEF456...
```

## 📱 Telegram Notifications

### Startup Message
```
🚀 Starting Bitcoin Address Finder with 500 threads
📊 Loaded 1500000 funded addresses
```

### Progress Updates
```
📊 Progress Update:
• Checked: 1000000 addresses
• Found: 0 matches
• Rate: 1250.50 checks/sec
• Elapsed: 13m20s
```

### Match Found
```
🎯 FOUND BITCOIN ADDRESS!
🔑 Private Key: abc123def456...
📍 Address: 1ABC123DEF456...
📊 Total Found: 1
```

## 📈 Performance

- **Speed**: Typically 1000-2000 checks per second per thread
- **Memory**: Loads all funded addresses into RAM for instant lookup
- **Efficiency**: No API calls or network requests during operation
- **Scalability**: Linear performance increase with thread count

## 🚨 Important Notes

### Security
- **NEVER** share your private keys
- **NEVER** use this tool on addresses you don't own
- **ALWAYS** verify the source of your funded addresses list

### Performance
- Higher thread counts may cause system instability
- Monitor CPU and memory usage during operation
- Adjust thread count based on your system capabilities

### Legal Compliance
- Ensure compliance with local laws and regulations
- Only test addresses you own or have permission to test
- Respect rate limits and system resources

## 🐛 Troubleshooting

### Common Issues

**Build Error:**
```bash
go: module lookup disabled by GOPROXY=off
```
**Solution:** Run `go mod tidy` to download dependencies

**File Not Found:**
```
Failed to load funded addresses: failed to open Bitcoin_addresses_LATEST.txt
```
**Solution:** Ensure the file exists in the project directory

**Telegram Error:**
```
Failed to send Telegram message
```
**Solution:** Verify bot token and chat ID are correct

### Performance Tips

1. **Start Small**: Begin with 100-200 threads and increase gradually
2. **Monitor Resources**: Watch CPU and memory usage
3. **Optimize Threads**: Find the sweet spot for your system
4. **Regular Updates**: Check for Go updates and dependency updates

## 📄 License

This project is for educational purposes only. Use responsibly and in compliance with applicable laws.

## 🤝 Contributing

Contributions are welcome! Please ensure any modifications maintain the educational and ethical nature of the project.

## 📞 Support

For issues or questions:
1. Check the troubleshooting section
2. Verify your configuration
3. Ensure compliance with usage guidelines

---

**Remember: This tool demonstrates the cryptographic principles behind Bitcoin. The probability of finding a match is virtually zero due to the vast size of the private key space.**
