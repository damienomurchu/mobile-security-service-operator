package mobilesecurityservice

import (
	"fmt"
	"math/rand"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//getAppLabels returns an string map with the labels which wil be associated to the kubernetes/ocp resource which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservice_cr": name}
}

//buildAppEnvVars is a helper to build the env vars which will be configured in the deployment of the Mobile Security Service Project
func buildAppEnvVars(m *mobilesecurityservicev1alpha1.MobileSecurityService) *[]corev1.EnvVar {
	res := []corev1.EnvVar{}
	for key := range getAppEnvVarsMap(m) {
		env := corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: utils.GetConfigMapName(m),
					},
					Key: key,
				},
			},
		}
		res = append(res, env)
	}
	return &res
}

//getAppEnvVarsMap is a helper to get a map[string]string with the key and values required/used to setup the Mobile Security Service Project
func getAppEnvVarsMap(m *mobilesecurityservicev1alpha1.MobileSecurityService) map[string]string {
	return map[string]string{
		"PGHOST":                           m.Spec.DatabaseHost,
		"LOG_LEVEL":                        m.Spec.LogLevel,
		"LOG_FORMAT":                       m.Spec.LogFormat,
		"ACCESS_CONTROL_ALLOW_ORIGIN":      m.Spec.AccessControlAllowOrigin,
		"ACCESS_CONTROL_ALLOW_CREDENTIALS": m.Spec.AccessControlAllowCredentials,
		"PGDATABASE":                       m.Spec.DatabaseName,
		"PGPASSWORD":                       m.Spec.DatabasePassword,
		"PGUSER":                           m.Spec.DatabaseUser,
	}
}

// getOAuthArgsMap is a helper to get the []string with values required/used to set OAuth for the Mobile Security Service Project
func getOAuthArgsMap(m *mobilesecurityservicev1alpha1.MobileSecurityService) []string {
	return []string{
		"--http-address=0.0.0.0:4180",
		"--https-address=",
		"--provider=openshift",
		fmt.Sprintf("--openshift-service-account=%s", m.Name),
		"--upstream=http://localhost:3000",
		"--cookie-secure=true",
		fmt.Sprintf("--cookie-secret=%s", RandStringBytes(16)),
		"--cookie-httponly=false",
		"--bypass-auth-for=/api/init",
		"--bypass-auth-for=/api/healthz",
		"--bypass-auth-for=/api/ping",
		"--pass-user-headers=true",
	}
}

// RandStringBytes will return a string of n random bytes
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return string(b)
}
