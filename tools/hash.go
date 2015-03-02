
package tools

import (
  "fmt"
  "crypto/rand"

  "golang.org/x/crypto/sha3"
)

func hash(clear string) string {
  h := make([]byte, 64)
  sha3.ShakeSum256(h, []byte(clear))

  return fmt.Sprintf("%x", h)
}

// GenSalt generate an hex salt from rand.Read
func GenSalt(size int) (string, error) {
  buf := make([]byte, size)
  _, err := rand.Read(buf)
  if err != nil {
    return "", err
  }
  return fmt.Sprintf("%x", buf), nil
}

// PasswordHash generates a hash from the username, clear password and salt
func PasswordHash(username string, password string, salt string) string {
  return hash(fmt.Sprintf("%s%s%s", username, password, salt))
}

