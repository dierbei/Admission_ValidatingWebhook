package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsKeyName  = "tls.key"
	tlsCertName = "tls.crt"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/validate", Validation)

	certDir := os.Getenv("CERT_DIR")

	if err := http.ListenAndServeTLS(":8080", fmt.Sprintf("%s/%s", certDir, tlsCertName), fmt.Sprintf("%s/%s", certDir, tlsKeyName), r); err != nil {
		panic(err)
	}
}

// func Validation(w http.ResponseWriter, r *http.Request) {
func Validation(ctx *gin.Context) {
	// 从 Context 对象中获取请求和响应信息
	r := ctx.Request
	w := ctx.Writer

	ar := new(admission.AdmissionReview)
	if err := json.NewDecoder(r.Body).Decode(ar); err != nil {
		handleError(w, nil, err)
		return
	}

	response := &admission.AdmissionResponse{
		UID:     ar.Request.UID,
		Allowed: true,
	}

	pod := &corev1.Pod{}
	if err := json.Unmarshal(ar.Request.Object.Raw, pod); err != nil {
		handleError(w, ar, err)
		return
	}

	// (?m) 表示多行模式，即将 ^ 和 $ 的匹配模式从匹配整个字符串改为匹配每一行的开头和结尾；
	// (nginx|nginx:\S+) 表示匹配以 nginx 或者以 nginx: 开头的字符串。其中 | 表示或的意思，\S 表示匹配任意非空白字符，+ 表示匹配一个或多个前面的表达式。
	re := regexp.MustCompile(`(?m)(nginx|nginx:\S+)`)
	for _, c := range pod.Spec.Containers {
		if !re.MatchString(c.Image) {
			response.Allowed = false
			break
		}
	}

	responseAR := &admission.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: response,
	}
	json.NewEncoder(w).Encode(responseAR)
}

func handleError(w http.ResponseWriter, ar *admission.AdmissionReview, err error) {
	if err != nil {
		log.Println("[Error]", err.Error())
	}
	response := &admission.AdmissionResponse{
		Allowed: false,
	}
	if ar != nil {
		response.UID = ar.Request.UID
	}
	ar.Response = response
	json.NewEncoder(w).Encode(ar)
}
