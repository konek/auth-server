
package tools

import (
  "io"
  "io/ioutil"
  "encoding/json"
)

// ParseBody is a fancy ReadAll/Unarshal for reading the request's bodies.
func ParseBody(body io.ReadCloser, v interface{}) error {
  chunk, err := ioutil.ReadAll(body)
  if err != nil {
    return err
  }

  err = json.Unmarshal(chunk, v)
  return err
}

