package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var musicalKeys [22]string = [22]string{
	"C major",
	"C minor",
	"C# major",
	"C# minor",
	"D major",
	"D minor",
	"D# major",
	"D# minor",
	"E major",
	"E minor",
	"F major",
	"F minor",
	"G major",
	"G minor",
	"G# major",
	"G# minor",
	"A major",
	"A minor",
	"A# major",
	"A# minor",
	"B major",
	"B minor",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generates a random string of length 6
func RandomUsername() string {
	return RandomString(6)
}

// Generates a random string of length 6
func RandomPassword() string {
	return RandomString(6)
}

// Generates a random string of format 6char + "@" + 6char + ".com"
func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(6) + ".com"
}

// Generates a random string of length 6
func RandomTitle() string {
	return RandomString(6)
}

func RandomGenre() string {
	return RandomString(6)
}

func RandomKey() string {
	k := len(musicalKeys)
	return musicalKeys[rand.Intn(k)]
}

func RandomBpm() int16 {
	return int16(RandomInt(20, 999))
}

func RandomTags() string {
	str := ""

	for i := 0; i < int(RandomInt(0, 5)); i++ {
		str += RandomString(int(RandomInt(3, 8))) + ","
	}
	return str
}

func RandomS3Key() string {
	return RandomString(12)
}
