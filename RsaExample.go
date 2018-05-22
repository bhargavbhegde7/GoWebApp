package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "os"
  "errors"
  "crypto/sha256"
  "crypto/rand"
  //"strconv"
  //"io"
  //"encoding/json"
)
var results []string
func ReceivePublicKeyFile(w http.ResponseWriter, req *http.Request){
  fmt.Println(req.FormValue("modulus"))
  fmt.Println(req.FormValue("exponent"))
}

func MessageEndPoint2(w http.ResponseWriter, req *http.Request){

  buf, err := ioutil.ReadAll(req.Body)
  if err!=nil {log.Fatal("request",err)}

  privKey, _ := ParseRsaPrivateKeyFromBinary(buf)

  priv_pem_str := ExportRsaPrivateKeyAsPemStr(privKey)
  writeToFile("priv_key", priv_pem_str)

  // pubKey_str := readKeyFromFile("pub_key")
  // pubKey, _ := convertStringPubKeyToRsaKey(pubKey_str)
  //
  // ciphertext := getEncrypted("Hello, there! from bhargav", pubKey)
  //
  // fmt.Printf("\n\n decrypted : \n%s\n", getDecrypted(ciphertext, privKey))
}

func MessageEndPoint(w http.ResponseWriter, req *http.Request){

  // buf, err := ioutil.ReadAll(req.Body)
  // if err!=nil {log.Fatal("request",err)}
  //
  // pubKey, _ := ParseRsaPublicKeyFromBinary(buf)
  //
  // pub_pem_str, _ := ExportRsaPublicKeyAsPemStr(pubKey)
  // writeToFile("pub_key", pub_pem_str)
  //
  // ciphertext := getEncrypted("Hello, there! from bhargav", pubKey)
  //
  //       w.Header().Set("Content-Type", "multipart/form-data")
  //       w.Header().Set("Content-Length", strconv.Itoa(len(ciphertext)))
  //       if _, err := w.Write(ciphertext); err != nil {
  //           log.Println("unable to write image.")
  //       }
}

func ParseRsaPublicKeyFromBinary(pubPEM []byte) (*rsa.PublicKey, error) {

    pub, err := x509.ParsePKIXPublicKey(pubPEM)
    if err != nil {
            return nil, err
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    return nil, errors.New("Key type is not RSA")
}

func ParseRsaPrivateKeyFromBinary(privPEM []byte) (*rsa.PrivateKey, error) {

    priv, err := x509.ParsePKCS1PrivateKey(privPEM)
    if err != nil {
            return nil, err
    }

    return priv, nil

    // switch priv := priv.(type) {
    // case *rsa.PrivateKey:
    //         return priv, nil
    // default:
    //         break // fall through
    // }
    // return nil, errors.New("Key type is not RSA")
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
    pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
    if err != nil {
            return "", err
    }
    pubkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PUBLIC KEY",
                    Bytes: pubkey_bytes,
            },
    )

    return string(pubkey_pem), nil
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
    privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}

func writeToFile(fileName string, text string){
  file, err := os.Create(fileName)
  if err != nil {
      log.Fatal("Cannot create file", err)
  }
  defer file.Close()

  fmt.Fprintf(file, text)
}

//encrypter
func getEncrypted(msg string, key *rsa.PublicKey) []byte{
  message := []byte(msg)
  label := []byte("")
  hash := sha256.New()
  ciphertext, err := rsa.EncryptOAEP(
      hash,
      rand.Reader,
      key,
      message,
      label)

  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  return ciphertext
}

func readKeyFromFile(fileName string) string{
  file, err := os.Open(fileName)
  if err != nil {
    fmt.Println(err)
  }
  defer file.Close()

  fileinfo, err := file.Stat()
  if err != nil {
    fmt.Println(err)
  }

  filesize := fileinfo.Size()
  buffer := make([]byte, filesize)

  _, err = file.Read(buffer)
  if err != nil {
    fmt.Println(err)
  }

  result := string(buffer)

  return result
}

func convertStringPubKeyToRsaKey(pubPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    return nil, errors.New("Key type is not RSA")
}

func getDecrypted(ciphertext []byte, key *rsa.PrivateKey) []byte{
  hash := sha256.New()
  label := []byte("")
  plainText, err := rsa.DecryptOAEP(
    hash,
    rand.Reader,
    key,
    ciphertext,
    label)
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}
return plainText
}
