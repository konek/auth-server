
package tools

import (
  "strings"
)

func splitDomain(domain string) []string {
  if domain[0] == '/' {
    domain = domain[1:]
  }
  if domain[len(domain)-1] == '/' {
    domain = domain[0:len(domain)-1]
  }
  path := strings.Split(domain, "/")
  return path
}

func domainsEqual(a []string, b[]string) bool {
  for i := 0; i < len(b); i++ {
    if a[i] != b[i] {
      return false
    }
  }
  return true
}

// CheckDomains verify if domain is accessible from domains
func CheckDomains(domains []string, domain string) bool {
  path := splitDomain(domain)
  for _, d := range domains {
    s := splitDomain(d)

    if len(path) < len(s) {
      continue
    }
    if domainsEqual(path, s) {
      return true
    }
  }
  return false
}

// CheckDomain verify if domain a is accessible from domain b
func CheckDomain(a string, b string) bool {
  sa := splitDomain(a)
  sb := splitDomain(b)

  if len(sa) < len(sb) {
    return false
  }
  if domainsEqual(sa, sb) == false {
    return false
  }
  return true
}

