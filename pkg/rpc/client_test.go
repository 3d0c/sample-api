package rpc

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

var starter sync.Once
var addr, saddr net.Addr

func testHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(500 * time.Millisecond)
	io.WriteString(w, "hello, world!\n")
}

func testDelayedHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(2100 * time.Millisecond)
	io.WriteString(w, "hello, world ... in a bit\n")
}

func postHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Expeted POST method", 400)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected Content-Type is application/json", 400)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if string(result) != "ping" {
		http.Error(w, "Expected `ping` command", 400)
		return
	}

	w.Write([]byte("pong"))
	return
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Expeted GET method", 400)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected Content-Type is application/json", 400)
		return
	}

	w.Write([]byte("pong"))
	return
}

func setupMockServer(t *testing.T) {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/test-delayed", testDelayedHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	rootTLSCert, err := tls.X509KeyPair([]byte(caPEM), []byte(keyPEM))
	if err != nil {
		t.Fatal(err)
	}

	lnTls, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{rootTLSCert},
	})

	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		err = http.Serve(lnTls, nil)
		if err != nil {
			t.Fatal(err)
		}
	}()

	addr = ln.Addr()
	saddr = lnTls.Addr()
}

func TestClientTimeouts(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })

	req, _ := http.NewRequest("GET", "http://"+addr.String()+"/test", nil)

	c, err := client(nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	req, _ = http.NewRequest("GET", "http://"+addr.String()+"/test-delayed", nil)

	c, err = client(&Config{time.Second * 1, time.Second * 1, nil, nil, nil, nil})
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req)
	if err == nil {
		t.Fatal("Expected error not found")
	}
}

func TestHttpPost(t *testing.T) {
	result, err := Post("http://"+addr.String()+"/post", []byte("ping"), nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "pong" {
		t.Fatalf("Expected result: pong. Obtained: %s", string(result))
	}
}

func TestHttpGet(t *testing.T) {
	result, err := Get("http://"+addr.String()+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "pong" {
		t.Fatalf("Expected result: pong. Obtained: %s", string(result))
	}
}

func TestHttpsPost(t *testing.T) {
	result, err := Post("https://"+saddr.String()+"/post", []byte("ping"),
		&Config{
			CertPEM:   []byte(clientPEM),
			KeyPEM:    []byte(clientKeyPEM),
			CaCertPEM: []byte(caPEM),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "pong" {
		t.Fatalf("Expected result: pong. Obtained: %s", string(result))
	}
}

const (
	caPEM = `-----BEGIN CERTIFICATE-----
MIIC7DCCAdSgAwIBAgIQNdrdrxPtg61341o7b3DUVDANBgkqhkiG9w0BAQsFADAX
MRUwEwYDVQQKEwxEQ0FTIFByb2plY3QwHhcNMTcwNzE4MTM1ODE5WhcNMjcwNzE4
MTM1ODE5WjAXMRUwEwYDVQQKEwxEQ0FTIFByb2plY3QwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDcnx4r2ykn9srjGcDmoi0j9PKvDQbDS+XFAVINqyen
f32E1qrGynJoT+H26U0eid2zZP9F+1S/v63W7TTvyH2Q20WQHR+NZIekalr9F8U2
IqH2NJETuFiYN5mAl6zDgchM9flaWhpSOWMK5fVXQ5+LpTG6tcnyn1v07tJBcmEp
xHHd9lYp9s520LzGhXiqQC5C5n7HE0oehtLuXo5Y87GCo9wJmR5rhTQz4y9lAL7o
b2OQTU2Qr7zo127CMm84rHE4qnrDp+xrxW1lIwRYpM1jZ0w5y5hShrHHKXZPNXFB
Dw/uUAtyu4bic8FTVPS5fuYbCt4SQRTbpgLi5x/pmgNrAgMBAAGjNDAyMA4GA1Ud
DwEB/wQEAwICBDAPBgNVHRMBAf8EBTADAQH/MA8GA1UdEQQIMAaHBH8AAAEwDQYJ
KoZIhvcNAQELBQADggEBAIJfTiSWhVldfhA3qFOp3CCd4i6pcB2nMomv/ifkRuvS
577BzLV79t3qmRGJU/j29V3FiZkvi/f0ehjuD7M7ZeIK/mjGjvZYroBquQV8PoQU
wf1OttYC0bmnzPxkcpZGAkrundRacaLI/jVdCmix93QfK4Za0ondBzo1+AYD5Ss6
q0k/Vl1XusyBpvVGxbQDrZ4JVwl675P6OsJOQU2N+rugZvaeItZI1ZFtB7idBwPx
gQjq5kU2vKEDKqiXoma/rudcrNuaHpp2Vn6MVaUy6Da8dCbzZ5fZ+YrwfxXRZSzb
+L7EY34u7SUy2KLLU8oGqhmIkH/IogBsdUlJECspFrw=
-----END CERTIFICATE-----`

	keyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA3J8eK9spJ/bK4xnA5qItI/Tyrw0Gw0vlxQFSDasnp399hNaq
xspyaE/h9ulNHonds2T/RftUv7+t1u0078h9kNtFkB0fjWSHpGpa/RfFNiKh9jSR
E7hYmDeZgJesw4HITPX5WloaUjljCuX1V0Ofi6UxurXJ8p9b9O7SQXJhKcRx3fZW
KfbOdtC8xoV4qkAuQuZ+xxNKHobS7l6OWPOxgqPcCZkea4U0M+MvZQC+6G9jkE1N
kK+86NduwjJvOKxxOKp6w6fsa8VtZSMEWKTNY2dMOcuYUoaxxyl2TzVxQQ8P7lAL
cruG4nPBU1T0uX7mGwreEkEU26YC4ucf6ZoDawIDAQABAoIBAF6Wkdb3taN+uem4
Ju465eOep1XJ3fZpWe+m60kU8oFrtaL4bmugbICwjw7PY9MOBNFfIdsdMG6tfZjC
RonviXZLrH1nHxn92fHx72THhjP5kTr/upub09AfNl7dqKKByCDG7MuCCfrqw73u
bcumIIc8oh+MxTKShFh09Aw/P80psj0wgscXDDruC2wteLtYQTfEcXgnTzduoe+8
MlTvXRyEK9wyjpqq5tjgAKAktmNQa2YpeQDv7RTzTniW6e3fudAcxYwAVHsu7OZV
t+BYcBbWBJ27JVZhL4FRrBLi7A+tq0nwI+OEqT/ZTUr4Db5lWzioV/9zXNyxcpu2
E44ZDYECgYEA5EZ7EiI/jY2UqhnFDWo8+2AL2E229Geis9G92i42vzDuKaTTXOSc
MWLRSI2dKicviqrdMEn75eVmes1UIhy1bzNdH9Lj2Jz/VlHftuY0xuckJcE1rPCM
jVr/CQff+5VFPd6m4emhtKYp3Ee55Z42cGPJdMtjdVyH78MvAdATOMsCgYEA92qq
u/PNBElxqROf4m6rfdCGosmtBv0oAQYorN4sGZarCJkoIP+wkFk3bNzFPu3nbYum
iFzewJxWhKhqKwbr8zUvIM4zNWUqM2FylibnoVLlxJEdUSNqm++EVuWjDul9hDh7
ro8N7vwtyTx29HGbys/p6/dn2g+0qsrxy/QhK+ECgYBcP280oMp19aUCKG/NQAVs
wB+JRb6NfePuLvA93zcYhDl6crVHcMr92iUg4LmGc1du/iVsgjldahrDvX4mWtun
GGalmZ+hxbAZvfReASGKz5V3/GAohv0FkqRFjf0huezFV9iwqq1CR3PbJNEmzYzK
Vkju/dIvdzkn1wSEAwYBiQKBgCHLuDakPsTvI09tFtHfPB2bdkiWM8RYoDZDmRrD
3lJAemxaP1kClCOjjCaaoXbPGGWmRcEqrmKw+EB2oMnv0BsQkLdycxxADVunW/eW
qN0obapECDUlGVLjjLgx9ev7iOGetYZKlCSo3bg3QihxvE4fyFwrF0x6CLurrQum
9UEBAoGAQFioLbAVxhZrCRaZafC7FMVLa8y2jJQkp2mfJJXULHqvrJLh1VW2dvJr
EWOWVWB7oB50mCFCr4NjOjLEl3EprqlYlgyRMRez2c8ij70d+qaLRu+ikayQz2CO
W+eTfZimmku/Zx9TJ38wVYpaqZ0xyoAWwFDWs6R1V4CNcid6LGY=
-----END RSA PRIVATE KEY-----`

	clientPEM = `-----BEGIN CERTIFICATE-----
MIIC2TCCAcGgAwIBAgIQEnh0Y1Gt4yDKasFmqqCK0zANBgkqhkiG9w0BAQsFADAX
MRUwEwYDVQQKEwxEQ0FTIFByb2plY3QwHhcNMTcwNzE4MTM1ODE5WhcNNDIwNzE4
MTM1ODE5WjAXMRUwEwYDVQQKEwxEQ0FTIFByb2plY3QwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCpcy6tcuV238hdjCygvVuGJBKEgNhf5JML4vq2+cTO
97zpvt4SrHBbKRZ+aOsP6i4/yeEL0VkVOEq6EYe5B4FBP/rUqujTzciLbssqo9Dx
iClW6LT0tGMn53puaZUvsYXz8czMqHbIE3k81VsE3uHwDWef4Mkg2p7Zm7j5kasO
nV3mL9RquHK/AZxKR05SnsZXNBuqeZOSb4AS+kKOhjyuRVCE81zO3RbwfgM61x0O
t6PsLJMNElTNhGe3tBmdwNYhfm9a7ktjj6OchLk1ds7gzZ7ATnRQ9NOTvsUAlSKx
QHY1EtowaqoLAtmJjGYk6GSZ+OY28SFG6Yeuga+2vpEtAgMBAAGjITAfMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjANBgkqhkiG9w0BAQsFAAOCAQEAhNan
Oi8n9fdT6rmxwODYV/cppgHsnvxtT9m5aIXcQyZJweSAEEuXE0qq4kcf1kdIxYYM
pMySycNNR9RtttiY3jr8zkGPtU11woN+UVS3PEIBCn98zVEA3bOznPxzkBq2kR/V
3qTR4SF2hH+XP9d86eNLZ5w9RbqsO/Mi7STnzoKaI//dt9PigNqjlJQWuAvrSO6a
smnaOvroVWxtdLOhDN9ITwgya8ModYKEH5ScSWWB6xDYBmxz9wulyrTb0mMS5dr0
cZsV7v6+2gdjUjDGnRv6YQ4hu1vjwrSjrjS3t63DUU5xfTuQ21wzklM7Zmy82xw/
EdV1YxbQtnvLXkKDqg==
-----END CERTIFICATE-----`

	clientKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAqXMurXLldt/IXYwsoL1bhiQShIDYX+STC+L6tvnEzve86b7e
EqxwWykWfmjrD+ouP8nhC9FZFThKuhGHuQeBQT/61Kro083Ii27LKqPQ8YgpVui0
9LRjJ+d6bmmVL7GF8/HMzKh2yBN5PNVbBN7h8A1nn+DJINqe2Zu4+ZGrDp1d5i/U
arhyvwGcSkdOUp7GVzQbqnmTkm+AEvpCjoY8rkVQhPNczt0W8H4DOtcdDrej7CyT
DRJUzYRnt7QZncDWIX5vWu5LY4+jnIS5NXbO4M2ewE50UPTTk77FAJUisUB2NRLa
MGqqCwLZiYxmJOhkmfjmNvEhRumHroGvtr6RLQIDAQABAoIBAQCezM0Py79z+iCi
Wr1wHkVpnuCjwMQGigWsfBkN5WElvVITlMY1pdjd9dGYweaY0ZRqP11iX90mX8zY
N5mEM55Ucgs1xr3P0OiCk2BfI4qq8DodspPzSCIswWghlV00hx/MD+0oVzCKLIXQ
Fvrnw6DmAQn98QFMgtWfhaqCaJtwY09dXjcYiAq31u7OeGk5lL7fBVE4/THW0elB
41I4Cjd6WigOObSkKmaFZPvvQ1jYQftSD8OaPrEPs6bJc+4Plh0S+c/MmvtaIL8x
DJfItvvXAJGFwFJbGUkvYx4wuZknZ9ukW3EtMRoU3dvPLfVcPyBmBkLX1706FhpZ
gDgHmYdZAoGBAMdYBgeEL8wM27ZtY0gM+6atbndXQaVADgexKM4nff98aa+F5uQo
JJlicWEALcjI919XbjBrLn3TvTkynYiAS05M/G/t4oLRmMZx1QubBCeL/1XTUFdf
UkfGQZJAJ7LuuIxKyrMZ3eKXyfMkDSo+JoXsgSB4yVkCEzVjLrrKUaXLAoGBANmc
HrpS98cEK6fVtURxvqo9H7+Gc6Yc0AazWE+VZ8IyxsUMnOcZhAkUqG8X1/hu0hWY
7It9F56YdPeHZxS/X+uLa2t2gSaS/2LwxDLcvC16t7l/gQr3AUmEj8o9ai+XT6ws
7Rr5B6SN7HpnSlATDfHDUqwVSuHP5eq5uDQJ3wXnAoGAe74OcBgEO5w3vzSEDrPQ
exTpn7dQjq4Gh4vXkrE9K50lVcm/HB43KefMDbS4twzZUhvJ0NCX2Y/YxGrBE3zg
QkjT4v8+PoqxVW5QG5YsrAfhhntYQgRv5RISniCpBA/gC2ZaEXebHw+uUvosGe7e
pv+64FdaVaBOIDPjTBxPNqUCgYEAiZrBi40fhcfLW0w2XduXd8tDIjeBNg+ONE6A
j4KopBK7wqshJLnr1lor0GRBe6WIT4PuQJ6Pqwg5HrWOp34Ex8vX15KORLg7qnMZ
fhg7Hn81YqWQEkDznWuoCXkghoumI/gczyuee54LZbOfOFd0P+cFhi/ItFZkyzM5
mh6L2w0CgYAHej+6RNxc+QI8GZjMnl3s0onWU+nU6I0215kvuRhF5PzxFkduKf4x
sqhITKH8B73fuz0J42exZ0S/ToHnkm7/JHBIG3En3oLYQQbbAJIyDMlCekHEA6Em
iiYlGPeKP1Nz5AyrIDEs22uk2/dzrLnetNIPIsu4rudPf/jj+La1+g==
-----END RSA PRIVATE KEY-----`
)
