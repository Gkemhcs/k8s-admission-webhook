package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValidatePrivilegedContainer checks if a container is privileged
func ValidatePrivilegedContainer(c *gin.Context) {
	var admissionReviewReq admissionv1.AdmissionReview
	if err := c.ShouldBindJSON(&admissionReviewReq); err != nil {
		logrus.Error("Failed to parse admission request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	reviewResponse := handleAdmissionRequest(admissionReviewReq.Request)
	admissionReviewResp := admissionv1.AdmissionReview{
		TypeMeta: admissionReviewReq.TypeMeta,
		Response: reviewResponse,
	}

	c.JSON(http.StatusOK, admissionReviewResp)
}

// handleAdmissionRequest processes the admission request
func handleAdmissionRequest(req *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	if req == nil {
		logrus.Error("Empty admission request received")
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result:  &v1.Status{Message: "Empty request"},
		}
	}

	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		logrus.Error("Failed to unmarshal pod JSON:", err)
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result:  &v1.Status{Message: "Invalid pod JSON"},
		}
	}

	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			logrus.Warnf("Denying privileged container in Pod: %s", pod.Name)
			return &admissionv1.AdmissionResponse{
				UID:     req.UID,
				Allowed: false,
				Result:  &v1.Status{Message: "Privileged containers are not allowed"},
			}
		}
	}

	logrus.Infof("Pod %s is allowed", pod.Name)
	return &admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
	}
}
