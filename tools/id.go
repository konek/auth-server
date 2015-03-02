
package tools

// CheckID verify if a string id is a 24-byte long hex string (mongodb id)
func CheckID(id string) bool {
  if len(id) != 24 {
    return false
  }
  for i := 0; i < len(id); i++ {
    if id[i] >= 'a' && id[i] <= 'f' {
      continue
    }
    if id[i] >= '0' && id[i] <= '9' {
      continue
    }
    return false
  }
  return true
}

