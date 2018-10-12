package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/hhkbp2/go-logging"
	"gopkg.in/mgo.v2"
)

var Config Configuration
var RedisClient *redis.Client
var MongoSession *mgo.Session
var logger logging.Logger
var _r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Configuration struct {
	DbHost            string
	Database          string
	DbUsername        string
	DbPassword        string
	DbTimeout         int
	DbName            string
	SessionTimeout    int
	ChannelBufferSize int
	RedisAddress      string
	ListenPort        string
	AccessLogPath     string
	ProxyDomain       string
	RedisPassword     string
}

func SettingsFromConfFile() bool {
	file, err := os.Open("conf.json")

	encryptedFields := []string{"RedisPassword"}

	if err != nil {
		logger.Errorf("SettingsFromConfFile failed to read file : %v", err.Error())
		return false
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		logger.Errorf("failed decoding config json : %v", err)
		return false
	}

	key := "5f4dcc3b5aa765d61d8327deb882cf99"

	r := reflect.ValueOf(&Config)

	for i := 0; i < len(encryptedFields); i++ {
		val := reflect.Indirect(r).FieldByName(encryptedFields[i]).String()

		ciphertext, err := hex.DecodeString(val)

		if err != nil {
			logger.Errorf("SettingsFromConfFile: field %s not in hex string %s", val, err.Error())
			return false
		}

		plaintext, err := decrypt(ciphertext, []byte(key))

		if err != nil {
			logger.Errorf("SettingsFromConfFile: failure decrypting the field %s due to %s", val, err.Error())
			return false
		}

		// logger.Debugf("SettingsFromConfFile: setting the field %s to %s", encryptedFields[i], plaintext)

		r.Elem().FieldByName(encryptedFields[i]).SetString(string(plaintext))
	}

	return true
}

func InitRedis() {
	if IsDevSetup() {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     strings.Split(Config.RedisAddress, ",")[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     strings.Split(Config.RedisAddress, ",")[0],
			Password: "",
			DB:       0, // use default DB
		})
	}
}

func InitMongoClient() bool {

	Host := strings.Split(Config.DbHost, ",")
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    Host,
		Database: Config.DbName,
		Username: Config.DbUsername,
		Password: Config.DbPassword,
		Timeout:  time.Second * 10,
	}
	var err error
	MongoSession, err = mgo.DialWithInfo(mongoDialInfo)

	if err != nil {
		logger.Errorf("InitMongoClient failed : %v", err.Error())
		return false
	}

	return true
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

//IsDevSetup is true or not
func IsDevSetup() bool {
	return Config.ProxyDomain == "192.68.24.21"
}
